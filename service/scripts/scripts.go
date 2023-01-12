package scripts

import (
	"bytes"
	"context"
	"fmt"
	"github.com/huoxue1/qinglong-go/service/client"
	"github.com/huoxue1/qinglong-go/service/env"
	"github.com/huoxue1/qinglong-go/utils"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
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

func Run(filePath, content string) error {
	err := os.WriteFile(path.Join("data", "scripts", filePath), []byte(content), 0666)
	if err != nil {
		return err
	}
	cmd := getCmd(filePath)
	cancelChan := make(chan int, 1)
	ctx := context.WithValue(context.Background(), "cancel", cancelChan)
	now := time.Now()
	utils.RunTask(ctx, cmd, env.GetALlEnv(), func(ctx context.Context) {
		writer := ctx.Value("log").(io.Writer)
		_, _ = writer.Write([]byte(fmt.Sprintf("##开始执行..  %s\n\n", now.Format("2006-01-02 15:04:05"))))
	}, func(ctx context.Context) {
		writer := ctx.Value("log").(io.Writer)
		_, _ = writer.Write([]byte(fmt.Sprintf("\n##执行结束..  %s，耗时%.1f秒\n\n", time.Now().Format("2006-01-02 15:04:05"), time.Now().Sub(now).Seconds())))
		_ = os.Remove(filePath)
	}, client.MyClient)
	return nil
}

func getCmd(filePath string) string {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".js":
		return fmt.Sprintf("%s %s", "node", filePath)
	case ".py":
		return fmt.Sprintf("%s %s", "python", filePath)

	}
	return ""
}

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
			if strings.HasPrefix(entry.Name(), "_") {
				continue
			}
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
