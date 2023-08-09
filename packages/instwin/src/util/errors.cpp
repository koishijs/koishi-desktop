#include "instwin/util/errors.hpp"

namespace InstWin {

void LogAndFailWithLastError(HWND hWnd, const wchar_t *message) {
  wchar_t *errMessage;
  unsigned long err = GetLastError();
  FormatMessageW(
      FORMAT_MESSAGE_ALLOCATE_BUFFER | FORMAT_MESSAGE_FROM_SYSTEM |
          FORMAT_MESSAGE_IGNORE_INSERTS,
      nullptr,
      err,
      0,
      reinterpret_cast<wchar_t *>(&errMessage),
      0,
      nullptr);
  auto errString =
      std::format(L"{}\nWinError {}: {}", message, err, errMessage);
  MessageBoxW(hWnd, errString.c_str(), L"Error", MB_ICONERROR);
  ExitProcess(1);
}

} // namespace InstWin
