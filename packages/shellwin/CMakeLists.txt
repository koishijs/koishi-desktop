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
# App
############################################################

add_executable(koishell
    src/main.cpp
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
