package system

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type FileSystem struct {
	gobin    string
	coverage bool
}

func (self *FileSystem) Walk(root string, step filepath.WalkFunc) {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if isMetaDirectory(info) {
			return filepath.SkipDir
		}

		return step(path, info, err)
	})

	if err != nil {
		log.Println("Error while walking file system:", err)
		panic(err)
	}
}
func isMetaDirectory(info os.FileInfo) bool {
	name := info.Name()
	return info.IsDir() && (strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") || name == "testdata")
}

func (self *FileSystem) Listing(directory string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(directory)
}

func (self *FileSystem) Exists(directory string) bool {
	info, err := os.Stat(directory)
	return err == nil && info.IsDir()
}

func (self *FileSystem) ReadGo(path string) (source, covered string, err error) {
	rawSource, err := ioutil.ReadFile(path)
	if err != nil {
		return "", "", err
	}
	source = string(rawSource)

	if !self.coverage {
		covered = source
		return
	}

	command := exec.Command(self.gobin, "tool", "cover", "-mode=set", "-var=GoConvey__coverage__", path)
	rawOutput, err := command.CombinedOutput()
	covered = string(rawOutput)
	if err != nil {
		source = ""
		covered = ""
	}
	return
}

func NewFileSystem(gobin string) *FileSystem {
	self := new(FileSystem)
	self.gobin = gobin
	self.coverage = goVersion_1_2_orGreater()
	return self
}
