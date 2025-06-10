async function transcribe(audio_data) {
    return [
    { onset_seconds: 0.1, midi_pitch: 60, velocity: 0.8 },
    { onset_seconds: 0.5, midi_pitch: 64, velocity: 0.7 },
  ];
}

module.exports = { transcribe };