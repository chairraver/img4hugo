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
	imgsizes      = []int{1024, 640, 320, 160}
)

func main() {

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

			width, height, err := getImageData(output)
			if err != nil {
				log.Fatal("file " + output + " is not accessible")
			}

			dest := strings.TrimSuffix(file, ext)
			dest = dest + "_" + width + "x" + height + ext

			err = os.Rename(output, dest)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

const template1 = `<div class="post-pic" data-src="{{.Original}}">
    <img alt="" width="{{.Thumbwidth}}" height="{{.Thumbheight}}"
	src="{{.Thumbnail}}"></img><br/>
{{if .Caption}}<p>{{.Caption}}</p>{{end}}
</div>
`

type HtmlImageProps struct {
	Original    string
	Caption     string
	Thumbnail   string
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

		webfullpath := cwd + string(os.PathSeparator) + file
		webfullpath = strings.Split(webfullpath, sep+staticSplit+sep)[1]
		webfullpath = filepath.ToSlash(filepath.Clean("/" + webfullpath))

		direntries, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}

		for j := 0; j < len(direntries); j++ {
			name := direntries[j].Name()

			if strings.HasPrefix(name, base_noext+"_") {
				fullpath := cwd + string(os.PathSeparator) +
					dir + string(os.PathSeparator) + name
				width, height, err := getImageData(fullpath)
				if err != nil {
					log.Fatal(err)
				}

				webpath := strings.Split(fullpath, sep+staticSplit+sep)[1]
				webpath = filepath.ToSlash(filepath.Clean("/" + webpath))

				pic := HtmlImageProps{
					Original:    webfullpath,
					Caption:     "",
					Thumbnail:   webpath,
					Thumbwidth:  width,
					Thumbheight: height,
				}
				err = template.Execute(os.Stdout, pic)
				if err != nil {
					log.Println("executing template:", err)
				}
			}
		}
	}
}

func getImageData(img string) (width, height string, err error) {

	_, err = os.Stat(img)
	if err != nil {
		return
	}

	identifyCmd := exec.Command(imIdentifyCmd, img)

	identifyOut, err := identifyCmd.Output()
	if err != nil {
		return
	}

	res := strings.Split(string(identifyOut), " ")[2]
	width = strings.Split(res, "x")[0]
	height = strings.Split(res, "x")[1]
	return
}
