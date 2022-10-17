#include "b64/decode.h"

#include "shellcomm/logger.hpp"
#include "shellcomm/parse.hpp"

namespace shellcomm {
bool parse(int argc, const char **argv, parse_result *result) {
  if (argc != 2) {
    log("argc not valid.");
    return false;
  }

  int b64_len = strlen(argv[1]) + 1;
  char *json = new char[b64_len];
  int cnt = 0;
  base64::base64_decodestate s;

  base64::base64_init_decodestate(&s);
  cnt = base64::base64_decode_block(argv[1], b64_len, json, &s);
  *(json + cnt) = 0;

  return true;
}
} // namespace shellcomm
