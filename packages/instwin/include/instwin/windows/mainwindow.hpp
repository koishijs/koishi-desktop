#ifndef _INSTWIN_WINDOWS_MAINWINDOW_
#define _INSTWIN_WINDOWS_MAINWINDOW_

#include "instwin/stdafx.hpp"

#include <QMainWindow>

#include "instwin/core/installer.hpp"
#include "instwin/util/errors.hpp"

QT_BEGIN_NAMESPACE
namespace Ui {
class MainWindow;
}
QT_END_NAMESPACE

class MainWindow : public QMainWindow {
  Q_OBJECT

public:
  MainWindow();
  ~MainWindow();

public slots:
  void navigateToProgressPage();
  void navigateToFinishPage();
  void appendProgressLog(std::string s);

private:
  Ui::MainWindow *ui;

  InstWin::Installer installer;

  void initializeWindowStyle();
};

#endif /* _INSTWIN_WINDOWS_MAINWINDOW_ */
