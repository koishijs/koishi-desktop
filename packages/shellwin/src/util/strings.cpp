#include "koishell/util/strings.hpp"

namespace KoiShell {

char *WideCharToUTF8(wchar_t *w) {
  int len = WideCharToMultiByte(
      CP_UTF8, WC_ERR_INVALID_CHARS, w, -1, nullptr, 0, nullptr, nullptr);
  char *result = new char[len + 1];
  result[len] = 0;
  WideCharToMultiByte(
      CP_UTF8, WC_ERR_INVALID_CHARS, w, -1, result, len, nullptr, nullptr);
  return result;
}

} // namespace KoiShell
