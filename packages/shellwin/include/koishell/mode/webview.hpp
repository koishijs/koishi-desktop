#ifndef _KOISHELL_MODE_WEBVIEW_
#define _KOISHELL_MODE_WEBVIEW_

#include <windows.h>

#include "nlohmann/json.hpp"

#include "koishell/util/logger.hpp"

using json = nlohmann::json;

namespace KoiShell {

int RunWebView(_In_ HINSTANCE hInstance, _In_ int nCmdShow, _In_ json arg);

LRESULT CALLBACK WndProcWebView(
    _In_ HWND hWnd, _In_ UINT message, _In_ WPARAM wParam, _In_ LPARAM lParam);

} // namespace KoiShell

#endif /* _KOISHELL_MODE_WEBVIEW_ */
