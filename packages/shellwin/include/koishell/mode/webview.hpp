#ifndef _KOISHELL_MODE_WEBVIEW_
#define _KOISHELL_MODE_WEBVIEW_

#include <sstream>
#include <wil/com.h>
#include <windows.h>
#include <wrl.h>

#include "WebView2.h"
#include "nlohmann/json.hpp"

#include "koishell/util/logger.hpp"
#include "koishell/util/strings.hpp"

using njson = nlohmann::json;

namespace KoiShell {

const wchar_t *const KoiShellWebViewClass = L"KoiShellWebViewClass";
const wchar_t *const KoiShellWebViewTitleSuffix = L" - Koishi";

class WebViewWindow {
  HINSTANCE hInstance;
  int nCmdShow;
  njson arg;

  WNDCLASSEXW wcex;
  HWND hWnd;
  wil::com_ptr<ICoreWebView2Controller> webviewController;
  wil::com_ptr<ICoreWebView2> webview;

  static LRESULT CALLBACK WndProc(
      _In_ HWND hWnd,
      _In_ UINT message,
      _In_ WPARAM wParam,
      _In_ LPARAM lParam);

public:
  WebViewWindow(_In_ HINSTANCE hInstance, _In_ int nCmdShow, _In_ njson arg);
  int Run();
};

int RunWebView(_In_ HINSTANCE hInstance, _In_ int nCmdShow, _In_ njson arg);

} // namespace KoiShell

#endif /* _KOISHELL_MODE_WEBVIEW_ */
