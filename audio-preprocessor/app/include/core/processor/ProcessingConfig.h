#pragma once

namespace audio {

struct ProcessingConfig {
    float threshold = 0.01f;
    int margin = 128;

    int sample_rate = 44100;

    float high_pass = 50.0f;

    bool use_bandpass = false;
    float band_low = 80.0f;
    float band_high = 1200.0f;

    int fade_samples = 512;
};

} // namespace audioproc
