#include "processor/processor.h"

#include <sndfile.h>
#include <soxr.h>
#include <cmath>
#include <stdexcept>

namespace audioproc {

AudioProcessor::AudioProcessor() {}
AudioProcessor::~AudioProcessor() {}

void AudioProcessor::processWav(
    const std::string& input_path,
    const std::string& output_path,
    ProcessingConfig& config
) {
    auto audio = this->readWav(input_path, config.sample_rate);
    if (config.use_bandpass)
        this->bandPassFilter(audio, config.band_high, config.band_low, config.sample_rate);
    else
        this->highPassFilter(audio, config.high_pass, config.sample_rate);

    this->trimSilence(audio, config.threshold, config.margin);
    this->normalize(audio);
    this->fadeInOut(audio, config.fade_samples);
    this->writeWav(output_path, audio, config.sample_rate);
}

std::vector<float> AudioProcessor::processBuffer(
    const std::vector<float>& buffer,
    const ProcessingConfig& config
) {
    auto audio = buffer;

    if (config.use_bandpass)
        this->bandPassFilter(audio, config.band_high, config.band_low, config.sample_rate);
    else
        this->highPassFilter(audio, config.high_pass, config.sample_rate);

    this->trimSilence(audio, config.threshold, config.margin);
    this->normalize(audio);
    this->fadeInOut(audio, config.fade_samples);
    return audio;
}


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
    if (!file)
        throw std::runtime_error("Failed to open output WAV");
    
    sf_writef_float(file, audio.data(), audio.size());
    sf_close(file);
}

void AudioProcessor::normalize(std::vector<float>& audio) {
    float maxVal = 0.0f;

    for (float s : audio) 
        maxVal = std::max(maxVal, std::abs(s));

    if (maxVal == 0)
        return;

    for (float& s : audio)
        s /= maxVal;
}

void AudioProcessor::trimSilence(std::vector<float>& audio, const float threshold, const int margin) {
    int start = 0;
    int end = static_cast<int>(audio.size()) - 1;

    while (start < end && std::abs(audio[start]) < threshold) {
        ++start;
    }

    while (end > start && std::abs(audio[end]) < threshold) {
        --end;
    }

    start = std::max(0, start - margin);
    end = std::min(static_cast<int>(audio.size()) - 1, end + margin);

    if (end > start) {
        audio = std::vector<float>(audio.begin() + start, audio.begin() + end + 1);
    }
}

void AudioProcessor::highPassFilter(std::vector<float>& audio, const float cut_off, const int sample_rate) {
    if (audio.size() < 10)
        throw std::runtime_error("Audio is too short");

    const float PI = 3.14159265358979f;
    float omega = 2.0f * PI * cut_off / sample_rate;
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
    a0 = 1.0f;

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

void AudioProcessor::lowPassFilter(std::vector<float>& audio, const float cut_off, const int sample_rate) {
    if (audio.size() < 10)
        throw std::runtime_error("Audio is too short");

    const float PI = 3.14159265358979f;
    float omega = 2.0f * PI * cut_off / sample_rate;
    float cos_omega = cosf(omega);
    float sin_omega = sinf(omega);
    float alpha = sin_omega / (2.0f * sqrtf(2.0f)); // Q = sqrt(2)/2 (Butterworth)

    float b0 =  (1 - cos_omega) / 2;
    float b1 =   1 - cos_omega;
    float b2 =  (1 - cos_omega) / 2;
    float a0 =   1 + alpha;
    float a1 =  -2 * cos_omega;
    float a2 =   1 - alpha;

    b0 /= a0; b1 /= a0; b2 /= a0;
    a1 /= a0; a2 /= a0;
    a0 = 1.0f;

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

void AudioProcessor::bandPassFilter(
    std::vector<float>& audio,
    const float highHz,
    const float lowHz,
    const int sample_rate
) {
    if (lowHz >= highHz)
        throw std::invalid_argument("bandPassFilter: lowHz must be less than highHz");

    this->lowPassFilter(audio, highHz, sample_rate);
    this->highPassFilter(audio, lowHz, sample_rate);
}


void AudioProcessor::fadeInOut(std::vector<float>& audio, int samples) {
    if (audio.size() < 2 * samples)
        return;

    const size_t SIZE = audio.size();

    for (int i = 0; i < samples; ++i) {
        float gain = static_cast<float>(i) / (samples - 1); 
        audio[i] *= gain;
        audio[SIZE - 1 - i] *= gain;
    }
}

} // namespace audioproc