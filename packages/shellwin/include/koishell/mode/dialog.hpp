#ifndef _KOISHELL_MODE_DIALOG_
#define _KOISHELL_MODE_DIALOG_

#include "koishell/stdafx.hpp"

#include "koishell/util/logger.hpp"
#include "koishell/util/strings.hpp"
#include "shellcomm/parse.hpp"

using njson = nlohmann::json;

namespace KoiShell {

int RunDialog(_In_ HINSTANCE hInstance, _In_ njson arg);

} // namespace KoiShell

#endif /* _KOISHELL_MODE_DIALOG_ */
