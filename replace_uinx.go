//go:build !windows
// +build !windows

package mockery

import "syscall"

var writeAccess = syscall.PROT_READ | syscall.PROT_WRITE | syscall.PROT_EXEC
var readAccess = syscall.PROT_READ | syscall.PROT_EXEC

//go:nosplit
func mprotectCrossPage(addr uintptr, length int, prot int) error {
	for p := pageStart(addr); p < addr+uintptr(length); p += uintptr(pageSize) {
		page := rawMemoryAccess(p, pageSize)
		if err := syscall.Mprotect(page, prot); err != nil {
			return err
		}
	}
	return nil
}

// this function is super unsafe
// It copies a slice to a raw memory location, disabling all memory protection before doing so.
func copyToLocation(location uintptr, data []byte) error {
	f := rawMemoryAccess(location, len(data))

	if err := mprotectCrossPage(location, len(data), writeAccess); err != nil {
		return err
	}
	copy(f, data[:])
	if err := mprotectCrossPage(location, len(data), readAccess); err != nil {
		return err
	}
	return nil
}
