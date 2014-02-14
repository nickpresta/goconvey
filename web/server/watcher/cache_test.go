package watcher

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCache(t *testing.T) {
	Convey("Subject: SourceCache", t, func() {

		Convey("When a stack trace contains a reference to a file has been scanned", func() {
			SkipConvey("The cache should correct the line numbers with the original line numbers", func() {

			})
		})

		Convey("When a stack trace contains no references to known files", func() {
			SkipConvey("It should return the same stack trace", func() {

			})
		})

	})
}

const originalSource = `package scratch

import "fmt"

func hi() {
	a := []int{}
	fmt.Println(a[123])
}
`

const instrumentedSource = `package scratch

import "fmt"

func hi() {
	GoCover.Count[0] = 1
	a := []int{}
	fmt.Println(a[123])
}

var GoCover = struct {
	Count     [1]uint32
	Pos       [3 * 1]uint32
	NumStmt   [1]uint16
} {
	Pos: [3 * 1]uint32{
		5, 8, 0x2000b, // [0]
	},
	NumStmt: [1]uint16{
		2, // 0
	},
}
`
const coverageBuildFailure = `# github.com/smartystreets/goconvey/scratch
/var/folders/rw/z5hccl114lzbvjjrqfwr9mcr0000gn/T/go-build011729042/github.com/smartystreets/goconvey/scratch/_test/scratch.go:5: a declared and not used
FAIL	github.com/smartystreets/goconvey/scratch [build failed]
`

const coveragePanic = `panic: runtime error: index out of range [recovered]
	panic: runtime error: index out of range

goroutine 3 [running]:
runtime.panic(0xf5bc0, 0x23ff97)
	/usr/local/go/src/pkg/runtime/panic.c:266 +0xb6
testing.func·005()
	/usr/local/go/src/pkg/testing/testing.go:385 +0xe8
runtime.panic(0xf5bc0, 0x23ff97)
	/usr/local/go/src/pkg/runtime/panic.c:248 +0x106
github.com/smartystreets/goconvey/scratch.hi()
	github.com/smartystreets/goconvey/scratch/_test/scratch.go:8 +0xd0
github.com/smartystreets/goconvey/scratch.TestStuff(0x210333000)
	/Users/mike/work/go/src/github.com/smartystreets/goconvey/scratch/scratch_test.go:6 +0x1a
testing.tRunner(0x210333000, 0x23b080)
	/usr/local/go/src/pkg/testing/testing.go:391 +0x8b
created by testing.RunTests
	/usr/local/go/src/pkg/testing/testing.go:471 +0x8b2

goroutine 1 [chan receive]:
testing.RunTests(0x1508f8, 0x23b080, 0x1, 0x1, 0xca201)
	/usr/local/go/src/pkg/testing/testing.go:472 +0x8d5
testing.Main(0x1508f8, 0x23b080, 0x1, 0x1, 0x242c80, ...)
	/usr/local/go/src/pkg/testing/testing.go:403 +0x84
main.main()
	github.com/smartystreets/goconvey/scratch/_test/_testmain.go:93 +0x11b
`
