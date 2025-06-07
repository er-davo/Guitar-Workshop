#include "processor/processor.h"
#include <iostream>

int main(int argc, char* argv[]) {
    if (argc != 3) {
        std::cerr << "Usage: audio-preprocessor input.wav output.wav\n";
        return 1;
    }

    try {
        auto proc = audioproc::AudioProcessor();
        auto input_file = argv[1],
            output_file = argv[2];

        proc.processWav(input_file, output_file);

        std::cout << "Audio processed successfully.\n";
    } catch (const std::exception& e) {
        std::cerr << "Error: " << e.what() << '\n';
        return 1;
    }

    return 0;
}
