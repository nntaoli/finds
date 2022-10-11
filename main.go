package main

import (
	"errors"
	"flag"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
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

func findAndReplaceAll(old, new string, replace bool) {
	oldReg, _ := regexp.Compile(old)
	for _, f := range files {
		data, err := ioutil.ReadFile(f)
		if err != nil {
			log.Println("[err] read  ", f, " err => ", err)
			continue
		}
		dataStr := string(data)
		if replace {
			log.Println("[info] replace ", f)
			newData := oldReg.ReplaceAllString(dataStr, new)
			err = ioutil.WriteFile(f, []byte(newData), 644)
			if err != nil {
				log.Println("[err] write ", f, " error ", err)
			}
		} else {
			if oldReg.MatchString(dataStr) {
				indes := oldReg.FindAllStringIndex(dataStr, len(old))
				for _, idx := range indes {
					if f == "/home/dev/data2/workspace/find-replace/find-replace" {
						log.Println(dataStr[idx[0]-12 : idx[1]+11])
					}

				}
				log.Println("[info] find ", f)
			}
		}
	}
}

func readAllFileName(dir string) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, f := range fileInfos {
		if f.IsDir() {
			allFiles(dir + "/" + f.Name())
			continue
		}
		files = append(files, path.Join(dir, f.Name()))
	}
}

func batchRename(old, new string) {
	for _, f := range files {
		dir, name := filepath.Split(f)
		newName := strings.ReplaceAll(name, old, new)
		err := os.Rename(f, filepath.Join(dir, newName))
		if err != nil {
			log.Println("[error] rename  [", f, "] . err=", err.Error())
			continue
		}
		log.Println("[info] rename success new :", filepath.Join(dir, newName))
	}
}

func main() {
	app := &cli.App{
		Name: "find-replace",
		Commands: []*cli.Command{
			{
				Name:  "content",
				Usage: "replace content",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "replace",
						Usage: "true: find and replace \n false: only find",
					},
					&cli.StringFlag{
						Name:     "old",
						Usage:    "old string",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "new",
						Usage: "new string",
					},
				},
				Action: func(context *cli.Context) error {
					if context.NArg() == 0 {
						return errors.New("please input search dir path")
					}
					allFiles(context.Args().Get(0))
					findAndReplaceAll(context.String("old"), context.String("new"), context.Bool("replace"))
					return nil
				},
			},
			{
				Name:  "rename",
				Usage: "rename file name",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "old",
						Usage:    "old string",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "new",
						Usage: "new string",
					},
				},
				Action: func(context *cli.Context) error {
					if context.NArg() == 0 {
						return errors.New("please input search dir path")
					}
					allFiles(context.Args().Get(0))
					batchRename(context.String("old"), context.String("new"))
					return nil
				},
			},
			{
				Name:  "empty",
				Usage: "find and delete empty dir or file",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "delete", Usage: "delete empty dir or file"},
					&cli.BoolFlag{Name: "dir", Usage: "include dir", Value: true},
					&cli.BoolFlag{Name: "file", Usage: "include file", Value: true},
				},
				Action: func(context *cli.Context) error {
					if context.Args().Len() == 0 {
						return errors.New("no  directory")
					}
					return findEmptyDirAndFile(context.Args().Get(0), context.Bool("delete"), context.Bool("dir"), context.Bool("file"))
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
