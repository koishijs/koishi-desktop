#ifndef _KOISHELL_MODE_WEBVIEW_
#define _KOISHELL_MODE_WEBVIEW_

#include <wil/com.h>
#include <windows.h>
#include <wrl.h>

#include "WebView2.h"
#include "nlohmann/json.hpp"

#include "koishell/util/logger.hpp"

using json = nlohmann::json;

namespace KoiShell {

class WebViewWindow {
  HINSTANCE hInstance;
  int nCmdShow;
  json arg;

  WNDCLASSEXW wcex;
  wil::com_ptr<ICoreWebView2Controller> webviewController;
  wil::com_ptr<ICoreWebView2> webview;

public:
  WebViewWindow(_In_ HINSTANCE hInstance, _In_ int nCmdShow, _In_ json arg);
  int Run();
};

int RunWebView(_In_ HINSTANCE hInstance, _In_ int nCmdShow, _In_ json arg);

LRESULT CALLBACK WndProcWebView(
    _In_ HWND hWnd, _In_ UINT message, _In_ WPARAM wParam, _In_ LPARAM lParam);

} // namespace KoiShell

#endif /* _KOISHELL_MODE_WEBVIEW_ */
