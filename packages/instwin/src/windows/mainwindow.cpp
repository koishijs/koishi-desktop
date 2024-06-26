#include "instwin/windows/mainwindow.hpp"
#include "ui_mainwindow.h"

MainWindow::MainWindow()
    : QMainWindow(nullptr), ui(new Ui::MainWindow), installer(this) {
  // Setup UI
  ui->setupUi(this);

  connect(
      &installer,
      &InstWin::Installer::onLog,
      this,
      &MainWindow::appendProgressLog);
  connect(
      ui->installButton,
      &QPushButton::clicked,
      &installer,
      &InstWin::Installer::start);

  initializeWindowStyle();
}

MainWindow::~MainWindow() {
  delete ui;
}

void MainWindow::initializeWindowStyle() {
  HWND hWnd = reinterpret_cast<HWND>(winId());

  // Fix window
  // setWindowFlag(Qt::MSWindowsFixedSizeDialogHint);

  // Enable DWM transparent
  OSVERSIONINFOEXW osvi{0};
  osvi.dwOSVersionInfoSize = sizeof(OSVERSIONINFOEXW);
  if (!GetVersionExW(reinterpret_cast<LPOSVERSIONINFOW>(&osvi)))
    InstWin::LogAndFailWithLastError(hWnd, L"Failed to get OS version info.");

  unsigned char supports = 0;

  if (osvi.dwBuildNumber >= 17763)
    supports++; // Supports Windows 10 immersive dark mode (19), supports = 1
  if (osvi.dwBuildNumber >= 18985)
    supports++; // Supports Windows 10 immersive dark mode (20), supports = 2
  if (osvi.dwBuildNumber >= 22000)
    supports++; // Supports Windows 11 Mica, supports = 3
  if (osvi.dwBuildNumber >= 22523)
    supports++; // Supports Windows 11 Mica Tabbed, supports = 4

  if (supports) {
    MARGINS dwmExtendFrameIntoClientAreaMargins = {-1};
    DwmExtendFrameIntoClientArea(hWnd, &dwmExtendFrameIntoClientAreaMargins);
    int dwmUseDarkMode = 1;
    DwmSetWindowAttribute(
        hWnd,
        supports >= 2 ? 20 // DWMWINDOWATTRIBUTE::DWMWA_USE_IMMERSIVE_DARK_MODE
                      :    // = 20 (starting from 18985)
                           // DWMWINDOWATTRIBUTE::DWMWA_USE_IMMERSIVE_DARK_MODE
            19,            // = 19 (before 18985)

        &dwmUseDarkMode,
        sizeof(dwmUseDarkMode));
    unsigned int dwmCornerPreference =
        2; // DWM_WINDOW_CORNER_PREFERENCE::DWMWCP_ROUND
    DwmSetWindowAttribute(
        hWnd,
        33, // DWMWINDOWATTRIBUTE::DWMWA_WINDOW_CORNER_PREFERENCE
        &dwmCornerPreference,
        sizeof(dwmCornerPreference));
    int dwmMica = 1;
    DwmSetWindowAttribute(
        hWnd,
        1029, // DWMWINDOWATTRIBUTE::DWMWA_MICA
        &dwmMica,
        sizeof(dwmMica));
    unsigned int dwmSystemBackdropType =
        supports >= 4 ? 4 : 2; // DWM_SYSTEMBACKDROP_TYPE::DWMSBT_MAINWINDOW
    DwmSetWindowAttribute(
        hWnd,
        38, // DWMWINDOWATTRIBUTE::DWMWA_SYSTEMBACKDROP_TYPE
        &dwmSystemBackdropType,
        sizeof(dwmSystemBackdropType));
  }
}

void MainWindow::navigateToProgressPage() {
  ui->centralWidget->setCurrentIndex(1);
}

void MainWindow::navigateToFinishPage() {
  ui->centralWidget->setCurrentIndex(2);
}

void MainWindow::appendProgressLog(std::string s) {
  ui->progressPageLog->moveCursor(QTextCursor::End);
  ui->progressPageLog->insertPlainText(s.c_str());
  ui->progressPageLog->moveCursor(QTextCursor::End);
  ui->progressPageLog->insertPlainText("\n");
  QScrollBar *verticalScrollBar = ui->progressPageLog->verticalScrollBar();
  verticalScrollBar->setValue(verticalScrollBar->maximum());
}
