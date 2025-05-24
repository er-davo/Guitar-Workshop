from torch.utils.data import Dataset, DataLoader
from sklearn.model_selection import train_test_split
from midiprocessor.midi import MidiDataPreprocessor
from model.model import RiffGenerator

import torch

device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

class RiffDataset(Dataset):
    def __init__(self, sequences, styles, tones):
        self.sequences = sequences
        self.styles = styles
        self.tones = tones
        
    def __len__(self):
        return len(self.sequences)
    
    def __getitem__(self, idx):
        seq = self.sequences[idx]
        return {
            'input': torch.LongTensor(seq[:-1]),
            'target': torch.LongTensor(seq[1:]),
            'style': torch.LongTensor([self.styles[idx]]),
            'tone': torch.LongTensor([self.tones[idx]])
        }

def train():
    # Инициализация
    preprocessor = MidiDataPreprocessor()
    sequences = preprocessor.create_sequences()
    
    # Загрузка метаданных (пример)
    styles = load_styles_metadata()  # [rock, blues, ...]
    tones = load_tones_metadata()    # [E, A, ...]
    
    # Разделение данных
    X_train, X_val, y_style_train, y_style_val, y_tone_train, y_tone_val = train_test_split(
        sequences, styles, tones, test_size=0.2
    )
    
    # Даталоадеры
    train_dataset = RiffDataset(X_train, y_style_train, y_tone_train)
    train_loader = DataLoader(train_dataset, batch_size=64, shuffle=True)
    
    # Модель
    model = RiffGenerator(
        vocab_size=len(preprocessor.event_vocab),
        style_size=6,  # 5 стилей + unknown
        tone_size=5     # 4 тональности + unknown
    ).to(device)
    
    optimizer = torch.optim.AdamW(model.parameters(), lr=1e-4)
    scheduler = torch.optim.lr_scheduler.ReduceLROnPlateau(optimizer, 'min')
    
    # Цикл обучения
    for epoch in range(100):
        model.train()
        total_loss = 0
        
        for batch in train_loader:
            optimizer.zero_grad()
            
            inputs = batch['input'].to(device)
            targets = batch['target'].to(device)
            styles = batch['style'].to(device)
            tones = batch['tone'].to(device)
            
            outputs, _ = model(inputs, styles, tones)
            loss = F.cross_entropy(outputs.view(-1, outputs.size(-1)), targets.view(-1))
            
            loss.backward()
            torch.nn.utils.clip_grad_norm_(model.parameters(), 1.0)
            optimizer.step()
            
            total_loss += loss.item()
        
        # Валидация и сохранение модели
        val_loss = validate(model, X_val, y_style_val, y_tone_val)
        scheduler.step(val_loss)
        
        print(f"Epoch {epoch+1} | Train Loss: {total_loss/len(train_loader):.4f} | Val Loss: {val_loss:.4f}")
        torch.save(model.state_dict(), f"riff_generator_{epoch}.pt")