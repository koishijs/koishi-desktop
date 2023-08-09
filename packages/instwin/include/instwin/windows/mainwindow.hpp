#ifndef _INSTWIN_WINDOWS_MAINWINDOW_
#define _INSTWIN_WINDOWS_MAINWINDOW_

#include "instwin/stdafx.hpp"

#include <QMainWindow>

#include "instwin/util/errors.hpp"

QT_BEGIN_NAMESPACE
namespace Ui {
class MainWindow;
}
QT_END_NAMESPACE

class MainWindow : public QMainWindow {
  Q_OBJECT

public:
  MainWindow(QWidget *parent = nullptr);
  ~MainWindow();

private:
  Ui::MainWindow *ui;

  void InitializeWindowStyle();
};

#endif /* _INSTWIN_WINDOWS_MAINWINDOW_ */
