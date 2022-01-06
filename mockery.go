package mockery

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

type Patch struct {
	targetBytes []byte
	target      *reflect.Value
	replacement *reflect.Value
}

var (
	lock    = sync.Mutex{}
	patches = make(map[uintptr]*Patch)
)

//go:linkname getInternalPtrFromValue reflect.(*Value).pointer
func getInternalPtrFromValue(v *reflect.Value) unsafe.Pointer

// 把target方法替换为replacement方法
func PatchMethod(target, replacement interface{}) (*Patch, error) {
	t := getValueFrom(target)
	r := getValueFrom(replacement)

	if err := isPatchable(&t, &r); err != nil {
		return nil, err
	}

	patch := &Patch{target: &t, replacement: &r}

	if err := applyPatch(patch); err != nil {
		return nil, err
	}
	return patch, nil
}

// 把target结构体的methodName方法替换为replacement方法
func PatchInstanceMethod(target reflect.Type, methodName string, replacement interface{}) (*Patch, error) {
	m, ok := target.MethodByName(methodName)

	if !ok && target.Kind() == reflect.Struct {
		target = reflect.PtrTo(target)
		m, ok = target.MethodByName(methodName)
	}
	if !ok {
		return nil, errors.New(fmt.Sprintf("Method '%v' not found", methodName))
	}

	return PatchMethodByReflect(m, replacement)
}

// 把target反射方法替换为replacement方法
func PatchMethodByReflect(target reflect.Method, replacement interface{}) (*Patch, error) {
	return PatchMethodByReflectValue(target.Func, replacement)
}

// 把target反射方法使用自定义函数来替换
func PatchMethodWithMakeFunc(target reflect.Method, fn func(args []reflect.Value) (results []reflect.Value)) (*Patch, error) {
	return PatchMethodByReflect(target, reflect.MakeFunc(target.Type, fn))
}

// 把target反射值使用自定义函数来替换
func PatchMethodWithMakeFuncValue(target reflect.Value, fn func(args []reflect.Value) (results []reflect.Value)) (*Patch, error) {
	return PatchMethodByReflectValue(target, reflect.MakeFunc(target.Type(), fn))
}

// 把target反射值替换为replacement方法
func PatchMethodByReflectValue(target reflect.Value, replacement interface{}) (*Patch, error) {
	t := &target
	r := getValueFrom(replacement)
	if err := isPatchable(t, &r); err != nil {
		return nil, err
	}
	patch := &Patch{target: t, replacement: &r}
	if err := applyPatch(patch); err != nil {
		return nil, err
	}
	return patch, nil
}

func (p *Patch) Patch() error {
	if p == nil {
		return errors.New("patch is nil")
	}
	if err := isPatchable(p.target, p.replacement); err != nil {
		return err
	}
	if err := applyPatch(p); err != nil {
		return err
	}
	return nil
}

func (p *Patch) Unpatch(target interface{}) error {
	if p == nil {
		return errors.New("patch is nil")
	}
	return applyUnpatch(p)
}

func getValueFrom(data interface{}) reflect.Value {
	if cValue, ok := data.(reflect.Value); ok {
		return cValue
	} else {
		return reflect.ValueOf(data)
	}
}

func isPatchable(target, replacement *reflect.Value) error {
	lock.Lock()
	defer lock.Unlock()

	if target.Kind() != reflect.Func {
		return errors.New("the target is not a Func")
	}

	if replacement.Kind() != reflect.Func {
		return errors.New("the replacement is not a Func")
	}

	if target.Type() != replacement.Type() {
		return errors.New(fmt.Sprintf("the target and redirection doesn't have the same type: %s != %s", target.Type(), replacement.Type()))
	}
	if _, ok := patches[target.Pointer()]; ok {
		return errors.New("the target is already patched")
	}
	return nil
}

func applyPatch(patch *Patch) error {
	lock.Lock()
	defer lock.Unlock()

	patch.targetBytes = replaceFunction(patch.target.Pointer(), (uintptr)(getInternalPtrFromValue(patch.replacement)))
	patches[patch.target.Pointer()] = patch

	return nil
}

func applyUnpatch(patch *Patch) error {
	lock.Lock()
	defer lock.Unlock()

	if patch.targetBytes == nil || len(patch.targetBytes) == 0 {
		return errors.New("the target is not patched")
	}
	tptr := patch.target.Pointer()
	if _, ok := patches[tptr]; !ok {
		return errors.New("the target is not patched")
	}
	unpatch(tptr, patch)
	delete(patches, tptr)
	return nil
}

func unpatch(target uintptr, p *Patch) {
	copyToLocation(target, p.targetBytes)
}
