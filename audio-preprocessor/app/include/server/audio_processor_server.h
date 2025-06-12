#pragma once

#include "processor.grpc.pb.h"
#include "core/processor/AudioProcessor.h"

#include <grpcpp/grpcpp.h>

namespace audioproc {

std::string sanitizeFilename(const std::string file_name);

class AudioProcessorServiceImpl final : public AudioProcessorService::Service {
public:
    grpc::Status ProcessAudio(grpc::ServerContext* context,
                               const ProcessAudioRequest* request,
                               ProcessAudioResponse* response) override;
    
    grpc::Status SplitIntoChunks(grpc::ServerContext* context,
                               const SplitAudioRequest* request,
                               SplitAudioResponse* response) override;

};

void RunServer(const std::string& address);

} // namespace audioproc
