// This code builds up maps of source code line numbers allowing
// code that has coverage statements to be compared with the original
// source code. This is necessary when correcting stack traces and build
// failure output generated from `go test -cover` in order to present it
// to the tester in the UI. See the tests for more details.
//
// The inspiration for this code:
//     https://groups.google.com/d/topic/golang-nuts/GBeHbrIOxU4/discussion
//
// Rob Pike's summary:
//     I think it's pretty easy, but comme ci comme ça.
//
// My apologies.
//
package watcher

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/smartystreets/goconvey/web/server/contract"
)

type SourceCache struct {
	fs        contract.FileSystem
	checksums map[string]int64       // [string] = package-relative path; int64 = checksum
	files     map[string]map[int]int // [string] = package-relative path; [int] = instrumented line; int = source line;
}

func (self *SourceCache) Update(path string, sum int64) {
	relativePath := contract.ResolvePackageName(path)
	if self.checksums[relativePath] == sum {
		return
	}
	log.Println("Caching source code for:", relativePath)
	self.checksums[relativePath] = sum
	log.Println("About to read:", path)
	original, coverage, err := self.fs.ReadGo(path)
	if err != nil {
		log.Println(err.Error())
		log.Println(original, coverage)
		return
	}

	rewrite := make(map[int]int)
	self.files[relativePath] = rewrite
	lines := strings.Split(original, "\n")
	covered := strings.Split(coverage, "\n")
	if strings.HasSuffix(relativePath, "/package_parser_go1.1_test.go") {
		log.Println(original)
		log.Println(coverage)
	}
	y := 0
	for x, line := range lines {
		if strings.Contains(covered[y], "GoConvey__coverage__") {
			y++
		} else if strings.HasPrefix(line, "// +build") {
			// go cover tool preserves build directive comments.
		} else if strings.HasPrefix(strings.TrimSpace(line), "//") {
			continue // go cover tool filters out all other comments.
		}

		rewrite[y] = x
		y++
	}
}

func (self *SourceCache) Rewrite(output string) string {
	preparedOutput := strings.Replace(output, "_test/", "", -1)
	outputLines := strings.Split(preparedOutput, "\n")
	for path, rewrite := range self.files {
		if strings.Contains(preparedOutput, path) {
			for number, line := range outputLines {
				if strings.Contains(line, path) {
					fields := strings.Split(line, ":")
					numberWord := strings.Split(strings.TrimSpace(fields[1]), " ")[0]
					referenceLineNumber, err := strconv.Atoi(numberWord)
					if err != nil {
						return output // just bail
					}
					prefix := ""
					if strings.HasPrefix(line, "\t") {
						prefix = "\t"
					}
					suffix := " "
					if len(fields) > 2 {
						suffix = ":" + fields[2]
					}
					correctLineNumber := rewrite[referenceLineNumber]
					correctedLine := fmt.Sprintf("%s%s:%d%s", prefix, path, correctLineNumber, suffix)
					outputLines[number] = correctedLine
				}
			}
		}
	}
	return strings.Join(outputLines, "\n")
}

func NewSourceCache(fs contract.FileSystem) *SourceCache {
	self := new(SourceCache)
	self.fs = fs
	self.checksums = make(map[string]int64)
	self.files = make(map[string]map[int]int)
	return self
}
