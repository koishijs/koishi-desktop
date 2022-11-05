#include "koishell/main.hpp"

int WINAPI wWinMain(
    _In_ HINSTANCE hInstance,
    _In_opt_ HINSTANCE hPrevInstance,
    _In_ PWSTR pCmdLine,
    _In_ int nCmdShow) {
  if (!AttachConsole(ATTACH_PARENT_PROCESS)) return 1;
  if (!SetConsoleCP(CP_UTF8)) return 1;
  if (!SetConsoleOutputCP(CP_UTF8)) return 1;

  int argc;
  wchar_t **argv = CommandLineToArgvW(GetCommandLineW(), &argc);

  if (argc != 2) KoiShell::LogAndFail(L"argc not valid.");

  char *rawArg = KoiShell::WideCharToUTF8(argv[1]);
  if (!rawArg)
    KoiShell::LogAndFailWithLastError(L"Failed to convert argv to string.");

  ShellComm::ParseResult parseResult;
  if (!ShellComm::Parse(rawArg, &parseResult)) return 1;
  delete rawArg;

  switch (parseResult.mode) {
  case ShellComm::MODE_WEBVIEW:
    return KoiShell::RunWebView(hInstance, nCmdShow, parseResult.json);
  default:
    KoiShell::LogAndFail(L"Unknown mode.");
  }
}
