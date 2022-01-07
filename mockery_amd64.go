//go:build amd64
// +build amd64

package mockery

// Assembles a jump to a function value
//go:nosplit
func jmpToFunctionValue(to uintptr) ([]byte, error) {
	return []byte{
		0x48, 0xBA,
		byte(to),
		byte(to >> 8),
		byte(to >> 16),
		byte(to >> 24),
		byte(to >> 32),
		byte(to >> 40),
		byte(to >> 48),
		byte(to >> 56), // movabs rdx,to
		0xFF, 0x22,     // jmp QWORD PTR [rdx]
	}, nil
}
