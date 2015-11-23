package main // -*- coding: utf-8 -*-

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var (
	imIdentifyCmd = "identify"
	imConvertCmd  = "convert"
	staticSplit   = "static"
	imgsizes      []int
)

func main() {

	imgsizes = []int{1024, 640, 320, 160}

	var img4hugoRootCmd = &cobra.Command{
		Use:   "img4hugo",
		Short: "img4hugo is an application to simply embedding images into hugo content.",
	}

	var resizeCmd = &cobra.Command{
		Use:   "resize image",
		Short: "Resize the image to the standard set of image sizes",
		Run: func(cmd *cobra.Command, args []string) {
			resize(args, imgsizes)
		},
	}

	var tohtml = &cobra.Command{
		Use:   "tohtml image",
		Short: "Produce a short HTML fragment for inclusion into a hugo post",
		Run: func(cmd *cobra.Command, args []string) {
			tohtml(args)
		},
	}

	img4hugoRootCmd.AddCommand(resizeCmd)
	img4hugoRootCmd.AddCommand(tohtml)
	img4hugoRootCmd.Execute()
}

func resize(args []string, imgsizes []int) {
	for i := 0; i < len(args); i++ {

		file := args[i]
		ext := filepath.Ext(file)

		_, err := os.Stat(file)
		if err != nil {
			log.Fatal("file " + file + " is not accessible")
		}

		for j := 0; j < len(imgsizes); j++ {
			sizeStr := fmt.Sprintf("%d", imgsizes[j])
			output := "output_" + sizeStr + ext

			convertCmd := exec.Command(imConvertCmd, "-resize", sizeStr, file, output)
			err := convertCmd.Run()
			if err != nil {
				log.Fatal("converting " + file + " to " + output + " failed")
			}

			sizeStr, err = getImageData(output)
			if err != nil {
				log.Fatal("file " + output + " is not accessible")
			}

			dest := strings.TrimSuffix(file, ext)
			dest = dest + "_" + sizeStr + ext

			err = os.Rename(output, dest)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

const template1 = `
<div class="post-pic" data-src="{{.Original}}">
    <img alt="" width="{{.Thumbwidth}}" height="{{.Thumbheight}}"
	src="/images/2015/11/IMG_20150613_132225_640x725.jpg"></img><br/>
{{if .Caption}}<p>{{.Caption}}</p>{{end}}
</div>
`

type ImageProps struct {
	Original    string
	Caption     string
	Thumbwidth  string
	Thumbheight string
}

func tohtml(args []string) {

	template := template.Must(template.New("imagediv").Parse(template1))

	for i := 0; i < len(args); i++ {
		file := args[i]

		_, err := os.Stat(file)
		if err != nil {
			log.Fatal("file " + file + " is not accessible")
		}

		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		base := filepath.Base(file)
		ext := filepath.Ext(base)
		base_noext := strings.TrimSuffix(base, ext)
		dir := filepath.Dir(file)
		sep := string(os.PathSeparator)

		direntries, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}

		for j := 0; j < len(direntries); j++ {
			name := direntries[j].Name()

			if strings.HasPrefix(name, base_noext) {
				fullpath := cwd + string(os.PathSeparator) +
					dir + string(os.PathSeparator) + name
				staticpath := strings.Split(fullpath, sep+staticSplit+sep)
				fmt.Println(fullpath)
				fmt.Println(filepath.ToSlash(filepath.Clean("/" + staticpath[1])))
				getImageData(dir + string(os.PathSeparator) + name)
			}
		}
	}
}

func getImageData(img string) (sizeStr string, err error) {

	_, err = os.Stat(img)
	if err != nil {
		return
	}

	identifyCmd := exec.Command(imIdentifyCmd, img)

	identifyOut, err := identifyCmd.Output()
	if err != nil {
		return
	}

	res := strings.Split(string(identifyOut), " ")
	fmt.Printf("%q\n", res)

	return res[2], nil
}
