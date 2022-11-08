#ifndef _SHELLCOMM_PARSE_
#define _SHELLCOMM_PARSE_

#include "b64/decode.h"
#include "b64/encode.h"
#include "nlohmann/json.hpp"

using njson = nlohmann::json;

namespace ShellComm {

enum Mode {
  MODE_UNKNOWN = 0,
  MODE_WEBVIEW = 1,
  MODE_DIALOG = 2
};

struct ParseResult {
  Mode mode;
  njson json;
};

bool ParseArgv(int argc, const char **argv, ParseResult *result);

bool Parse(const char *arg, ParseResult *result);

void SetOutput(std::string str);

} // namespace ShellComm

#endif /* _SHELLCOMM_PARSE_ */
