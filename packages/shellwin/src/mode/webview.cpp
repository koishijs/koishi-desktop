#include "koishell/mode/webview.hpp"

namespace KoiShell {

WebViewWindow::WebViewWindow(
    _In_ HINSTANCE hInstance, _In_ int nCmdShow, _In_ njson arg)
    : hInstance(hInstance), nCmdShow(nCmdShow), arg(arg) {
}

int WebViewWindow::Run() {
  std::string nameS = arg["name"];
  wchar_t *nameC = KoiShell::UTF8ToWideChar(const_cast<char *>(nameS.c_str()));
  if (!nameC) {
    LogWithLastError(L"Failed to parse nameC.");
    return 1;
  }
  std::wostringstream nameStream;
  nameStream << nameC << KoiShellWebViewTitleSuffix;
  std::wstring name = nameStream.str();

  std::string urlS = arg["url"];
  wchar_t *url = KoiShell::UTF8ToWideChar(const_cast<char *>(urlS.c_str()));
  if (!url) {
    LogWithLastError(L"Failed to parse url.");
    return 1;
  }

  wcex.cbSize = sizeof(WNDCLASSEXW);
  wcex.style = CS_HREDRAW | CS_VREDRAW;
  wcex.lpfnWndProc = WndProc;
  wcex.cbClsExtra = 0;
  wcex.cbWndExtra = 0;
  wcex.hInstance = hInstance;
  wcex.hIcon = LoadIconW(hInstance, IDI_APPLICATION);
  wcex.hCursor = LoadCursorW(hInstance, IDC_ARROW);
  wcex.hbrBackground = (HBRUSH)(COLOR_WINDOW + 1);
  wcex.lpszMenuName = nullptr;
  wcex.lpszClassName = KoiShellWebViewClass;
  wcex.hIconSm = LoadIconW(hInstance, IDI_APPLICATION);

  if (!RegisterClassExW(&wcex)) {
    LogWithLastError(L"Failed to register window class.");
    return 1;
  }

  hWnd = CreateWindowExW(
      WS_EX_OVERLAPPEDWINDOW,
      KoiShellWebViewClass,
      name.c_str(),
      WS_OVERLAPPEDWINDOW,
      CW_USEDEFAULT,
      CW_USEDEFAULT,
      1366,
      768,
      nullptr,
      nullptr,
      hInstance,
      this);

  if (!hWnd) {
    LogWithLastError(L"Failed to create window.");
    return 1;
  }

  ShowWindow(hWnd, nCmdShow);
  UpdateWindow(hWnd);

  CreateCoreWebView2Environment(
      Microsoft::WRL::Callback<
          ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler>(
          [this,
           &url](HRESULT result, ICoreWebView2Environment *env) -> HRESULT {
            env->CreateCoreWebView2Controller(
                hWnd,
                Microsoft::WRL::Callback<
                    ICoreWebView2CreateCoreWebView2ControllerCompletedHandler>(
                    [this, &url](
                        HRESULT result,
                        ICoreWebView2Controller *controller) -> HRESULT {
                      if (controller) {
                        webviewController = controller;
                        webviewController->get_CoreWebView2(&webview);
                      }

                      wil::com_ptr<ICoreWebView2Settings> settings;
                      webview->get_Settings(&settings);
                      settings->put_IsScriptEnabled(1);
                      settings->put_AreDefaultScriptDialogsEnabled(1);
                      settings->put_IsWebMessageEnabled(1);
                      settings->put_AreDevToolsEnabled(1);

                      RECT bounds;
                      GetClientRect(hWnd, &bounds);
                      webviewController->put_Bounds(bounds);

                      webview->Navigate(url);

                      return 0;
                    })
                    .Get());

            return 0;
          })
          .Get());

  MSG msg;
  while (GetMessageW(&msg, nullptr, 0, 0)) {
    TranslateMessage(&msg);
    DispatchMessageW(&msg);
  }

  return (int)msg.wParam;
}

LRESULT CALLBACK WebViewWindow::WndProc(
    _In_ HWND hWnd, _In_ UINT message, _In_ WPARAM wParam, _In_ LPARAM lParam) {
  WebViewWindow *pThis;

  if (message == WM_NCCREATE) {
    pThis = static_cast<WebViewWindow *>(
        reinterpret_cast<CREATESTRUCTW *>(lParam)->lpCreateParams);

    SetLastError(0);
    if (!SetWindowLongPtrW(
            hWnd, GWLP_USERDATA, reinterpret_cast<long long>(pThis)))
      if (GetLastError() != 0) {
        LogWithLastError(L"Failed to set window user data.");
        return 0;
      }
  } else
    pThis = reinterpret_cast<WebViewWindow *>(
        GetWindowLongPtrW(hWnd, GWLP_USERDATA));

  if (!pThis) return DefWindowProcW(hWnd, message, wParam, lParam);

  switch (message) {
  case WM_DESTROY:
    PostQuitMessage(0);
    return 0;

  case WM_SIZE:
    if (pThis->webviewController) {
      RECT bounds;
      GetClientRect(hWnd, &bounds);
      pThis->webviewController->put_Bounds(bounds);
    }

  default:
    return DefWindowProcW(hWnd, message, wParam, lParam);
  }
}

int RunWebView(_In_ HINSTANCE hInstance, _In_ int nCmdShow, _In_ njson arg) {
  WebViewWindow webviewWindow = WebViewWindow(hInstance, nCmdShow, arg);
  return webviewWindow.Run();
}

} // namespace KoiShell