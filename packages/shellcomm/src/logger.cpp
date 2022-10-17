#include <iostream>

#include "shellcomm/logger.hpp"

namespace shellcomm {

void log(const char *messages) {
  std::cerr << messages << std::endl;
}

} // namespace shellcomm
