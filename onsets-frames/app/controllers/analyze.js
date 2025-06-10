const onsetsFramesService = require('../services/onsets-frames')
const tf = require('@tensorflow/tfjs-node');
const mm = require('@magenta/music');

function createAnalyzer(model) {
  return async function analyzer(call, callback) {
    try {
      const { audio_bytes } = call.request.audio_data;
      const audioBuffer = Buffer.from(audio_bytes);

      const decoded = await require('wav-decoder').decode(audioBuffer);
      const audioData = decoded.channelData[0];

      const notes = await model.transcribeFromAudio(audioData);

      const responseNotes = notes.map(n => ({
        onset_seconds: n.startTime,
        midi_pitch: n.pitch,
        velocity: n.velocity ?? 100,
      }));

      callback(null, { notes: responseNotes });
    } catch (err) {
      console.error('Ошибка анализа:', err);
      callback(err);
    }
  };
}

module.exports = createAnalyzer;
