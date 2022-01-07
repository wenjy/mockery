//go:build 386
// +build 386

package mockery

// Assembles a jump to a function value
//go:nosplit
func jmpToFunctionValue(to uintptr) ([]byte, error) {
	return []byte{
		0xBA,
		byte(to),
		byte(to >> 8),
		byte(to >> 16),
		byte(to >> 24), // mov edx,to
		0xFF, 0x22,     // jmp DWORD PTR [edx]
	}
}, nil
