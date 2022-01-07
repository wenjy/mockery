//go:build !386 && !amd64
// +build !386,!amd64

package mockery

import (
	"fmt"
	"runtime"
)

// Assembles a jump to a function value
//go:nosplit
func jmpToFunctionValue(to uintptr) ([]byte, error) {
	return nil, fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
}
