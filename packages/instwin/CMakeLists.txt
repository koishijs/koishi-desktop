cmake_minimum_required(VERSION 3.5)

############################################################
# Project
############################################################

project(instwin
    VERSION 0.1.0
    DESCRIPTION "Cordis Desktop Installer."
    HOMEPAGE_URL "https://koishi.chat"
    LANGUAGES C CXX
)

set(CMAKE_AUTOUIC ON)
set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTORCC ON)

set(CMAKE_C_STANDARD 23)
set(CMAKE_CXX_STANDARD 23)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

set(CMAKE_MESSAGE_LOG_LEVEL VERBOSE)

# Use static link
set(CMAKE_MSVC_RUNTIME_LIBRARY "MultiThreaded$<$<CONFIG:Debug>:Debug>")

############################################################
# Dependencies
############################################################

# Qt Dependencies
find_package(QT NAMES Qt6 Qt5 REQUIRED COMPONENTS Widgets LinguistTools)
find_package(Qt${QT_VERSION_MAJOR} REQUIRED COMPONENTS Widgets LinguistTools)

add_subdirectory(../shellcomm lib/shellcomm)

############################################################
# App
############################################################

set(TS_FILES
    src/instwin_en_US.ts)

set(PROJECT_SOURCES
    src/main.cpp
    include/instwin/core/installer.hpp
    src/core/installer.cpp
    src/util/errors.cpp
    include/instwin/widgets/iconwidget.hpp
    src/widgets/iconwidget.cpp
    src/widgets/iconwidget.ui
    include/instwin/windows/mainwindow.hpp
    src/windows/mainwindow.cpp
    src/windows/mainwindow.ui
    ${TS_FILES}
    src/instwin.exe.manifest
    src/instwin.qrc
    src/instwin.rc
)

qt_add_executable(instwin
    MANUAL_FINALIZATION
        ${PROJECT_SOURCES}
)

qt_create_translation(
    QM_FILES
        ${CMAKE_SOURCE_DIR}
        ${TS_FILES}
)

target_link_libraries(instwin
    PRIVATE
        shellcomm
        Qt${QT_VERSION_MAJOR}::Widgets
        dwmapi
)

target_include_directories(instwin
    PRIVATE
        ${PROJECT_SOURCE_DIR}/include
)

set_target_properties(instwin
    PROPERTIES
        ${BUNDLE_ID_OPTION}
        WIN32_EXECUTABLE TRUE
)

############################################################
# Compile Flags
############################################################

target_compile_definitions(instwin
    PRIVATE
        UNICODE
        _UNICODE
        WIN32
        _WIN32
        WINDOWS
        _WINDOWS
)

if (MSVC)
    target_compile_options(instwin
        PRIVATE
            /utf-8
    )
else()
    target_compile_options(instwin
        PRIVATE
            -fexec-charset=UTF-8
            -finput-charset=UTF-8
    )
endif()

############################################################
# Support Visual Studio "Edit and Continue"
############################################################

if (CMAKE_GENERATOR MATCHES "Visual Studio")
    target_compile_options(instwin
        PRIVATE
            /ZI
    )

    target_link_options(instwin
        PRIVATE
            /INCREMENTAL
    )
endif()

############################################################
# Finalize
############################################################

qt_finalize_executable(instwin)
