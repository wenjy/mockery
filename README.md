# mockery
Go monkery patching

**只用于Go测试**

- 代码请参考：

[monkey](https://github.com/bouk/monkey) 

[go-mpatch](https://github.com/undefinedlabs/go-mpatch)

- 相关技术博客 

[go patch 译](https://blog.wenjy.top/2022/01/05/golang/go-patch.html) 

[go patch](https://bou.ke/blog/monkey-patching-in-go/)

## 兼容性

- **Go版本** 已完成 `go1.7` 到 `go1.17` 的测试

- **系统架构** `x86`、`amd64`

- **操作系统** `macos`、`linux`、`windows`

## 限制

- 目表函数如果使用内联（inlined），将不能替换，可以使用指令`//go:noinline`或者`gcflags=-l`来构建，告诉go编译器禁用内联

- 需要对包含可执行代码的内存页有写权限，一些操作系统可能会限制这种访问

- 不是线程安全的

## 使用示例

### 替换一个函数

```go
//go:noinline
func methodA() int { return 1 }

//go:noinline
func methodB() int { return 2 }

func TestPatcher(t *testing.T) {
	patch, err := mpatch.PatchMethod(methodA, methodB)
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
```

### 使用`reflect.ValueOf`替换函数

```go
//go:noinline
func methodA() int { return 1 }

//go:noinline
func methodB() int { return 2 }

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
```

### 使用自定义函数来替换

```go
//go:noinline
func methodA() int { return 1 }

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
```

### 替换结构体指针的方法

```go
type myStruct struct {
}

//go:noinline
func (s *myStruct) Method() int {
	return 1
}

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
```


### 替换结构体的方法

```go
type myStruct struct {
}

//go:noinline
func (s myStruct) ValueMethod() int {
	return 1
}

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
```

### 使用反射来替换结构体指针方法

```go

type myStruct struct {
}

//go:noinline
func (s *myStruct) Method() int {
	return 1
}
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
```

### 使用反射方法+自定义函数来替换结构体指针方法

```go
type myStruct struct {
}

//go:noinline
func (s *myStruct) Method() int {
	return 1
}

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
```

### 使用反射值+自定义函数来替换结构体指针方法

```go
type myStruct struct {
}

//go:noinline
func (s *myStruct) Method() int {
	return 1
}
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
```

### 使用反射值来替换结构体指针方法

```go
type myStruct struct {
}

//go:noinline
func (s *myStruct) Method() int {
	return 1
}

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
```