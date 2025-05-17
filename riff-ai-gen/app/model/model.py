import torch
import torch.nn as nn
import torch.nn.functional as F

class RiffGenerator(nn.Module):
    def __init__(self, vocab_size, style_size, tone_size, 
                 embed_size=256, hidden_size=512, n_layers=3):
        super().__init__()
        self.embed = nn.Embedding(vocab_size, embed_size)
        self.style_embed = nn.Embedding(style_size, embed_size)
        self.tone_embed = nn.Embedding(tone_size, embed_size)
        
        self.lstm = nn.LSTM(
            input_size=embed_size*3,
            hidden_size=hidden_size,
            num_layers=n_layers,
            dropout=0.4,
            batch_first=True
        )
        
        self.attention = nn.Linear(hidden_size + embed_size*3, hidden_size)
        self.fc = nn.Linear(hidden_size*2, vocab_size)
        
        self.dropout = nn.Dropout(0.3)
        self.init_weights()

    def init_weights(self):
        for name, param in self.named_parameters():
            if 'weight' in name:
                nn.init.xavier_normal_(param)

    def forward(self, x, styles, tones, hidden=None):
        batch_size = x.size(0)
        
        # Эмбеддинги
        style_emb = self.style_embed(styles).unsqueeze(1)
        tone_emb = self.tone_embed(tones).unsqueeze(1)
        note_emb = self.embed(x)
        
        # Конкатенация условий
        cond_emb = torch.cat([style_emb, tone_emb], dim=-1)
        cond_emb = cond_emb.expand(-1, x.size(1), -1)
        
        full_emb = torch.cat([note_emb, cond_emb], dim=-1)
        
        # LSTM + Attention
        out, hidden = self.lstm(full_emb, hidden)
        attn_weights = F.softmax(self.attention(torch.cat([out, full_emb], dim=-1)), dim=1)
        context = torch.sum(attn_weights * out, dim=1)
        
        # Final prediction
        out = self.dropout(torch.cat([out[:, -1], context], dim=-1))
        return self.fc(out), hidden

    def generate(self, style, tone, length=100, temperature=0.9):
        self.eval()
        with torch.no_grad():
            # Инициализация
            input_seq = torch.LongTensor([[self.event_vocab['<start>']]]).to(device)
            hidden = None
            
            # Конвертируем условия
            style_tensor = torch.LongTensor([style]).to(device)
            tone_tensor = torch.LongTensor([tone]).to(device)
            
            output = []
            for _ in range(length):
                out, hidden = self.forward(input_seq, style_tensor, tone_tensor, hidden)
                probs = F.softmax(out / temperature, dim=-1)
                next_token = torch.multinomial(probs, 1)
                output.append(next_token.item())
                input_seq = next_token.unsqueeze(0)
                
            return self.decode(output)
    
    def decode(self, tokens):
        inv_vocab = {v:k for k,v in self.event_vocab.items()}
        return [inv_vocab.get(t, '<unk>') for t in tokens]