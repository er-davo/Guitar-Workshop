#include "processor/processor.h"

#include <sndfile.h>
#include <soxr.h>
#include <cmath>
#include <stdexcept>

namespace audioproc {

std::vector<float> AudioProcessor::readWav(const std::string& path, int& sr) {
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

void AudioProcessor::writeWav(const std::string& path, const std::vector<float>& audio, int sr) {
    SF_INFO info = {};
    info.samplerate = sr;
    info.channels = 1;
    info.format = SF_FORMAT_WAV | SF_FORMAT_PCM_16;

    SNDFILE* file = sf_open(path.c_str(), SFM_WRITE, &info);
    if (!file) throw std::runtime_error("Failed to open output WAV");
    sf_writef_float(file, audio.data(), audio.size());
    sf_close(file);
}

void AudioProcessor::normalize(std::vector<float>& audio) {
    float maxVal = 0.0f;
    for (float s : audio) maxVal = std::max(maxVal, std::abs(s));
    if (maxVal == 0) return;
    for (float& s : audio) s /= maxVal;
}

void AudioProcessor::trimSilence(std::vector<float>& audio) {
    int start = 0;
    int end = static_cast<int>(audio.size()) - 1;

    while (start < end && std::abs(audio[start]) < this->threshold) {
        ++start;
    }

    while (end > start && std::abs(audio[end]) < this->threshold) {
        --end;
    }

    start = std::max(0, start - this->margin);
    end = std::min(static_cast<int>(audio.size()) - 1, end + this->margin);

    if (end > start) {
        audio = std::vector<float>(audio.begin() + start, audio.begin() + end + 1);
    }
}

// High Pass Filter
void AudioProcessor::HPS(std::vector<float>& audio) {
    if (audio.size() < 10)
        throw std::runtime_error("Audio is too short");

    const float PI = 3.14159265358979f;
    float omega = 2.0f * PI * this->cutoff / this->sample_rate;
    float cos_omega = cosf(omega);
    float sin_omega = sinf(omega);
    float alpha = sin_omega / (2.0f * sqrtf(2.0f)); // Q = sqrt(2)/2 (Butterworth)

    float b0 =  (1 + cos_omega) / 2;
    float b1 = -(1 + cos_omega);
    float b2 =  (1 + cos_omega) / 2;
    float a0 =   1 + alpha;
    float a1 =  -2 * cos_omega;
    float a2 =   1 - alpha;

    b0 /= a0; b1 /= a0; b2 /= a0;
    a1 /= a0; a2 /= a0;
    a0 = 1;

    std::vector<float> out(audio.size(), 0.0f);

    for (size_t n = 2; n < audio.size(); ++n) {
        out[n] = b0 * audio[n]
            + b1 * audio[n - 1]
            + b2 * audio[n - 2]
            - a1 * out[n - 1]
            - a2 * out[n - 2];
    }

    audio = out;
}

AudioProcessor::AudioProcessor(
    float threshold,
    int margin,
    float cutoff,
    int sample_rate
) {
    this->threshold = threshold;
    this->margin = margin;
    this->cutoff = cutoff;
    this->sample_rate = sample_rate;
}

void AudioProcessor::processWav(
    const std::string& inputPath,
    const std::string& outputPath
) {
    auto audio = this->readWav(inputPath, sample_rate);
    this->HPS(audio);
    this->trimSilence(audio);
    this->normalize(audio);
    this->writeWav(outputPath, audio, sample_rate);
}

std::vector<float> AudioProcessor::processBuffer(const std::vector<float>& buffer, int sr) {
    auto audio = buffer;
    this->sample_rate = sr;
    this->HPS(audio);
    this->trimSilence(audio);
    this->normalize(audio);
    return audio;
}

} // namespace audioproc