#include "koishell/main.hpp"

int WINAPI wWinMain(
    _In_ HINSTANCE hInstance,
    _In_opt_ HINSTANCE hPrevInstance,
    _In_ PWSTR pCmdLine,
    _In_ int nCmdShow) {
  if (!AttachConsole(ATTACH_PARENT_PROCESS)) return 1;

  int argc;
  wchar_t **argv = CommandLineToArgvW(GetCommandLineW(), &argc);

  if (argc != 2) {
    ShellComm::Log("argc not valid.");
    return 1;
  }

  char *arg = KoiShell::WideCharToUTF8(argv[1]);

  ShellComm::ParseResult parseResult;
  if (!ShellComm::Parse(arg, &parseResult)) return 1;
  delete arg;

  switch (parseResult.mode) {
  case ShellComm::MODE_WEBVIEW:
    break;
  default:
    ShellComm::Log("Unknown mode.");
    return 1;
  }

  return 0;
}
