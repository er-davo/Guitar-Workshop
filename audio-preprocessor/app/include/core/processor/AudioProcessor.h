#pragma once

#include "ProcessingConfig.h"
#include "core/audio_file/AudioFileIO.h"

#include <vector>
#include <string>

namespace audio {

class AudioProcessor : public core::AudioFileIO {
public:
    AudioProcessor();
    ~AudioProcessor();
    
    void processWav(
        const std::string& input_path,
        const std::string& output_path,
        ProcessingConfig& config
    );
    
    std::vector<float> processBuffer(
        const std::vector<float>& buffer,
        const ProcessingConfig& config
    );

    std::vector<std::vector<float>> splitIntoChunks(
        const std::vector<float>& audio,
        const int sample_rate,
        const float chunk_duration_sec,
        const float overlap_duration_sec = 0.0f
    );

    std::vector<std::vector<float>> splitIntoChunksWithQuietPriority(
        const std::vector<float>& audio,
        const int sample_rate,
        const float threshold,
        const float chunk_min_duration_sec,
        const float chunk_max_duration_sec,
        const float overlap_duration_sec = 0.0f
    );
  
    void normalize(std::vector<float>& audio);
    void trimSilence(std::vector<float>& audio, const float threshold, const int margin);

    void highPassFilter(std::vector<float>& audio, const float cut_off, const int sample_rate);
    void lowPassFilter(std::vector<float>& audio, const float cut_off, const int sample_rate);
    void bandPassFilter(
        std::vector<float>& audio,
        const float highHz,
        const float lowHz,
        const int sample_rate
    );

    void fadeInOut(std::vector<float>& audio, int samples = 512);
};

} // namespace audioproc