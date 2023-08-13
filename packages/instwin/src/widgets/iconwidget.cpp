#include "instwin/widgets/iconwidget.hpp"
#include "ui_iconwidget.h"

IconWidget::IconWidget(QWidget *parent)
    : QWidget(parent), ui(new Ui::IconWidget) {
  ui->setupUi(this);
}

IconWidget::~IconWidget() {
  delete ui;
}
