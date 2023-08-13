#include "instwin/main.hpp"

#include <QApplication>
#include <QFile>
#include <QLocale>
#include <QTranslator>

int main(int argc, char *argv[]) {
  QApplication app(argc, argv);

  // Load global stylesheet
  QFile qssFile(":/instwin.qss");
  qssFile.open(QFile::ReadOnly | QFile::Text);
  QTextStream qssTextStream(&qssFile);
  app.setStyleSheet(qssTextStream.readAll());

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
