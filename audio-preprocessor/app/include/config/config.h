#pragma once

#include <string>

namespace audioproc {

struct Config {
    std::string port;
};

Config loadConfig();

} // namespace audioproc