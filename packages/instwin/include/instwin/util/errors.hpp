#ifndef _INSTWIN_UTIL_ERRORS_
#define _INSTWIN_UTIL_ERRORS_

#include "instwin/stdafx.hpp"

namespace InstWin {

void LogAndFailWithLastError(HWND hWnd, const wchar_t *message);

} // namespace InstWin

#endif /* _INSTWIN_UTIL_ERRORS_ */
