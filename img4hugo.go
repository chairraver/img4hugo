package main // -*- coding: utf-8 -*-

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
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

		img, err := imaging.Open(file)
		if err != nil {
			log.Fatal(err)
		}

		for j := 0; j < len(imgsizes); j++ {
			resized := imaging.Resize(img, imgsizes[j], 0, imaging.Lanczos)
			rect := resized.Bounds().Max
			out := fmt.Sprintf("%s_%dx%d%s", strings.TrimSuffix(file, ext), rect.X, rect.Y, ext)
			err = imaging.Save(resized, out)
			log.Println("saved " + out)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

const template1 = `<div class="post-pic" data-src="{{.Original}}"{{if .Caption}} data-html-sub="{{.Caption}}"{{end}}>
    <img{{if .Caption}} alt="{{.Caption}}"{{end}} width="{{.Thumbwidth}}" height="{{.Thumbheight}}"
	src="{{.Thumbnail}}"/><br/>
{{if .Caption}}<p><em>{{.Caption}}</em></p>{{end}}
</div>
`

type HtmlImageProps struct {
	Original    string
	Caption     string
	Thumbnail   string
	Thumbwidth  int
	Thumbheight int
}

func tohtml(args []string) {

	// template := template.Must(template.New("imagediv").Parse(template1))

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

				img, err := imaging.Open(fullpath)
				if err != nil {
					log.Fatal(err)
				}

				width := img.Bounds().Max.X
				height := img.Bounds().Max.Y

				webpath := strings.Split(fullpath, sep+staticSplit+sep)[1]
				webpath = filepath.ToSlash(filepath.Clean("/" + webpath))

				fmt.Printf("{{< imgdiv class=\"%s\" href=\"%s\" alt=\"%s\"\n", "", webfullpath, "")
				fmt.Printf("    src=\"%s\" width=\"%d\" height=\"%d\" >}}\n", webpath, width, height)
				// picProps := HtmlImageProps{
				// 	Original:    webfullpath,
				// 	Caption:     "",
				// 	Thumbnail:   webpath,
				// 	Thumbwidth:  width,
				// 	Thumbheight: height,
				// }
				// err = template.Execute(os.Stdout, picProps)
				// if err != nil {
				// 	log.Println("executing template:", err)
				// }
			}
		}
	}
}
