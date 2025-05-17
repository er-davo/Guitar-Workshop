import os
import pretty_midi
import numpy as np
from collections import defaultdict

class MidiDataPreprocessor:
    def __init__(self, data_dir="data/midi", seq_length=32):
        self.data_dir = data_dir
        self.seq_length = seq_length
        self.event_vocab = defaultdict(lambda: len(self.event_vocab))
        self._init_special_tokens()
        
    def _init_special_tokens(self):
        self.event_vocab['<pad>'] = 0
        self.event_vocab['<start>'] = 1
        self.event_vocab['<end>'] = 2

    def _quantize_duration(self, duration):
        # Квантование длительностей нот
        bins = [0.25, 0.5, 0.75, 1.0, 1.5, 2.0]
        return min(bins, key=lambda x: abs(x - duration))

    def process_midi(self, file_path):
        midi = pretty_midi.PrettyMIDI(file_path)
        events = []
        
        for instrument in midi.instruments:
            for note in instrument.notes:
                # Фильтрация гитарного диапазона (E2 - E5)
                if 40 <= note.pitch <= 76:  
                    quantized = self._quantize_duration(note.end - note.start)
                    events.append(f"{note.pitch}_{quantized:.2f}")
        
        return events

    def build_vocabulary(self):
        all_events = []
        for fname in os.listdir(self.data_dir):
            if fname.endswith(".mid"):
                all_events.extend(self.process_midi(os.path.join(self.data_dir, fname)))
        
        # Создаем словарь
        for event in all_events:
            _ = self.event_vocab[event]
            
        self.vocab_size = len(self.event_vocab)

    def create_sequences(self):
        sequences = []
        for fname in os.listdir(self.data_dir):
            events = self.process_midi(fname)
            encoded = [self.event_vocab[e] for e in events]
            
            # Разбиваем на последовательности фиксированной длины
            for i in range(0, len(encoded) - self.seq_length, self.seq_length//2):
                seq = encoded[i:i+self.seq_length]
                if len(seq) == self.seq_length:
                    sequences.append(seq)
        
        return np.array(sequences)