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
  } else if (mode == "dialog") {
    result->mode = Mode::MODE_DIALOG;
  } else {
    result->mode = Mode::MODE_UNKNOWN;
  }

  return true;
}

void SetOutput(std::string str) {
  const char *input = str.c_str();
  unsigned long long len = str.length();
  char *output = new char[len * 4];
  char *c = output;
  int cnt = 0;
  base64::base64_encodestate s;

  base64_init_encodestate(&s);
  cnt = base64::base64_encode_block(input, len, c, &s);
  c += cnt;
  cnt = base64_encode_blockend(c, &s);
  c += cnt;

  *c = 0;

  std::cout << output << std::endl;

  delete[] output;
}

} // namespace ShellComm
