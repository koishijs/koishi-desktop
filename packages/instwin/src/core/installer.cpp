#include "instwin/core/installer.hpp"

namespace InstWin {

void InstallerWorker::run() {
  emit onProgress(0);
  emit onLog(std::format("Installation started at {}.", timeString()));

  emit onLog(std::format("Installation finished at {}.", timeString()));
  emit onResult(InstallResult::Success);
}

std::string InstallerWorker::timeString() {
  return std::format("{:%Y-%m-%d %H:%M:%OS}", std::chrono::system_clock::now());
}

Installer::Installer() {
  worker.moveToThread(&thread);
  connect(this, &Installer::start, &worker, &InstallerWorker::run);
  connect(&worker, &InstallerWorker::onLog, this, &Installer::onLog);
  connect(&worker, &InstallerWorker::onProgress, this, &Installer::onProgress);
  connect(&worker, &InstallerWorker::onResult, this, &Installer::onResult);
}

void Installer::start() {
  thread.start();
}

} // namespace InstWin
