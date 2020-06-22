// +build minimal

#define protected public
#define private public

#include "core-minimal.h"
#include "_cgo_export.h"

#ifndef QT_CORE_LIB
	#error ------------------------------------------------------------------
	#error please run: '$(go env GOPATH)/bin/qtsetup'
	#error more info here: https://github.com/therecipe/qt/wiki/Installation
	#error ------------------------------------------------------------------
#endif
#include <QTextDocument>
#include <QObject>
#include <QStringList>

