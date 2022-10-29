#include "koishell/main.hpp"

int WINAPI wWinMain(
    HINSTANCE hInstance,
    HINSTANCE hPrevInstance,
    PWSTR pCmdLine,
    int nCmdShow) {
  if (!AttachConsole(ATTACH_PARENT_PROCESS)) return 1;

  int argc;
  wchar_t **argv = CommandLineToArgvW(GetCommandLineW(), &argc);

  if (argc != 2) {
    ShellComm::Log("argc not valid.");
    return 1;
  }
}
