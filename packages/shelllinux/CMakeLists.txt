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

############################################################
# App
############################################################

add_executable(koishell
    src/main.cpp
)

target_link_libraries(koishell
    PRIVATE
        shellcomm
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
)

target_compile_options(koishell
    PUBLIC
        -fexec-charset=UTF-8
        -finput-charset=UTF-8
)
