#ifndef _KOISHELL_STRINGS_
#define _KOISHELL_STRINGS_

#include <windows.h>

namespace KoiShell {

char *WideCharToUTF8(wchar_t *w);

wchar_t *UTF8ToWideChar(_In_ char *s);

} // namespace KoiShell

#endif /* _KOISHELL_STRINGS_ */
