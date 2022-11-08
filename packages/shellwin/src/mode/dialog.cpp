#include "koishell/mode/dialog.hpp"

namespace KoiShell {

int RunDialog(_In_ HINSTANCE hInstance, _In_ njson arg) {
  std::string titleS = arg["title"];
  wchar_t *title = KoiShell::UTF8ToWideChar(const_cast<char *>(titleS.c_str()));
  if (!title) LogAndFailWithLastError(L"Failed to parse title.");

  std::string text1S = arg["text1"];
  wchar_t *text1 = KoiShell::UTF8ToWideChar(const_cast<char *>(text1S.c_str()));
  if (!text1) LogAndFailWithLastError(L"Failed to parse text1.");

  std::string text2S = arg["text2"];
  wchar_t *text2 = KoiShell::UTF8ToWideChar(const_cast<char *>(text2S.c_str()));
  if (!text2) LogAndFailWithLastError(L"Failed to parse text2.");

  unsigned int buttonCount = arg["buttonCount"];
  TASKDIALOG_BUTTON *buttons = new TASKDIALOG_BUTTON[buttonCount];

  if (buttonCount >= 1) {
    std::string textS = arg.value("button1Text", "OK");
    wchar_t *text = KoiShell::UTF8ToWideChar(const_cast<char *>(textS.c_str()));
    if (!text) LogAndFailWithLastError(L"Failed to parse button1Text.");
    buttons[0].nButtonID = 10;
    buttons[0].pszButtonText = text;
  }

  if (buttonCount >= 2) {
    std::string textS = arg.value("button2Text", "Cancel");
    wchar_t *text = KoiShell::UTF8ToWideChar(const_cast<char *>(textS.c_str()));
    if (!text) LogAndFailWithLastError(L"Failed to parse button2Text.");
    buttons[0].nButtonID = 11;
    buttons[0].pszButtonText = text;
  }

  if (buttonCount >= 3) {
    std::string textS = arg.value("button3Text", "Don't Save");
    wchar_t *text = KoiShell::UTF8ToWideChar(const_cast<char *>(textS.c_str()));
    if (!text) LogAndFailWithLastError(L"Failed to parse button3Text.");
    buttons[0].nButtonID = 12;
    buttons[0].pszButtonText = text;
  }

  wchar_t *icon;
  std::string style = arg["style"];
  if (style == "info")
    icon = TD_INFORMATION_ICON;
  else if (style == "warn")
    icon = TD_WARNING_ICON;
  else if (style == "error")
    icon = TD_ERROR_ICON;

  TASKDIALOGCONFIG dlg;
  dlg.cbSize = sizeof(dlg);
  dlg.hwndParent = nullptr;
  dlg.hInstance = hInstance;
  dlg.dwFlags = TDF_CAN_BE_MINIMIZED | TDF_USE_COMMAND_LINKS;
  dlg.dwCommonButtons = 0;
  dlg.pszWindowTitle = title;
  dlg.pszMainIcon = icon;
  dlg.pszMainInstruction = text1;
  dlg.pszContent = text2;
  dlg.cButtons = buttonCount;
  dlg.pButtons = buttons;
  dlg.nDefaultButton = 10;
  dlg.cRadioButtons = 0;
  dlg.pRadioButtons = nullptr;
  dlg.nDefaultRadioButton = 0;
  dlg.pszVerificationText = nullptr;
  dlg.pszExpandedInformation = nullptr;
  dlg.pszExpandedControlText = nullptr;
  dlg.pszCollapsedControlText = nullptr;
  dlg.hFooterIcon = nullptr;
  dlg.pszFooter = nullptr;
  dlg.pfCallback = nullptr;
  dlg.lpCallbackData = 0;
  dlg.cxWidth = 0;

  int result = 0;
  long success = TaskDialogIndirect(&dlg, &result, nullptr, nullptr);
  result -= 9;
  if (success) {
    std::wstringstream s;
    s << L"[WinError " << success << "] " << L"Failed to process task dialog.";
    LogAndFail(s.str());
  }

  std::stringstream s;
  s << "{\"result\":" << ((result > 0 && result <= buttonCount) ? result : 2)
    << '}';
  ShellComm::SetOutput(s.str());

  return 0;
}

} // namespace KoiShell
