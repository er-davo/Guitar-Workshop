#pragma once

#include "core/processor/AudioConfigs.h"
#include "core/audio_file/AudioFileIO.h"

#include <vector>
#include <string>
#include <utility>

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

    std::vector<std::pair<float, std::vector<float>>> splitIntoChunks(
        const std::vector<float>& audio,
        const ChunkingConfig config
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