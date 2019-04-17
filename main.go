package main

import (
	"fmt"
	_ "github.com/oov/psd"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func getAllFile(pathname string, f func(string)) error {
	rd, err := ioutil.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			dir, _ := filepath.Abs(pathname + "/" + fi.Name() + "/")
			err := getAllFile(dir, f)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			dir, _ := filepath.Abs(pathname + "/" + fi.Name())
			f(dir)
		}
	}
	return err
}
func pds2png(pdsfile string, pngfile string) {
	file, err := os.Open(pdsfile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	out, err := os.Create(pngfile)
	if err != nil {
		panic(err)
	}
	err = png.Encode(out, img)
	if err != nil {
		panic(err)
	}
}
func main() {

	if len(os.Args) > 1 {
		dir, err := filepath.Abs(os.Args[1])
		if err != nil {
			log.Fatal(err.Error())
		}
		outdir, err := filepath.Abs(dir + "/psd2png")
		if err != nil {
			log.Fatal(err.Error())
		}
		_, err = os.Stat(outdir)
		if err != nil {
			err = os.MkdirAll(outdir, os.ModeDir)
			if err != nil {
				log.Fatal(err.Error())
			}
		}

		_ = getAllFile(dir, func(s string) {
			if path.Ext(s) == ".psd" {
				newdir := strings.ReplaceAll(s, dir, "")
				newdir, err := filepath.Abs(outdir + newdir)
				if err != nil {
					log.Fatal(err.Error())
				}
				_, err = os.Stat(filepath.Dir(newdir))
				if err != nil {
					err = os.MkdirAll(filepath.Dir(newdir), os.ModeDir)
					if err != nil {
						log.Fatal(err.Error())
					}
				}

				pds2png(s, newdir+".png")
				log.Println(s, "-->", newdir+".png")
			}
		})
		if runtime.GOOS == "windows" {
			exec.Command("explorer.exe", outdir).Start()
		}
	}
}
