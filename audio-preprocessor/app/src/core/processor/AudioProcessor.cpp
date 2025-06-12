#include "core/processor/AudioProcessor.h"

#include <sndfile.h>
#include <soxr.h>
#include <cmath>
#include <stdexcept>
#include <utility>
#include <vector>

namespace audio {

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

std::vector<std::pair<float, std::vector<float>>> AudioProcessor::splitIntoChunks(
    const std::vector<float>& audio,
    const ChunkingConfig config
) {
    std::vector<std::pair<float, std::vector<float>>> chunks;

    const size_t chunk_min_size = static_cast<size_t>(config.chunk_min_duration_sec * config.sample_rate);
    const size_t chunk_max_size = static_cast<size_t>(config.chunk_max_duration_sec * config.sample_rate);
    const size_t overlap_size = static_cast<size_t>(config.overlap_duration_sec * config.sample_rate);
    const size_t step_size = chunk_max_size - overlap_size;

    size_t breakpoint = 0;
    for (size_t i = 0; i < audio.size(); i++) {
        if (std::abs(audio[i]) < config.threshold && i - breakpoint >= chunk_min_size) {
            float start_time = static_cast<float>(breakpoint) / config.sample_rate;

            chunks.emplace_back(
                start_time,
                std::vector<float>(audio.begin() + breakpoint, audio.begin() + i)
            );

            while (i < audio.size() && std::abs(audio[i]) < config.threshold)
                i++;
            
            breakpoint = i;
        } else if (i - breakpoint >= chunk_max_size) {
            size_t end = breakpoint + chunk_max_size;
            if (end > audio.size())
                end = audio.size();

            float start_time = static_cast<float>(breakpoint) / config.sample_rate;
            chunks.emplace_back(
                start_time,
                std::vector<float>(audio.begin() + breakpoint, audio.begin() + end)
            );

            breakpoint = end - overlap_size;
            i = breakpoint;
        }
    }

    if (breakpoint < audio.size()) {
        float start_time = static_cast<float>(breakpoint) / config.sample_rate;
        chunks.emplace_back(
            start_time,
            std::vector<float>(audio.begin() + breakpoint, audio.end())
        );
    }

    return chunks;
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