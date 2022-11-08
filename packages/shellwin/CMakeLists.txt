cmake_minimum_required(VERSION 3.0)

############################################################
# Project
############################################################

project(koishell
    VERSION 0.1.0
    DESCRIPTION "Shell for Koishi Desktop."
    HOMEPAGE_URL "https://koishi.chat"
    LANGUAGES C CXX
)

set(CMAKE_C_STANDARD 23)
set(CMAKE_CXX_STANDARD 23)

set(CMAKE_MESSAGE_LOG_LEVEL VERBOSE)

############################################################
# Dependencies
############################################################

add_subdirectory(../shellcomm lib/shellcomm)
add_subdirectory(../wil lib/wil)
add_subdirectory(../webview2-cmake lib/webview2-cmake)

############################################################
# App
############################################################

add_executable(koishell
    WIN32
        src/main.cpp
        src/util/logger.cpp
        src/util/strings.cpp
        src/mode/webview.cpp
        src/mode/dialog.cpp
        src/koishell.rc
)

target_link_libraries(koishell
    PRIVATE
        shellcomm
        WIL
        webview2_static
        ShLwApi
        ComCtl32
)

target_include_directories(koishell
    PRIVATE
        ${PROJECT_SOURCE_DIR}/include
)

############################################################
# Compile Flags
############################################################

target_compile_definitions(koishell
    PRIVATE
        UNICODE
        _UNICODE
        WIN32
        _WIN32
        WINDOWS
        _WINDOWS
)

target_compile_options(koishell
    PRIVATE
        /utf-8
)

# Support Visual Studio "Edit and Continue"
target_compile_options(koishell PUBLIC /ZI)
target_link_options(koishell PUBLIC /INCREMENTAL)
