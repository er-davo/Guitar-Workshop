#include "core/processor/AudioProcessor.h"
#include "server/audio_processor_server.h"
#include "config/config.h"

#include <iostream>

int main(int argc, char* argv[]) {

#ifdef DEBUG

    // CLI

    if (argc != 3) {
        std::cerr << "Usage: audio-preprocessor input.wav output.wav\n";
        return 1;
    }

    try {
        auto proc = audioproc::AudioProcessor();
        audioproc::ProcessingConfig config;
        auto input_file = argv[1],
            output_file = argv[2];
      
        config.sample_rate = 44100;
        config.threshold = 0.01f;
        config.margin = 128;
        config.use_bandpass = true;
        config.band_low = 70.0f;
        config.band_high = 1500.0f;

        proc.processWav(
            input_file,
            output_file,
            config
        );

        std::cout << "Audio processed successfully.\n";
    } catch (const std::exception& e) {
        std::cerr << "Error: " << e.what() << '\n';
        return 1;
    }

#else
    // gRPC server

    auto config = audioproc::loadConfig();

    audioproc::RunServer("0.0.0.0:" + config.port);

#endif

    return 0;
}
