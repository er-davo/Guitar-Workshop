#pragma once

#include <string>
#include <vector>

namespace audio {

namespace core {

class AudioFileIO {
public:
    AudioFileIO();
    ~AudioFileIO();

    std::vector<float> readWav(const std::string& path, int& sr);
    void writeWav(const std::string& path, const std::vector<float>& audio, int sr);
};

} // namespace core

} // namespace audio
