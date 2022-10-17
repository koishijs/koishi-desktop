#include "koishell/main.hpp"

int WINAPI wWinMain(
    HINSTANCE hInstance,
    HINSTANCE hPrevInstance,
    PWSTR pCmdLine,
    int nCmdShow) {
  int argc;
  wchar_t **argv = CommandLineToArgvW(GetCommandLineW(), &argc);

  if (argc != 2) {
    shellcomm::log("argc not valid.");
    return 1;
  }
}
