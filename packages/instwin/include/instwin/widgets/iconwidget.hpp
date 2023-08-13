#ifndef _INSTWIN_WIDGETS_ICONWIDGET_
#define _INSTWIN_WIDGETS_ICONWIDGET_

#include "instwin/stdafx.hpp"

#include <QWidget>

QT_BEGIN_NAMESPACE
namespace Ui {
class IconWidget;
}
QT_END_NAMESPACE

class IconWidget : public QWidget {
  Q_OBJECT

public:
  explicit IconWidget(QWidget *parent = nullptr);
  ~IconWidget();

private:
  Ui::IconWidget *ui;
};

#endif /* _INSTWIN_WIDGETS_ICONWIDGET_ */
