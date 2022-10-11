package main

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var emptyDirFiles []string

func findEmptyDirAndFile(dir string, del, includeDir, includeFile bool) error {
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		if includeDir && d.IsDir() {
			fileInfos, err := ioutil.ReadDir(path)
			if err != nil {
				return err
			}

			if len(fileInfos) == 0 {
				log.Println("[info] find idr", path, " is empty")
				emptyDirFiles = append(emptyDirFiles, path)
			}

			return nil
		}

		if includeFile && !info.IsDir() && info.Size() == 0 {
			log.Println("[info] find file", path, "size=0")
			emptyDirFiles = append(emptyDirFiles, path)
		}

		return nil
	})

	if del {
		for _, f := range emptyDirFiles {
			err = os.Remove(f)
			if err != nil {
				log.Printf("[err] remove %s fail ,err=%s", f, err.Error())
				continue
			}
			log.Printf("[info] remove %s success", f)
		}
	}
	return err
}
