package mockery

import (
	"reflect"
	"syscall"
	"unsafe"
)

var pageSize = syscall.Getpagesize()

//go:nosplit
func rawMemoryAccess(p uintptr, length int) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: p,
		Len:  length,
		Cap:  length,
	}))
}

func pageStart(ptr uintptr) uintptr {
	return ptr & ^(uintptr(pageSize - 1))
}

func replaceFunction(from, to uintptr) ([]byte, error) {
	jumpData, err := jmpToFunctionValue(to)
	if err != nil {
		return nil, err
	}
	f := rawMemoryAccess(from, len(jumpData))
	original := make([]byte, len(f))
	copy(original, f)

	copyToLocation(from, jumpData)
	return original, nil
}
