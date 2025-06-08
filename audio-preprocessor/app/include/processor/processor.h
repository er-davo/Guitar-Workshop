#pragma once

#include "config.h"

#include <vector>
#include <string>

namespace audioproc {

class AudioProcessor {
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
    
    std::vector<float> readWav(const std::string& path, int& sr);
    void writeWav(const std::string& path, const std::vector<float>& audio, int sr);
    
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