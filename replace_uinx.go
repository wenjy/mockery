//go:build !windows
// +build !windows

package mockery

import "syscall"

var writeAccess = syscall.PROT_READ | syscall.PROT_WRITE | syscall.PROT_EXEC
var readAccess = syscall.PROT_READ | syscall.PROT_EXEC

//go:nosplit
func mprotectCrossPage(addr uintptr, length int, prot int) {
	for p := pageStart(addr); p < addr+uintptr(length); p += uintptr(pageSize) {
		page := rawMemoryAccess(p, pageSize)
		err := syscall.Mprotect(page, prot)
		if err != nil {
			panic(err)
		}
	}
}

// this function is super unsafe
// aww yeah
// It copies a slice to a raw memory location, disabling all memory protection before doing so.
func copyToLocation(location uintptr, data []byte) {
	f := rawMemoryAccess(location, len(data))

	mprotectCrossPage(location, len(data), writeAccess)
	copy(f, data[:])
	mprotectCrossPage(location, len(data), readAccess)
}
