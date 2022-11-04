#include "shellcomm/parse.hpp"
#include "shellcomm/logger.hpp"

using njson = nlohmann::json;

namespace ShellComm {

bool ParseArgv(int argc, const char **argv, ParseResult *result) {
  if (argc != 2) {
    Log("argc not valid.");
    return false;
  }

  return Parse(argv[1], result);
}

bool Parse(const char *arg, ParseResult *result) {
  int b64_len = strlen(arg) + 1;
  char *raw_json = new char[b64_len];
  int cnt = 0;
  base64::base64_decodestate s;

  base64::base64_init_decodestate(&s);
  cnt = base64::base64_decode_block(arg, b64_len, raw_json, &s);
  *(raw_json + cnt) = 0;

  result->json = njson::parse(raw_json, nullptr, false);
  if (result->json.is_discarded()) {
    Log("Failed to parse arg.");
    return false;
  }

  std::string mode = result->json["mode"];
  if (mode == "webview") {
    result->mode = Mode::MODE_WEBVIEW;
  } else {
    result->mode = Mode::MODE_UNKNOWN;
  }

  return true;
}

} // namespace ShellComm
