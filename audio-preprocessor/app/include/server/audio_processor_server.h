#pragma once

#include "audio_processor.grpc.pb.h"
#include "processor/processor.h"

#include <grpcpp/grpcpp.h>

namespace audioproc {

std::string sanitizeFilename(const std::string file_name);

class AudioProcessorServiceImpl final : public AudioProcessorService::Service {
public:
    grpc::Status ProcessAudio(grpc::ServerContext* context,
                               const ProcessAudioRequest* request,
                               ProcessAudioResponse* response) override;

};

void RunServer(const std::string& address);

} // namespace audioproc
