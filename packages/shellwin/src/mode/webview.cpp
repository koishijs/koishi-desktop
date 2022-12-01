#include "koishell/mode/webview.hpp"

namespace KoiShell {

WebViewWindow::WebViewWindow(
    _In_ HINSTANCE hInstance, _In_ int nCmdShow, _In_ njson arg)
    : hInstance(hInstance), nCmdShow(nCmdShow), arg(arg) {
}

int WebViewWindow::Run() {
  wchar_t cwd[MAX_PATH];
  if (!GetCurrentDirectoryW(MAX_PATH, cwd))
    LogAndFailWithLastError(L"Failed to get cwd.");
  wchar_t udf[MAX_PATH];
  if (!PathCombineW(udf, cwd, L"data\\home\\WebView2"))
    LogAndFailWithLastError(L"Failed to combine udf.");
  int udfErr = SHCreateDirectoryExW(nullptr, udf, nullptr);
  if (udfErr != ERROR_SUCCESS && udfErr != ERROR_FILE_EXISTS &&
      udfErr != ERROR_ALREADY_EXISTS) {
    std::wstringstream s;
    s << L"[WinError " << udfErr << "] "
      << L"Failed to create directory for udf.";
    LogAndFail(s.str());
  }

  std::string nameS = arg["name"];
  wchar_t *nameC = KoiShell::UTF8ToWideChar(const_cast<char *>(nameS.c_str()));
  if (!nameC) LogAndFailWithLastError(L"Failed to parse nameC.");
  std::wostringstream nameStream;
  nameStream << nameC << KoiShellWebViewTitleSuffix;
  std::wstring name = nameStream.str();

  std::string urlS = arg["url"];
  wchar_t *url = KoiShell::UTF8ToWideChar(const_cast<char *>(urlS.c_str()));
  if (!url) LogAndFailWithLastError(L"Failed to parse url.");

  HRSRC userscriptRc =
      FindResourceW(hInstance, MAKEINTRESOURCEW(102), MAKEINTRESOURCEW(256));
  HGLOBAL userscriptRcData = LoadResource(hInstance, userscriptRc);
  unsigned long userscriptSize = SizeofResource(hInstance, userscriptRc);
  const char *userscriptData =
      static_cast<const char *>(LockResource(userscriptRcData));
  char *userscriptS = new char[userscriptSize + 1];
  memcpy_s(userscriptS, userscriptSize, userscriptData, userscriptSize);
  userscriptS[userscriptSize] = 0;
  wchar_t *userscript = KoiShell::UTF8ToWideChar(userscriptS);
  delete[] userscriptS;

  wcex.cbSize = sizeof(WNDCLASSEXW);
  wcex.style = CS_HREDRAW | CS_VREDRAW;
  wcex.lpfnWndProc = WndProc;
  wcex.cbClsExtra = 0;
  wcex.cbWndExtra = 0;
  wcex.hInstance = hInstance;
  wcex.hIcon = static_cast<HICON>(
      LoadImageW(hInstance, MAKEINTRESOURCEW(101), IMAGE_ICON, 0, 0, 0));
  wcex.hCursor = LoadCursorW(hInstance, IDC_ARROW);
  // wcex.hbrBackground = static_cast<HBRUSH>(GetStockObject(NULL_BRUSH));
  wcex.hbrBackground = static_cast<HBRUSH>(GetStockObject(BLACK_BRUSH));
  wcex.lpszMenuName = nullptr;
  wcex.lpszClassName = KoiShellWebViewClass;
  wcex.hIconSm = static_cast<HICON>(
      LoadImageW(hInstance, MAKEINTRESOURCEW(101), IMAGE_ICON, 0, 0, 0));

  if (!RegisterClassExW(&wcex))
    LogAndFailWithLastError(L"Failed to register window class.");

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

  if (!hWnd) LogAndFailWithLastError(L"Failed to create window.");

  ShowWindow(hWnd, nCmdShow);
  UpdateWindow(hWnd);

  MARGINS dwmExtendFrameIntoClientAreaMargins = {-1};
  DwmExtendFrameIntoClientArea(hWnd, &dwmExtendFrameIntoClientAreaMargins);
  int dwmUseDarkMode = 0;
  DwmSetWindowAttribute(
      hWnd,
      DWMWINDOWATTRIBUTE::DWMWA_USE_IMMERSIVE_DARK_MODE,
      &dwmUseDarkMode,
      sizeof(dwmUseDarkMode));
  unsigned int dwmCornerPreference = DWM_WINDOW_CORNER_PREFERENCE::DWMWCP_ROUND;
  DwmSetWindowAttribute(
      hWnd,
      DWMWINDOWATTRIBUTE::DWMWA_WINDOW_CORNER_PREFERENCE,
      &dwmCornerPreference,
      sizeof(dwmCornerPreference));
  unsigned int dwmSystemBackdropType =
      2; // DWM_SYSTEMBACKDROP_TYPE::DWMSBT_MAINWINDOW
  DwmSetWindowAttribute(
      hWnd,
      38, // DWMWINDOWATTRIBUTE::DWMWA_SYSTEMBACKDROP_TYPE
      &dwmSystemBackdropType,
      sizeof(dwmSystemBackdropType));
  dwmSystemBackdropType = 4; // DWM_SYSTEMBACKDROP_TYPE::DWMSBT_TABBEDWINDOW
  DwmSetWindowAttribute(
      hWnd,
      38, // DWMWINDOWATTRIBUTE::DWMWA_SYSTEMBACKDROP_TYPE
      &dwmSystemBackdropType,
      sizeof(dwmSystemBackdropType));

  auto options = Microsoft::WRL::Make<CoreWebView2EnvironmentOptions>();
  CheckFailure(
      CreateCoreWebView2EnvironmentWithOptions(
          nullptr,
          udf,
          options.Get(),
          Microsoft::WRL::Callback<
              ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler>(
              [this, &url, &userscript](
                  HRESULT result, ICoreWebView2Environment *env) -> HRESULT {
                CheckFailure(
                    env->CreateCoreWebView2Controller(
                        hWnd,
                        Microsoft::WRL::Callback<
                            ICoreWebView2CreateCoreWebView2ControllerCompletedHandler>(
                            [this, &url, &userscript](
                                HRESULT result,
                                ICoreWebView2Controller *controller)
                                -> HRESULT {
                              if (controller) {
                                webviewController = controller;
                                CheckFailure(
                                    webviewController->get_CoreWebView2(
                                        &webview),
                                    L"Failed to get WebView2.");
                              }

                              wil::com_ptr<ICoreWebView2Controller2>
                                  controller2 =
                                      webviewController
                                          .query<ICoreWebView2Controller2>();
                              CheckFailure(
                                  controller2->put_DefaultBackgroundColor({}),
                                  L"Failed to set transparent background.");

                              wil::com_ptr<ICoreWebView2Settings> settings;
                              webview->get_Settings(&settings);
                              settings->put_IsScriptEnabled(1);
                              settings->put_AreDefaultScriptDialogsEnabled(1);
                              settings->put_IsWebMessageEnabled(1);
                              settings->put_AreDevToolsEnabled(1);

                              webview->AddScriptToExecuteOnDocumentCreated(
                                  userscript, nullptr);
                              delete[] userscript;

                              EventRegistrationToken eventRegistrationToken;
                              CheckFailure(
                                  webview->add_WebMessageReceived(
                                      Microsoft::WRL::Callback<
                                          ICoreWebView2WebMessageReceivedEventHandler>(
                                          [this](
                                              ICoreWebView2 * /*sender*/,
                                              ICoreWebView2WebMessageReceivedEventArgs
                                                  *args) {
                                            // Check URI
                                            wil::unique_cotaskmem_string uriRaw;
                                            CheckFailure(
                                                args->get_Source(&uriRaw),
                                                L"Failed to get webview URI.");
                                            std::wstring uri = uriRaw.get();

                                            wil::unique_cotaskmem_string
                                                messageRaw;
                                            CheckFailure(
                                                args->TryGetWebMessageAsString(
                                                    &messageRaw),
                                                L"Failed to get webview "
                                                L"message.");
                                            std::wstring message =
                                                messageRaw.get();

                                            OnMessage(&message);

                                            return S_OK;
                                          })
                                          .Get(),
                                      &eventRegistrationToken),
                                  L"Failed to register webview message "
                                  L"handler.");

                              RECT bounds;
                              GetClientRect(hWnd, &bounds);
                              webviewController->put_Bounds(bounds);

                              webview->Navigate(url);

                              return 0;
                            })
                            .Get()),
                    L"Failed to create WebView2 controller.");

                return 0;
              })
              .Get()),
      L"Failed to create WebView2 environment.");

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
      if (GetLastError() != 0)
        LogAndFailWithLastError(L"Failed to set window user data.");
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
    return 0;

  default:
    return DefWindowProcW(hWnd, message, wParam, lParam);
  }
}

void WebViewWindow::OnMessage(std::wstring *message) {
  int dwmUseDarkMode;
  if ((*message) == L"TD") {
    dwmUseDarkMode = 1;
    DwmSetWindowAttribute(
        hWnd,
        DWMWINDOWATTRIBUTE::DWMWA_USE_IMMERSIVE_DARK_MODE,
        &dwmUseDarkMode,
        sizeof(dwmUseDarkMode));
    return;
  }

  if ((*message) == L"TL" || (*message) == L"TR") {
    dwmUseDarkMode = 0;
    DwmSetWindowAttribute(
        hWnd,
        DWMWINDOWATTRIBUTE::DWMWA_USE_IMMERSIVE_DARK_MODE,
        &dwmUseDarkMode,
        sizeof(dwmUseDarkMode));
    return;
  }
}

int RunWebView(_In_ HINSTANCE hInstance, _In_ int nCmdShow, _In_ njson arg) {
  WebViewWindow webviewWindow = WebViewWindow(hInstance, nCmdShow, arg);
  return webviewWindow.Run();
}

} // namespace KoiShell
