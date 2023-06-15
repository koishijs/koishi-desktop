#include "koishell/util/strings.hpp"

namespace KoiShell {

char *WideCharToUTF8(_In_ wchar_t *w) {
  const int len = WideCharToMultiByte(
      CP_UTF8, WC_ERR_INVALID_CHARS, w, -1, nullptr, 0, nullptr, nullptr);
  if (!len) return nullptr;
  char *s = new char[len + 1];
  s[len] = 0;
  WideCharToMultiByte(
      CP_UTF8, WC_ERR_INVALID_CHARS, w, -1, s, len, nullptr, nullptr);
  return s;
}

wchar_t *UTF8ToWideChar(_In_ const char *s) {
  const int len =
      MultiByteToWideChar(CP_UTF8, MB_ERR_INVALID_CHARS, s, -1, nullptr, 0);
  if (!len) return nullptr;
  wchar_t *w = new wchar_t[len + 1];
  w[len] = 0;
  MultiByteToWideChar(CP_UTF8, MB_ERR_INVALID_CHARS, s, -1, w, len);
  return w;
}

} // namespace KoiShell
