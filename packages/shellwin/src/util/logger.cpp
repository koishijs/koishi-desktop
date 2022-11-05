#include "koishell/util/logger.hpp"

namespace KoiShell {

void FailWithLastError() {
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
  ExitProcess(1);
}

void LogAndFailWithLastError(const wchar_t *messages) {
  std::wcerr << messages << L" Last error:" << std::endl;
  FailWithLastError();
}

void LogAndFail(const wchar_t *messages) {
  std::wcerr << messages << std::endl;
  ExitProcess(1);
}

void LogAndFail(const std::wstring &messages) {
  std::wcerr << messages << std::endl;
  ExitProcess(1);
}

// https://github.com/MicrosoftEdge/WebView2Samples/blob/main/SampleApps/WebView2APISample/CheckFailure.cpp
void CheckFailure(HRESULT hr, const std::wstring &message) {
  if (FAILED(hr)) {
    std::wcerr << "[WinError HRESULT " << hr << "] " << message << std::endl;
    ExitProcess(1);
  }
}

} // namespace KoiShell
