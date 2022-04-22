# Benchmark in Golang
## 1、 Fib 测试用例
```go
//fib.go
package main

func fib(n int) int {
	if n == 0 || n == 1 {
		return n
	}
	return fib(n-2) + fib(n-1)
}
```
```go
//fib_test.go
package main

import "testing"

func BenchmarkFib(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fib(30)
	}
}
```
1. benchmark 和普通的单元测试用例一样，位于 ```_test.go``` 文件中；  
2. 函数名以```Benchmark```开头，参数是```（b *testing.B）```
## 2、 运行测试用例
1. 运行当前```package```内的用例
   ```shell
   go test packageName 或
   go test.
   ```
2. 运行子```package```内的用例
   ```shell
   go test packageName/sunPackageName 或
   go test ./subPackageName
   ```
3. 递归测试当前目录下所有```package```
   ```shell
   go test packageName/... 或
   go test ./...
> ```go test``` 命令默认不运行 benchmark 测试用例，如果想运行，需要加上```-bench```参数    
> ```-bench```参数支持传入一个正则表达式，从而运行指定的Benchmark，例如只运行以Fib结尾的benchmark：
```shell
go test -bench='Fib$' .
```
benchmark 用例的参数 b *testing.B，有个属性 b.N 表示这个用例需要运行的次数。b.N 对于每个用例都是不一样的。

那这个值是如何决定的呢？b.N 从 1 开始，如果该用例能够在 1s 内完成，b.N 的值便会增加，再次执行。b.N 的值大概以 1, 2, 3, 5, 10, 20, 30, 50, 100 这样的序列递增，越到后面，增加得越快。我们仔细观察上述例子的输出：
```
BenchmarkFib-8               202           5980669 ns/op
```
BenchmarkFib-8 中的 -8 即 GOMAXPROCS，默认等于 CPU 核数。可以通过 -cpu 参数改变 GOMAXPROCS，-cpu 支持传入一个列表作为参数，例如：
```
$ go test -bench='Fib$' -cpu=2,4 .
goos: darwin
goarch: amd64
pkg: example
BenchmarkFib-2               206           5774888 ns/op
BenchmarkFib-4               205           5799426 ns/op
PASS
ok      example 3.563s
```
在这个例子中，改变 CPU 的核数对结果几乎没有影响，因为这个 Fib 的调用是串行的。

202 和 5980669 ns/op 表示用例执行了 202 次，每次花费约 0.006s。总耗时比 1s 略多。
### 2.1、 指定运行时间
默认运行时间为1s
```shell
go test -bench='Fib$' -benchtime=5s .

goos: darwin
goarch: arm64
pkg: benchmarktest
BenchmarkFib-8   	     859	   6922489 ns/op
PASS
ok  	benchmarktest	7.098s
```
> 实际执行的时间是大于5s，测试用例编译、执行、销毁等是需要时间的。
### 2.2、 指定运行次数
```
$ go test -bench='Fib$' -benchtime=50x .
goos: darwin
goarch: amd64
pkg: example
BenchmarkFib-8                50           6121066 ns/op
PASS
ok      example 0.319s
```
### 2.3、 指定运行轮数
> -count 参数可以用来设置 benchmark 的轮数。例如，进行 3 轮 benchmark。
```
$ go test -bench='Fib$' -benchtime=5s -count=3 .
goos: darwin
goarch: amd64
pkg: example
BenchmarkFib-8               975           5946624 ns/op
BenchmarkFib-8              1023           5820582 ns/op
BenchmarkFib-8               961           6096816 ns/op
PASS
ok      example 19.463s
```
## 3、 查看内存分配情况
-benchmem 参数可以度量内存分配的次数。内存分配次数也性能也是息息相关的，例如不合理的切片容量，将导致内存重新分配，带来不必要的开销。  
在下面的例子中，generateWithCap 和 generate 的作用是一致的，生成一组长度为 n 的随机序列。唯一的不同在于，generateWithCap 创建切片时，将切片的容量(capacity)设置为 n，这样切片就会一次性申请 n 个整数所需的内存。
```go
// generate_test.go
package main

import (
	"math/rand"
	"testing"
	"time"
)

func generateWithCap(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, 0, n)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}

func generate(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, 0)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}

func BenchmarkGenerateWithCap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		generateWithCap(1000000)
	}
}

func BenchmarkGenerate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		generate(1000000)
	}
}
```
运行结果：
```
➜  benchmarktest git:(master) ✗ go test -bench='Generate' -benchmem .
goos: darwin
goarch: arm64
pkg: benchmarktest
BenchmarkGenerateWithCap-8   	      62	  17984146 ns/op	 8003625 B/op	       1 allocs/op
BenchmarkGenerate-8          	      57	  19988124 ns/op	41678204 B/op	      39 allocs/op
PASS
ok  	benchmarktest	2.390s
```
Generate 分配的内存是 GenerateWithCap 的 6 倍，设置了切片容量，内存只分配一次，而不设置切片容量，内存分配了 39 次。
## 4、 重置计时器
### 4.1 ResetTimer()
```go
func BenchmarkFib(b *testing.B) {
	time.Sleep(time.Second * 3) // 模拟耗时准备任务
	b.ResetTimer() // 重置定时器
	for n := 0; n < b.N; n++ {
		fib(30) // run fib(30) b.N times
	}
}
```
### 4.2 StopTimer() & StartTimer()
```go
// sort_test.go
package main

import (
	"math/rand"
	"testing"
	"time"
)

func generateWithCap(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, 0, n)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}

func bubbleSort(nums []int) {
	for i := 0; i < len(nums); i++ {
		for j := 1; j < len(nums)-i; j++ {
			if nums[j] < nums[j-1] {
				nums[j], nums[j-1] = nums[j-1], nums[j]
			}
		}
	}
}

func BenchmarkBubbleSort(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		nums := generateWithCap(10000)
		b.StartTimer()
		bubbleSort(nums)
	}
}
```