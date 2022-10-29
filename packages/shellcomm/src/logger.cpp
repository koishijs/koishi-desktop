#include <iostream>

#include "shellcomm/logger.hpp"

namespace ShellComm {

void Log(const char *messages) {
  std::cerr << messages << std::endl;
}

} // namespace ShellComm
