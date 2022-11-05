#include "koishell/util/logger.hpp"

namespace KoiShell {

void LogW(const wchar_t *messages) {
  std::wcerr << messages << std::endl;
}

void LogLastError() {
  wchar_t *message;

  unsigned long err = GetLastError();

  FormatMessageW(
      FORMAT_MESSAGE_ALLOCATE_BUFFER | FORMAT_MESSAGE_FROM_SYSTEM |
          FORMAT_MESSAGE_IGNORE_INSERTS,
      nullptr,
      err,
      0,
      (wchar_t *)&message,
      0,
      nullptr);

  std::wcerr << "[WinError " << err << "] " << message << std::endl;

  LocalFree(message);
}

void LogWithLastError(const wchar_t *messages) {
  std::wcerr << messages << L" Last error:" << std::endl;
  LogLastError();
}

// https://github.com/MicrosoftEdge/WebView2Samples/blob/main/SampleApps/WebView2APISample/CheckFailure.cpp
void CheckFailure(HRESULT hr, const std::wstring &message) {
  if (FAILED(hr)) {
    std::wcerr << "[WinError HRESULT " << hr << "] " << message << std::endl;
    FAIL_FAST();
  }
}

} // namespace KoiShell
