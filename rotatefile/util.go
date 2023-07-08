package rotatefile

import (
	"compress/gzip"
	"io"
	"io/fs"
	"os"

	"github.com/gookit/goutil/fsutil"
)

const compressSuffix = ".gz"

func compressFile(srcPath, dstPath string) error {
	srcFile, err := os.OpenFile(srcPath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// create and open a gz file
	gzFile, err := fsutil.QuickOpenFile(dstPath)
	if err != nil {
		return err
	}
	defer gzFile.Close()

	zw := gzip.NewWriter(gzFile)

	fileSt, err := srcFile.Stat()
	if err != nil {
		return err
	}

	zw.Name = fileSt.Name()
	zw.ModTime = fileSt.ModTime()

	if _, err = io.Copy(zw, srcFile); err != nil {
		return err
	}
	return zw.Close()
}

type filterFunc func(fPath string, fi os.FileInfo) bool
type handleFunc func(fPath string, fi os.FileInfo) error

// from the go pkg: path/filepath.glob()
// TODO replace use fsutil.FindInDir()
func findFilesInDir(dir string, handleFn handleFunc, filters ...filterFunc) (err error) {
	fi, err := os.Stat(dir)
	if err != nil {
		return // ignore I/O error
	}
	if !fi.IsDir() {
		return // ignore I/O error
	}

	d, err := os.Open(dir)
	if err != nil {
		return // ignore I/O error
	}
	defer d.Close()

	// names, _ := d.Readdirnames(-1)
	// sort.Strings(names)

	stats, _ := d.Readdir(-1)
	for _, fi := range stats {
		baseName := fi.Name()
		filePath := dir + baseName

		// call filters
		if len(filters) > 0 {
			var filtered = false
			for _, filter := range filters {
				if !filter(filePath, fi) {
					filtered = true
					break
				}
			}

			if filtered {
				continue
			}
		}

		if err := handleFn(filePath, fi); err != nil {
			return err
		}
	}
	return nil
}

// TODO replace to fsutil.FileInfo
type fileInfo struct {
	fs.FileInfo
	filePath string
}

// Path get file full path. eg: "/path/to/file.go"
func (fi *fileInfo) Path() string {
	return fi.filePath
}

func newFileInfo(filePath string, fi fs.FileInfo) fileInfo {
	return fileInfo{filePath: filePath, FileInfo: fi}
}

// modTimeFInfos sorts by oldest time modified in the fileInfo.
// eg: [old_220211, old_220212, old_220213]
type modTimeFInfos []fileInfo

// Less check
func (fis modTimeFInfos) Less(i, j int) bool {
	return fis[j].ModTime().After(fis[i].ModTime())
}

// Swap value
func (fis modTimeFInfos) Swap(i, j int) {
	fis[i], fis[j] = fis[j], fis[i]
}

// Len get
func (fis modTimeFInfos) Len() int {
	return len(fis)
}
