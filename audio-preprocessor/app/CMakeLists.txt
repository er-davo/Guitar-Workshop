cmake_minimum_required(VERSION 3.15)
project(audio_preprocessor)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

#add_compile_definitions(DEBUG)

# подключаем зависимости через CMakeDeps
find_package(SndFile REQUIRED)
find_package(soxr REQUIRED)
find_package(Protobuf REQUIRED)
find_package(gRPC REQUIRED)

# исходники
file(GLOB_RECURSE SOURCES CONFIGURE_DEPENDS src/*.cpp)

list(APPEND SOURCES
    proto/audio_processor.pb.cc
    proto/audio_processor.grpc.pb.cc
)

add_executable(audio-preprocessor ${SOURCES})

# инклюды
target_include_directories(audio-preprocessor 
    PRIVATE 
    include
    proto
)

# линковка
target_link_libraries(audio-preprocessor
    PRIVATE
    SndFile::sndfile
    soxr::soxr
    protobuf::libprotobuf
    gRPC::grpc++
)