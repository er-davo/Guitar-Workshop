FROM conanio/gcc12-ubuntu16.04 AS builder

WORKDIR /audio-preprocessor

COPY app/ .

RUN conan profile detect && \
    conan install . --output-folder=build --build=missing

RUN cmake -B build -DCMAKE_BUILD_TYPE=Release && \
    cmake --build build --target audio-preprocessor -j$(nproc)

FROM alpine

WORKDIR /app

RUN apk add --no-cache libstdc++

COPY --from=builder /audio-preprocessor/build/audio-preprocessor .

RUN mkdir temp

CMD [ "./audio-preprocessor" ]