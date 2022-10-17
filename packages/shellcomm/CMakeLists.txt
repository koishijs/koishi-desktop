cmake_minimum_required(VERSION 3.0)

############################################################
# Project
############################################################

project(shellcomm
    VERSION 0.1.0
    DESCRIPTION "Common utilities for Koishi Desktop Shell."
    HOMEPAGE_URL "https://koishi.chat"
    LANGUAGES C CXX
)

set(CMAKE_C_STANDARD 23)
set(CMAKE_CXX_STANDARD 23)

set(CMAKE_MESSAGE_LOG_LEVEL VERBOSE)

############################################################
# Dependencies
############################################################

add_subdirectory(../libb64 lib/libb64)

############################################################
# App
############################################################

add_library(shellcomm STATIC
)

target_link_libraries(shellcomm
    PRIVATE
        libb64
)

target_include_directories(shellcomm
    PRIVATE
        ${PROJECT_SOURCE_DIR}/include
)

############################################################
# Compile Flags
############################################################

target_compile_definitions(shellcomm
    PRIVATE
        UNICODE
        _UNICODE
)

if(MSVC)
    target_compile_definitions(shellcomm
        PRIVATE
            WIN32
            _WIN32
            WINDOWS
            _WINDOWS
    )
    target_compile_options(shellcomm
        PRIVATE
            /utf-8
    )
else()
    target_compile_options(shellcomm
        PUBLIC
            -fexec-charset=UTF-8
            -finput-charset=UTF-8
    )
endif()
