# CMakeList.txt: проект CMake для fx-processor; включите исходный код и определения,
# укажите здесь логику для конкретного проекта.
#

add_subdirectory(juce)
# include_directories(juce/modules)
set(JUCE_MODULES
    juce_audio_utils
    juce_dsp
)

find_package(gRPC REQUIRED)
find_package(Protobuf REQUIRED)

add_executable (fx-processor
    "src/main.cpp"
)

target_link_libraries(fx-processor
    PRIVATE
    $(JUCE_MODULES)
    grpc++
    grpc
    protobuf
)

if (CMAKE_VERSION VERSION_GREATER 3.12)
  set_property(TARGET fx-processor PROPERTY CXX_STANDARD 20)
endif()

# TODO: Добавьте тесты и целевые объекты, если это необходимо.
