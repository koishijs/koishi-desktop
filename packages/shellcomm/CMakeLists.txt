cmake_minimum_required(VERSION 3.0)

############################################################
# Project
############################################################

project(shellcomm
    VERSION 0.1.0
    DESCRIPTION "Common utilities for Cordis Desktop Shell."
    HOMEPAGE_URL "https://koishi.chat"
    LANGUAGES C CXX
)

set(CMAKE_C_STANDARD 23)
set(CMAKE_CXX_STANDARD 23)

set(CMAKE_MESSAGE_LOG_LEVEL VERBOSE)

############################################################
# Dependencies
############################################################

set(JSON_Install OFF CACHE INTERNAL "")

add_subdirectory(../libb64 lib/libb64)
add_subdirectory(../nlohmann_json lib/json)

############################################################
# App
############################################################

add_library(shellcomm STATIC
    src/logger.cpp
    src/parse.cpp
)

target_link_libraries(shellcomm
    PUBLIC
        libb64
        nlohmann_json::nlohmann_json
)

target_include_directories(shellcomm
    PUBLIC
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
