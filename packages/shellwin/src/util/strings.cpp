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

wchar_t *UTF8ToWideChar(_In_ char *s) {
  int len =
      MultiByteToWideChar(CP_UTF8, MB_ERR_INVALID_CHARS, s, -1, nullptr, 0);
  if (!len) return nullptr;
  wchar_t *w = new wchar_t[len + 1];
  w[len] = 0;
  MultiByteToWideChar(CP_UTF8, MB_ERR_INVALID_CHARS, s, -1, result, len);
  return w;
}

} // namespace KoiShell
