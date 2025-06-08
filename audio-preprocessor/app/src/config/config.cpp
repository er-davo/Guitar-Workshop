#include "config/config.h"

namespace audioproc {

Config loadConfig() {
    Config cfg;

    cfg.port = std::getenv("PORT");

    return cfg;
}

} // namespace audioproc
