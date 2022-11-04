#include "koishell/mode/webview.hpp"

namespace KoiShell {

const wchar_t *const KoiShellWebViewClass = L"KoiShellWebViewClass";
const wchar_t *const KoiShellWebViewTitle = L"Koishi";

WebViewWindow::WebViewWindow(
    _In_ HINSTANCE hInstance, _In_ int nCmdShow, _In_ json arg)
    : hInstance(hInstance), nCmdShow(nCmdShow), arg(arg) {
}

int WebViewWindow::Run() {
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
      KoiShellWebViewTitle,
      WS_OVERLAPPEDWINDOW,
      CW_USEDEFAULT,
      CW_USEDEFAULT,
      1366,
      768,
      nullptr,
      nullptr,
      hInstance,
      nullptr);

  if (!hWnd) {
    LogWithLastError(L"Failed to create window.");
    return 1;
  }

  ShowWindow(hWnd, nCmdShow);
  UpdateWindow(hWnd);

  MSG msg;
  while (GetMessageW(&msg, nullptr, 0, 0)) {
    TranslateMessage(&msg);
    DispatchMessageW(&msg);
  }

  return (int)msg.wParam;
}

LRESULT CALLBACK WebViewWindow::WndProc(
    _In_ HWND hWnd, _In_ UINT message, _In_ WPARAM wParam, _In_ LPARAM lParam) {
  switch (message) {
  case WM_DESTROY:
    PostQuitMessage(0);
    return 0;

  default:
    return DefWindowProcW(hWnd, message, wParam, lParam);
  }
}

int RunWebView(_In_ HINSTANCE hInstance, _In_ int nCmdShow, _In_ json arg) {
  WebViewWindow webviewWindow = WebViewWindow(hInstance, nCmdShow, arg);
  return webviewWindow.Run();
}

} // namespace KoiShell
