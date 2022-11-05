#ifndef _KOISHELL_STRINGS_
#define _KOISHELL_STRINGS_

#include "koishell/stdafx.hpp"

namespace KoiShell {

char *WideCharToUTF8(_In_ wchar_t *w);

wchar_t *UTF8ToWideChar(_In_ char *s);

} // namespace KoiShell

#endif /* _KOISHELL_STRINGS_ */
