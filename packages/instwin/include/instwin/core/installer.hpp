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

class InstallerWorker : public QThread {
  Q_OBJECT

public:
  InstallerWorker(QObject *parent);

signals:
  void onLog(std::string s);
  void onProgress(int progress);
  void onResult(InstallResult result);

private:
  void run() override;

  std::string timeString();
};

class Installer : public QObject {
  Q_OBJECT

public:
  Installer(QObject *parent);

public slots:
  void start();

signals:
  void onLog(std::string s);
  void onProgress(int progress);
  void onResult(InstallResult result);

private:
  InstallerWorker worker;
};

} // namespace InstWin

#endif /* _INSTWIN_CORE_INSTALLER_ */
