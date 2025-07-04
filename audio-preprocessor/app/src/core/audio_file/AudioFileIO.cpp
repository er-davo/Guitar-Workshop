#include "core/audio_file/AudioFileIO.h"

#include <sndfile.h>
#include <stdexcept>

namespace audio {

namespace core {

AudioFileIO::AudioFileIO() {};
AudioFileIO::~AudioFileIO() {};

std::vector<float> AudioFileIO::readWav(const std::string& path, int& sr) {
    SF_INFO info;
    SNDFILE* file = sf_open(path.c_str(), SFM_READ, &info);
    if (!file)
        throw std::runtime_error("Failed to open input WAV");

    sr = info.samplerate;
    std::vector<float> buffer(info.frames * info.channels);
    sf_readf_float(file, buffer.data(), info.frames);
    sf_close(file);

    // моно (если stereo — усредняем)
    if (info.channels == 2) {
        std::vector<float> mono(info.frames);
        for (int i = 0; i < info.frames; ++i)
            mono[i] = 0.5f * (buffer[2*i] + buffer[2*i + 1]);
        return mono;
    }

    return buffer;
}

void AudioFileIO::writeWav(const std::string& path, const std::vector<float>& audio, int sr) {
    SF_INFO info = {};
    info.samplerate = sr;
    info.channels = 1;
    info.format = SF_FORMAT_WAV | SF_FORMAT_PCM_16;

    SNDFILE* file = sf_open(path.c_str(), SFM_WRITE, &info);
    if (!file)
        throw std::runtime_error("Failed to open output WAV");

    sf_count_t frames_written = sf_writef_float(file, audio.data(), audio.size() / info.channels);

    if (frames_written != static_cast<sf_count_t>(audio.size() / info.channels)) {
        sf_close(file);
        throw std::runtime_error("Failed to write full WAV file");
    }

    sf_close(file);
}

    
} // namespace core

} // namespace audio
