FROM ubuntu:22.04 AS builder

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    build-essential \
    cmake \
    python3 \
    python3-pip
#    libssl-dev \
#    zlib1g-dev && \
#    rm -rf /var/lib/apt/lists/*

RUN pip3 install conan

WORKDIR /audio-preprocessor

COPY app/conanfile.txt /audio-preprocessor/

RUN conan profile detect --force

RUN conan install . --output-folder=build --build=missing

COPY app/ /audio-preprocessor/

RUN cmake -S . -B build -DCMAKE_TOOLCHAIN_FILE=build/conan_toolchain.cmake -DCMAKE_BUILD_TYPE=Release && \
    cmake --build build --target audio-preprocessor -j$(nproc)

# RUN mkdir -p temp

# CMD [ "/build/Release/audio-preprocessor" ]

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /audio-preprocessor/build/audio-preprocessor /app/

RUN mkdir -p temp

CMD [ "/app/audio-preprocessor" ]