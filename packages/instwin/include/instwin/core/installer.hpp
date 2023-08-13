#ifndef _INSTWIN_CORE_INSTALLER_
#define _INSTWIN_CORE_INSTALLER_

#include "instwin/stdafx.hpp"

#include <chrono>

#include <QThread>

namespace InstWin {

enum InstallResult {
  Undefined = 0,
  Success = 1,
  Fail = 2,
  Warning = 3,
};

class InstallerWorker : public QObject {
  Q_OBJECT

public slots:
  void run();

signals:
  void onLog(std::string s);
  void onProgress(int progress);
  void onResult(InstallResult result);

private:
  std::string timeString();
};

class Installer : public QObject {
  Q_OBJECT

public:
  Installer();

public slots:
  void start();

signals:
  void onLog(std::string s);
  void onProgress(int progress);
  void onResult(InstallResult result);

private:
  QThread thread;
  InstallerWorker worker;
};

} // namespace InstWin

#endif /* _INSTWIN_CORE_INSTALLER_ */
