#include "instwin/main.hpp"

#include <QApplication>
#include <QLocale>
#include <QTranslator>

int main(int argc, char *argv[]) {
  const QApplication app(argc, argv);

  QTranslator translator;
  const QStringList uiLanguages = QLocale::system().uiLanguages();
  for (const QString &locale : uiLanguages) {
    const QString baseName = "QtPlaygroundA_" + QLocale(locale).name();
    if (translator.load(":/i18n/" + baseName)) {
      app.installTranslator(&translator);
      break;
    }
  }
  MainWindow mainWindow;
  mainWindow.show();
  return app.exec();
}
