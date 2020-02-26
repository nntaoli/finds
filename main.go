package main

import (
	"flag"
	"io/ioutil"
	"log"
	"regexp"
)

var (
	old           = flag.String("old", "", "old char")
	new           = flag.String("new", "", "new char")
	replaceAllCmd = flag.Bool("replaceAll", false, "replace")

	files []string
)

func allFiles(dir string) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, f := range fileInfos {
		if f.IsDir() {
			allFiles(dir + "/" + f.Name())
			continue
		}
		files = append(files, dir+"/"+f.Name())
	}
}

func findAndReplaceAll(files []string) {
	oldReg, _ := regexp.Compile(*old)
	for _, f := range files {
		data, err := ioutil.ReadFile(f)
		if err != nil {
			log.Println("[err] read  ", f, " err => ", err)
			continue
		}

		if *replaceAllCmd {
			log.Println("[info] replace ", f)
			newData := oldReg.ReplaceAllString(string(data), *new)
			err = ioutil.WriteFile(f, []byte(newData), 644)
			if err != nil {
				log.Println("[err] write ", f, " error ", err)
			}
		} else {
			if oldReg.MatchString(*old) {
				log.Println("[info] find ", f)
			}
		}
	}
}

func main() {
	flag.Parse()
	allFiles(flag.Arg(len(flag.Args()) - 1))
	findAndReplaceAll(files)
}
