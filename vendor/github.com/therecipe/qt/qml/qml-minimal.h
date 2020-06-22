// +build minimal

#pragma once

#ifndef GO_QTQML_H
#define GO_QTQML_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

struct QtQml_PackedString { char* data; long long len; void* ptr; };
struct QtQml_PackedList { void* data; long long len; };

#ifdef __cplusplus
}
#endif

#endif