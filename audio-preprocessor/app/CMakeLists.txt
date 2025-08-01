cmake_minimum_required(VERSION 3.15)
project(audio_preprocessor)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

#add_compile_definitions(DEBUG)

find_package(SndFile REQUIRED)
find_package(soxr REQUIRED)
find_package(Protobuf REQUIRED)
find_package(gRPC REQUIRED)
find_package(spdlog REQUIRED)

file(GLOB_RECURSE SOURCES CONFIGURE_DEPENDS src/*.cpp)

list(APPEND SOURCES
    proto/processor.pb.cc
    proto/processor.grpc.pb.cc
)

add_executable(audio-preprocessor ${SOURCES})

target_include_directories(audio-preprocessor 
    PRIVATE 
    include
    proto
)

target_link_libraries(audio-preprocessor
    PRIVATE
    SndFile::sndfile
    soxr::soxr
    protobuf::libprotobuf
    gRPC::grpc++
    spdlog::spdlog
)