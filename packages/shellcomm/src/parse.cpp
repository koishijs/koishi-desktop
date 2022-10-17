#include "shellcomm/parse.hpp"
#include "shellcomm/logger.hpp"

using json = nlohmann::json;

namespace shellcomm {
bool parse(int argc, const char **argv, parse_result *result) {
  if (argc != 2) {
    log("argc not valid.");
    return false;
  }

  int b64_len = strlen(argv[1]) + 1;
  char *raw_json = new char[b64_len];
  int cnt = 0;
  base64::base64_decodestate s;

  base64::base64_init_decodestate(&s);
  cnt = base64::base64_decode_block(argv[1], b64_len, raw_json, &s);
  *(raw_json + cnt) = 0;

  result->json = json::parse(raw_json, nullptr, false);
  if (result->json.is_discarded()) {
    log("Failed to parse arg.");
    return false;
  }

  std::string mode = result->json["mode"];
  if (mode == "webview") {
    result->mode = mode::MODE_WEBVIEW;
  } else {
    result->mode = mode::MODE_UNKNOWN;
  }

  return true;
}
} // namespace shellcomm
