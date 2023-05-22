#include "koishell/mode/webview.hpp"

namespace KoiShell {

WebViewWindow::WebViewWindow(
    _In_ HINSTANCE hInstance, _In_ int nCmdShow, _In_ njson arg)
    : hInstance(hInstance), nCmdShow(nCmdShow), arg(arg), supports(0) {
}

int WebViewWindow::Run() {
  OSVERSIONINFOEXW osvi;
  memset(&osvi, 0, sizeof(OSVERSIONINFOEXW));
  osvi.dwOSVersionInfoSize = sizeof(OSVERSIONINFOEXW);
  if (!GetVersionExW(reinterpret_cast<LPOSVERSIONINFOW>(&osvi)))
    LogAndFailWithLastError(L"Failed to get OS version info.");

  if (osvi.dwBuildNumber >= 17763)
    supports++; // Supports Windows 10 immersive dark mode (19), supports = 1
  if (osvi.dwBuildNumber >= 18985)
    supports++; // Supports Windows 10 immersive dark mode (20), supports = 2
  if (osvi.dwBuildNumber >= 22000)
    supports++; // Supports Windows 11 Mica, supports = 3
  if (osvi.dwBuildNumber >= 22523)
    supports++; // Supports Windows 11 Mica Tabbed, supports = 4

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
  wchar_t *userscriptW = KoiShell::UTF8ToWideChar(userscriptS);
  delete[] userscriptS;
  std::wstring userscript = std::wstring(userscriptW);
  delete[] userscriptW;
  userscript.replace(
      userscript.find(L"KOISHELL_RUNTIME_SUPPORTS"),
      25,
      supports >= 3 ? L"['enhance', 'enhanceColor']"
      : supports    ? L"['enhance']"
                    : L"[]");

  wcex.cbSize = sizeof(WNDCLASSEXW);
  wcex.style = CS_HREDRAW | CS_VREDRAW;
  wcex.lpfnWndProc = WndProc;
  wcex.cbClsExtra = 0;
  wcex.cbWndExtra = 0;
  wcex.hInstance = hInstance;
  wcex.hIcon = static_cast<HICON>(
      LoadImageW(hInstance, MAKEINTRESOURCEW(101), IMAGE_ICON, 0, 0, 0));
  wcex.hCursor = LoadCursorW(hInstance, IDC_ARROW);
  // Why?
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
                                  userscript.c_str(), nullptr);

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
  case WM_CREATE:
    if (pThis->supports) {
      MARGINS dwmExtendFrameIntoClientAreaMargins = {-1};
      DwmExtendFrameIntoClientArea(hWnd, &dwmExtendFrameIntoClientAreaMargins);
      int dwmUseDarkMode = 0;
      DwmSetWindowAttribute(
          hWnd,
          pThis->supports >= 2
              ? 20
              :   // DWMWINDOWATTRIBUTE::DWMWA_USE_IMMERSIVE_DARK_MODE
                  // = 20 (starting from 18985)
              19, // DWMWINDOWATTRIBUTE::DWMWA_USE_IMMERSIVE_DARK_MODE
                  // = 19 (before 18985)
          &dwmUseDarkMode,
          sizeof(dwmUseDarkMode));
      unsigned int dwmCornerPreference =
          2; // DWM_WINDOW_CORNER_PREFERENCE::DWMWCP_ROUND
      DwmSetWindowAttribute(
          hWnd,
          33, // DWMWINDOWATTRIBUTE::DWMWA_WINDOW_CORNER_PREFERENCE
          &dwmCornerPreference,
          sizeof(dwmCornerPreference));
      int dwmMica = 1;
      DwmSetWindowAttribute(
          hWnd,
          1029, // DWMWINDOWATTRIBUTE::DWMWA_MICA
          &dwmMica,
          sizeof(dwmMica));
      unsigned int dwmSystemBackdropType =
          pThis->supports >= 4
              ? 4
              : 2; // DWM_SYSTEMBACKDROP_TYPE::DWMSBT_MAINWINDOW
      DwmSetWindowAttribute(
          hWnd,
          38, // DWMWINDOWATTRIBUTE::DWMWA_SYSTEMBACKDROP_TYPE
          &dwmSystemBackdropType,
          sizeof(dwmSystemBackdropType));
    }

    return 0;

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
  unsigned long long const len = message->length();

  if ((*message)[0] == 'T') {
    // Sync theme
    int dwmUseDarkMode;

    switch ((*message)[1]) {
    case 'D':
      dwmUseDarkMode = 1;
      DwmSetWindowAttribute(
          hWnd,
          supports >= 2 ? 20
                        : // DWMWINDOWATTRIBUTE::DWMWA_USE_IMMERSIVE_DARK_MODE
                          // = 20 (starting from 18985)
              19,         // DWMWINDOWATTRIBUTE::DWMWA_USE_IMMERSIVE_DARK_MODE
                          // = 19 (before 18985),
          &dwmUseDarkMode,
          sizeof(dwmUseDarkMode));
      break;
    case 'L':
      dwmUseDarkMode = 0;
      DwmSetWindowAttribute(
          hWnd,
          supports >= 2 ? 20
                        : // DWMWINDOWATTRIBUTE::DWMWA_USE_IMMERSIVE_DARK_MODE
                          // = 20 (starting from 18985)
              19,         // DWMWINDOWATTRIBUTE::DWMWA_USE_IMMERSIVE_DARK_MODE
                          // = 19 (before 18985),
          &dwmUseDarkMode,
          sizeof(dwmUseDarkMode));
      break;
    case 'R':
      dwmUseDarkMode = 0;
      DwmSetWindowAttribute(
          hWnd,
          supports >= 2 ? 20
                        : // DWMWINDOWATTRIBUTE::DWMWA_USE_IMMERSIVE_DARK_MODE
                          // = 20 (starting from 18985)
              19,         // DWMWINDOWATTRIBUTE::DWMWA_USE_IMMERSIVE_DARK_MODE
                          // = 19 (before 18985),
          &dwmUseDarkMode,
          sizeof(dwmUseDarkMode));
      break;
    }

    if (len < 3) return;
    // Sync theme color
    unsigned long color;

    // Border
    color = std::stoul(
        std::wstring() + (*message)[7] + (*message)[8] + (*message)[5] +
            (*message)[6] + (*message)[3] + (*message)[4],
        nullptr,
        16);
    DwmSetWindowAttribute(
        hWnd,
        34, // DWMWINDOWATTRIBUTE::DWMWA_BORDER_COLOR
        &color,
        sizeof(color));
    // Caption
    color = std::stoul(
        std::wstring() + (*message)[13] + (*message)[14] + (*message)[11] +
            (*message)[12] + (*message)[9] + (*message)[10],
        nullptr,
        16);
    DwmSetWindowAttribute(
        hWnd,
        35, // DWMWINDOWATTRIBUTE::DWMWA_CAPTION_COLOR
        &color,
        sizeof(color));
    // Caption text
    color = std::stoul(
        std::wstring() + (*message)[19] + (*message)[20] + (*message)[17] +
            (*message)[18] + (*message)[15] + (*message)[16],
        nullptr,
        16);
    DwmSetWindowAttribute(
        hWnd,
        36, // DWMWINDOWATTRIBUTE::DWMWA_TEXT_COLOR
        &color,
        sizeof(color));
  }
}

int RunWebView(_In_ HINSTANCE hInstance, _In_ int nCmdShow, _In_ njson arg) {
  WebViewWindow webviewWindow = WebViewWindow(hInstance, nCmdShow, arg);
  return webviewWindow.Run();
}

} // namespace KoiShell
