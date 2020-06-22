// +build minimal

package qml

//#include <stdint.h>
//#include <stdlib.h>
//#include <string.h>
//#include "qml-minimal.h"
import "C"
import (
	"github.com/therecipe/qt/core"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

func cGoFreePacked(ptr unsafe.Pointer) { core.NewQByteArrayFromPointer(ptr).DestroyQByteArray() }
func cGoUnpackString(s C.struct_QtQml_PackedString) string {
	defer cGoFreePacked(s.ptr)
	if int(s.len) == -1 {
		return C.GoString(s.data)
	}
	return C.GoStringN(s.data, C.int(s.len))
}
func cGoUnpackBytes(s C.struct_QtQml_PackedString) []byte {
	defer cGoFreePacked(s.ptr)
	if int(s.len) == -1 {
		gs := C.GoString(s.data)
		return *(*[]byte)(unsafe.Pointer(&gs))
	}
	return C.GoBytes(unsafe.Pointer(s.data), C.int(s.len))
}
func unpackStringList(s string) []string {
	if len(s) == 0 {
		return make([]string, 0)
	}
	return strings.Split(s, "¡¦!")
}

var (
	helper      *core.QObject
	helperMutex sync.Mutex
	helperMap   []string
)

func init() {

	helper = core.NewQObject(nil)
	helper.ConnectObjectNameChanged(func(pl string) {
		for _, p := range strings.Split(pl, "|") {
			ptr, err := strconv.ParseUint(p, 10, 64)
			if err != nil || ptr == 0 {
				return
			}
			C.QJSValue_DestroyQJSValue(unsafe.Pointer(uintptr(ptr)))
			C.free(unsafe.Pointer(uintptr(ptr)))
		}
	})
}
