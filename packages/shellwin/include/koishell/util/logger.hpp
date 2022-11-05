#ifndef _KOISHELL_LOGGER_
#define _KOISHELL_LOGGER_

#include "koishell/stdafx.hpp"

namespace KoiShell {

void FailWithLastError();

void LogAndFailWithLastError(const wchar_t *messages);

void LogAndFail(const wchar_t *messages);

void CheckFailure(HRESULT hr, const std::wstring &message = L"Error");

} // namespace KoiShell

#define CHECK_FAILURE_STRINGIFY(arg) #arg
#define CHECK_FAILURE_FILE_LINE(file, line)                                    \
  ([](HRESULT hr) {                                                            \
    CheckFailure(                                                              \
        hr,                                                                    \
        L"Failure at " CHECK_FAILURE_STRINGIFY(                                \
            file) L"(" CHECK_FAILURE_STRINGIFY(line) L")");                    \
  })
#define CHECK_FAILURE CHECK_FAILURE_FILE_LINE(__FILE__, __LINE__)

#endif /* _KOISHELL_LOGGER_ */
