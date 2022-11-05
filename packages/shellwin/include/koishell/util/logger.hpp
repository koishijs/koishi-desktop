#ifndef _KOISHELL_LOGGER_
#define _KOISHELL_LOGGER_

#include "koishell/stdafx.hpp"

namespace KoiShell {

void FailWithLastError();

void LogAndFailWithLastError(const wchar_t *messages);

void LogAndFail(const wchar_t *messages);

void LogAndFail(const std::wstring &messages);

void CheckFailure(HRESULT hr, const std::wstring &message = L"Error");

} // namespace KoiShell

#endif /* _KOISHELL_LOGGER_ */
