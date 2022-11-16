package scripts

import "os"

type File struct {
	Key      string  `json:"key"`
	Parent   string  `json:"parent"`
	Title    string  `json:"title"`
	Type     string  `json:"type"`
	IsLeaf   bool    `json:"is_leaf"`
	Children []*File `json:"children"`
}

func GetFiles() []*File {
	var files []*File
	dir, err := os.ReadDir("data/scripts/")
	if err != nil {
		return []*File{}
	}
	for _, entry := range dir {
		if entry.IsDir() {
			f := &File{
				Key:      entry.Name(),
				Parent:   "",
				Title:    entry.Name(),
				Type:     "directory",
				IsLeaf:   true,
				Children: []*File{},
			}
			twoDir, err := os.ReadDir("data/scripts/" + entry.Name())
			if err != nil {
				continue
			}
			for _, dirEntry := range twoDir {
				f.Children = append(f.Children, &File{
					Key:    entry.Name() + "/" + dirEntry.Name(),
					Parent: entry.Name(),
					Title:  dirEntry.Name(),
					Type: func() string {
						if dirEntry.IsDir() {
							return "directory"
						} else {
							return "file"
						}
					}(),
					IsLeaf:   true,
					Children: []*File{},
				})
			}
			files = append(files, f)

		} else {
			files = append(files, &File{
				Key:      entry.Name(),
				Parent:   "",
				Title:    entry.Name(),
				Type:     "file",
				IsLeaf:   true,
				Children: []*File{},
			})
		}
	}
	return files
}
