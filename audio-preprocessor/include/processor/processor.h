#pragma once

#include <vector>
#include <string>

namespace audioproc {

class AudioProcessor {
private:
    float threshold;
    int margin;
    float cutoff;
    int sample_rate;

    std::vector<float> readWav(const std::string& path, int& sr);
    void writeWav(const std::string& path, const std::vector<float>& audio, int sr);
    
    void normalize(std::vector<float>& audio);
    void trimSilence(std::vector<float>& audio);
    void HPS(std::vector<float>& audio);

public:
    AudioProcessor(
        float threshold = 0.01f,
        int margin = 128,
        float cutoff = 50.0f,
        int sample_rate = 44100
    );

    void processWav(
        const std::string& inputPath,
        const std::string& outputPath
    );

    std::vector<float> processBuffer(const std::vector<float>& buffer, int sr);
};

} // namespace audioproc