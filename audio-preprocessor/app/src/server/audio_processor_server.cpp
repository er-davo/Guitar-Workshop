#include "server/audio_processor_server.h"

#include <fstream>
#include <sstream>
#include <iostream>
#include <string>
#include <cctype>

namespace audioproc {

std::string  sanitizeFilename(const std::string file_name) {
    std::string sanitizad_name;

    for (char c : file_name) {
        if (std::isalnum(static_cast<unsigned char>(c)) || 
            c == '.' || c == '_' || c == '-') {
            sanitizad_name += c;
        } else {
            sanitizad_name += '_';
        }
    }

    return sanitizad_name;
}

grpc::Status AudioProcessorServiceImpl::ProcessAudio(
    grpc::ServerContext* context,
    const ProcessAudioRequest* request,
    ProcessAudioResponse* response
) {
    std::string file_name = sanitizeFilename(request->file_name());
    std::string input_path = "temp/" + file_name;
    std::string output_path = "temp/_processed_" + file_name;

    std::ofstream input_file(input_path, std::ios::binary);
    if (!input_file.is_open()) {
        return grpc::Status(grpc::StatusCode::INTERNAL, "Failed to write input file");
    }
    input_file.write(request->wav_data().data(), request->wav_data().size());
    input_file.close();

    audio::ProcessingConfig cfg;
    cfg.threshold = request->config().threshold();
    cfg.margin = request->config().margin();
    cfg.high_pass = request->config().high_pass();
    cfg.use_bandpass = request->config().use_bandpass();
    cfg.band_low = request->config().band_low();
    cfg.band_high = request->config().band_high();
    cfg.fade_samples = request->config().fade_samples();
    cfg.sample_rate = request->config().sample_rate();

    audio::AudioProcessor proc;

    try {
        proc.processWav(input_path, output_path, cfg);
        std::remove(input_path.c_str());
    } catch (std::exception& e) {
        return grpc::Status(grpc::StatusCode::INTERNAL, e.what());
    }

    std::ifstream output_file(output_path, std::ios::binary | std::ios::ate);
    if (!output_file.is_open()) {
        return grpc::Status(grpc::StatusCode::INTERNAL, "Failed to read output file");
    }
    std::streamsize size = output_file.tellg();
    output_file.seekg(0, std::ios::beg);
    std::string buffer(size, '\0');
    output_file.read(buffer.data(), size);

    response->set_wav_data(std::move(buffer));

    std::remove(output_path.c_str());

    return grpc::Status::OK;
}

void RunServer(const std::string& address) {
    AudioProcessorServiceImpl service;

    grpc::ServerBuilder builder;
    builder.AddListeningPort(address, grpc::InsecureServerCredentials());
    builder.RegisterService(&service);

    std::unique_ptr<grpc::Server> server(builder.BuildAndStart());
    std::cout << "AudioProcessorService listening on " << address << '\n';

    server->Wait();
}

} // namespace audioproc
