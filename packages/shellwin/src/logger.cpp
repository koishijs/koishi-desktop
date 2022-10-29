#include "koishell/logger.hpp"

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
      message,
      0,
      nullptr);

  std::wcerr << "[WinError " << err << "] " << message << std::endl;

  LocalFree(message);
}

void LogWithLastError(const wchar_t *messages) {
  std::wcerr << messages << L" Last error:" << std::endl;
  LogLastError();
}

} // namespace KoiShell
