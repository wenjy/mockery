package mockery

import (
	"reflect"
	"testing"
)

//go:noinline
func methodA() int { return 1 }

//go:noinline
func methodB() int { return 2 }

type myStruct struct {
}

//go:noinline
func (s *myStruct) Method() int {
	return 1
}

//go:noinline
func (s myStruct) ValueMethod() int {
	return 1
}

// 测试替换方法
func TestPatcher(t *testing.T) {
	patch, err := PatchMethod(methodA, methodB)
	if err != nil {
		t.Fatal(err)
	}
	if methodA() != 2 {
		t.Fatal("The patch did not work")
	}

	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
	if methodA() != 1 {
		t.Fatal("The unpatch did not work")
	}
}

// 测试用反射替换方法
func TestPatcherUsingReflect(t *testing.T) {
	reflectA := reflect.ValueOf(methodA)
	patch, err := PatchMethodByReflectValue(reflectA, methodB)
	if err != nil {
		t.Fatal(err)
	}
	if methodA() != 2 {
		t.Fatal("The patch did not work")
	}

	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
	if methodA() != 1 {
		t.Fatal("The unpatch did not work")
	}
}

// 测试使用自定义函数替换方法
func TestPatcherUsingMakeFunc(t *testing.T) {
	reflectA := reflect.ValueOf(methodA)
	patch, err := PatchMethodWithMakeFuncValue(reflectA,
		func(args []reflect.Value) (results []reflect.Value) {
			return []reflect.Value{reflect.ValueOf(42)}
		})
	if err != nil {
		t.Fatal(err)
	}
	if methodA() != 42 {
		t.Fatal("The patch did not work")
	}

	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
	if methodA() != 1 {
		t.Fatal("The unpatch did not work")
	}
}

// 测试替换结构体指针的方法
func TestInstancePatcher(t *testing.T) {
	mStruct := myStruct{}

	var patch *Patch
	var err error
	patch, err = PatchInstanceMethod(reflect.TypeOf(mStruct), "Method", func(m *myStruct) int {
		patch.Unpatch()
		defer patch.Patch()
		return 41 + m.Method()
	})
	if err != nil {
		t.Fatal(err)
	}

	if mStruct.Method() != 42 {
		t.Fatal("The patch did not work")
	}
	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
	if mStruct.Method() != 1 {
		t.Fatal("The unpatch did not work")
	}
}

// 测试替换结构体的方法
func TestInstanceValuePatcher(t *testing.T) {
	mStruct := myStruct{}

	var patch *Patch
	var err error
	patch, err = PatchInstanceMethod(reflect.TypeOf(mStruct), "ValueMethod", func(m myStruct) int {
		patch.Unpatch()
		defer patch.Patch()
		return 41 + m.Method()
	})
	if err != nil {
		t.Fatal(err)
	}

	if mStruct.ValueMethod() != 42 {
		t.Fatal("The patch did not work")
	}
	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
	if mStruct.ValueMethod() != 1 {
		t.Fatal("The unpatch did not work")
	}
}

// 测试使用反射来替换结构体指针方法
func TestPatchMethodByReflect(t *testing.T) {
	mStruct := myStruct{}

	target := reflect.TypeOf(mStruct)
	target = reflect.PtrTo(target)
	m, _ := target.MethodByName("Method")

	var patch *Patch
	var err error
	patch, err = PatchMethodByReflect(m, func(m *myStruct) int {
		patch.Unpatch()
		defer patch.Patch()
		return 41 + m.Method()
	})

	if err != nil {
		t.Fatal(err)
	}

	if mStruct.Method() != 42 {
		t.Fatal("The patch did not work")
	}
	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
	if mStruct.Method() != 1 {
		t.Fatal("The unpatch did not work")
	}
}

// 测试使用反射方法+自定义函数来替换结构体指针方法
func TestPatchMethodWithMakeFunc(t *testing.T) {
	mStruct := myStruct{}

	target := reflect.TypeOf(mStruct)
	target = reflect.PtrTo(target)
	m, _ := target.MethodByName("Method")

	var patch *Patch
	var err error
	patch, err = PatchMethodWithMakeFunc(m, func(args []reflect.Value) (results []reflect.Value) {
		return []reflect.Value{reflect.ValueOf(42)}
	})

	if err != nil {
		t.Fatal(err)
	}

	if mStruct.Method() != 42 {
		t.Fatal("The patch did not work")
	}
	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
	if mStruct.Method() != 1 {
		t.Fatal("The unpatch did not work")
	}
}

// 测试使用反射值+自定义函数来替换结构体指针方法
func TestPatchMethodWithMakeFuncValue(t *testing.T) {
	mStruct := myStruct{}

	target := reflect.TypeOf(mStruct)
	target = reflect.PtrTo(target)
	m, _ := target.MethodByName("Method")

	var patch *Patch
	var err error
	patch, err = PatchMethodWithMakeFuncValue(m.Func, func(args []reflect.Value) (results []reflect.Value) {
		return []reflect.Value{reflect.ValueOf(42)}
	})

	if err != nil {
		t.Fatal(err)
	}

	if mStruct.Method() != 42 {
		t.Fatal("The patch did not work")
	}
	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
	if mStruct.Method() != 1 {
		t.Fatal("The unpatch did not work")
	}
}

// 测试使用反射值来替换结构体指针方法
func TestPatchMethodByReflectValue(t *testing.T) {
	mStruct := myStruct{}

	target := reflect.TypeOf(mStruct)
	target = reflect.PtrTo(target)
	m, _ := target.MethodByName("Method")

	var patch *Patch
	var err error
	patch, err = PatchMethodByReflectValue(m.Func, func(m *myStruct) int {
		patch.Unpatch()
		defer patch.Patch()
		return 41 + m.Method()
	})

	if err != nil {
		t.Fatal(err)
	}

	if mStruct.Method() != 42 {
		t.Fatal("The patch did not work")
	}
	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
	if mStruct.Method() != 1 {
		t.Fatal("The unpatch did not work")
	}
}
