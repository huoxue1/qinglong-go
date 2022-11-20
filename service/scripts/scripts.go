package scripts

import (
	"bytes"
	"github.com/huoxue1/qinglong-go/utils"
	"os"
	"path"
	"sort"
)

type File struct {
	Key      string  `json:"key"`
	Parent   string  `json:"parent"`
	Title    string  `json:"title"`
	Type     string  `json:"type"`
	IsLeaf   bool    `json:"is_leaf"`
	Children []*File `json:"children"`
}

var (
	excludedFiles = []string{
		"node_modules",
		"__pycache__",
	}
)

func GetFiles(base, p string) []*File {
	var files Files
	dir, err := os.ReadDir(path.Join(base, p))
	if err != nil {
		return []*File{}
	}
	for _, entry := range dir {
		if utils.In(entry.Name(), excludedFiles) {
			continue
		}
		if entry.IsDir() {
			f := &File{
				Key:      path.Join(p, entry.Name()),
				Parent:   p,
				Title:    entry.Name(),
				Type:     "directory",
				IsLeaf:   true,
				Children: GetFiles(base, path.Join(p, entry.Name())),
			}
			files = append(files, f)

		} else {
			files = append(files, &File{
				Key:      path.Join(p, entry.Name()),
				Parent:   p,
				Title:    entry.Name(),
				Type:     "file",
				IsLeaf:   true,
				Children: []*File{},
			})
		}
	}
	sort.Sort(files)
	return files
}

type Files []*File

func (a Files) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a Files) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}
func (a Files) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	if a[i].Type != a[j].Type {
		if a[i].Type == "file" {
			return false
		} else {
			return true
		}
	} else {
		return bytes.Compare([]byte(a[i].Title), []byte(a[j].Title)) > 0
	}
}
