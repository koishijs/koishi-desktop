#ifndef _KOISHELL_LOGGER_
#define _KOISHELL_LOGGER_

#include <iostream>
#include <windows.h>

namespace KoiShell {

void LogW(const wchar_t *messages);

void LogLastError();

void LogWithLastError(const wchar_t *messages);

} // namespace KoiShell

#endif /* _KOISHELL_LOGGER_ */
