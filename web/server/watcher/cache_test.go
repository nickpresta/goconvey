package watcher

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/goconvey/web/server/system"
)

func TestCache(t *testing.T) {
	var (
		cache *SourceCache
		fs    *system.FakeFileSystem
	)

	Convey("Subject: SourceCache", t, func() {
		fs = system.NewFakeFileSystem()
		cache = NewSourceCache(fs)

		Convey("When build failure output contains a reference to a file that has been scanned", func() {
			fs.Populate(actualFilePath, originalBuildFailureSource, coveredBuildFailureSource, nil)
			cache.Update(actualFilePath, 42)
			actual := cache.Rewrite(coverageBuildFailure)

			Convey("The cache should correct the line numbers with the original line numbers", func() {
				So(actual, ShouldEqual, expectedBuildFailure)
			})
		})

		Convey("When build failure output contains no references to known files", func() {
			actual := cache.Rewrite("blah blah blah")

			Convey("It should return the same stack trace", func() {
				So(actual, ShouldEqual, "blah blah blah")
			})
		})

		Convey("When panic output contains a reference to a file that has been scanned", func() {
			fs.Populate(actualFilePath, originalPanicSource, coveredPanicSource, nil)
			cache.Update(actualFilePath, 42)
			actual := cache.Rewrite(coveragePanic)

			Convey("The cache should correct the line numbers with the original line numbers", func() {
				So(actual, ShouldEqual, expectedPanic)
			})
		})

		Convey("When panic output contains no references to known files", func() {
			SkipConvey("It should return the same stack trace", func() {

			})
		})

	})
}

const actualFilePath = "/root/gopath/src/github.com/smartystreets/goconvey/scratch/scratch.go"

const originalBuildFailureSource = `package scratch

func hi() {
	a := []int{}
	// fmt.Println(a[123])
}
`

const coveredBuildFailureSource = `package scratch

func hi() {
	GoConvey__coverage__.Count[0] = 1
	a := []int{}
	// fmt.Println(a[123])
}

var GoConvey__coverage__ = struct {
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

const expectedBuildFailure = `# github.com/smartystreets/goconvey/scratch
github.com/smartystreets/goconvey/scratch/scratch.go:4: a declared and not used
FAIL	github.com/smartystreets/goconvey/scratch [build failed]
`

const originalPanicSource = `package scratch

import "fmt"

func hi() {
	a := []int{}
	fmt.Println(a[123])
}
`

const coveredPanicSource = `package scratch

import "fmt"

func hi() {
	GoConvey__coverage__.Count[0] = 1
	a := []int{}
	fmt.Println(a[123])
}

var GoConvey__coverage__ = struct {
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

const coveragePanic = `panic: runtime error: index out of range [recovered]
	panic: runtime error: index out of range

goroutine 3 [running]:
runtime.panic(0x00, 0x00)
	/usr/local/go/src/pkg/runtime/panic.c:266 +0x00
testing.func·005()
	/usr/local/go/src/pkg/testing/testing.go:385 +0x00
runtime.panic(0x00, 0x00)
	/usr/local/go/src/pkg/runtime/panic.c:248 +0x00
github.com/smartystreets/goconvey/scratch.hi()
	github.com/smartystreets/goconvey/scratch/_test/scratch.go:8 +0x00
github.com/smartystreets/goconvey/scratch.TestScratch(0x00)
	/Users/mike/work/go/src/github.com/smartystreets/goconvey/scratch/scratch_test.go:6 +0x00
testing.tRunner(0x00, 0x00)
	/usr/local/go/src/pkg/testing/testing.go:391 +0x00
created by testing.RunTests
	/usr/local/go/src/pkg/testing/testing.go:471 +0x00

goroutine 1 [chan receive]:
testing.RunTests(0x00, 0x00, 0x00, 0x00, 0x00)
	/usr/local/go/src/pkg/testing/testing.go:472 +0x00
testing.Main(0x00, 0x00, 0x00, 0x00, 0x00, ...)
	/usr/local/go/src/pkg/testing/testing.go:403 +0x00
main.main()
	github.com/smartystreets/goconvey/scratch/_testmain.go:47 +0x00
`

const expectedPanic = `panic: runtime error: index out of range [recovered]
	panic: runtime error: index out of range

goroutine 3 [running]:
runtime.panic(0x00, 0x00)
	/usr/local/go/src/pkg/runtime/panic.c:266 +0x00
testing.func·005()
	/usr/local/go/src/pkg/testing/testing.go:385 +0x00
runtime.panic(0x00, 0x00)
	/usr/local/go/src/pkg/runtime/panic.c:248 +0x00
github.com/smartystreets/goconvey/scratch.hi()
	github.com/smartystreets/goconvey/scratch/scratch.go:7 
github.com/smartystreets/goconvey/scratch.TestScratch(0x00)
	/Users/mike/work/go/src/github.com/smartystreets/goconvey/scratch/scratch_test.go:6 +0x00
testing.tRunner(0x00, 0x00)
	/usr/local/go/src/pkg/testing/testing.go:391 +0x00
created by testing.RunTests
	/usr/local/go/src/pkg/testing/testing.go:471 +0x00

goroutine 1 [chan receive]:
testing.RunTests(0x00, 0x00, 0x00, 0x00, 0x00)
	/usr/local/go/src/pkg/testing/testing.go:472 +0x00
testing.Main(0x00, 0x00, 0x00, 0x00, 0x00, ...)
	/usr/local/go/src/pkg/testing/testing.go:403 +0x00
main.main()
	github.com/smartystreets/goconvey/scratch/_testmain.go:47 +0x00
`
