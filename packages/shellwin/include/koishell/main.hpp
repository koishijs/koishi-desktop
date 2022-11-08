#ifndef _KOISHELL_MAIN_
#define _KOISHELL_MAIN_

#include "koishell/stdafx.hpp"

#include "koishell/mode/dialog.hpp"
#include "koishell/mode/webview.hpp"
#include "koishell/util/logger.hpp"
#include "koishell/util/strings.hpp"
#include "shellcomm/logger.hpp"
#include "shellcomm/parse.hpp"

int WINAPI wWinMain(
    HINSTANCE hInstance, HINSTANCE hPrevInstance, PWSTR pCmdLine, int nCmdShow);

#endif /* _KOISHELL_MAIN_ */
