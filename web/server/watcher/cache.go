package watcher

import "github.com/smartystreets/goconvey/web/server/contract"

type SourceCache struct {
	fs        contract.FileSystem
	checksums map[string]int64       // [string] = package-relative path; int64 = checksum
	files     map[string]map[int]int // [string] = package-relative path; [int] = instrumented line; int = source line;
}

func (self *SourceCache) Update(path string, sum int64) {
	// relativePath := contract.ResolvePackageName(path)
	// if self.checksums[relativePath] != sum {
	// 	self.checksums[relativePath] = sum
	// 	// TODO: reload file, which means we discard the old map and create a new one from scratch.
	// }
}

func (self *SourceCache) Rewrite(stack string) string {
	// for path, _ /*rewrite*/ := range self.files {
	// 	if strings.Contains(stack, path) {
	// 	}
	// }
	return ""
}
