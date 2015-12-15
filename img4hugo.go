package main // -*- coding: utf-8 -*-

import (
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/disintegration/imaging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	imIdentifyCmd   = "identify"
	imConvertCmd    = "convert"
	staticSplit     = "static"
	stdsize         = []int{1920, 1080}
	stdsizecfg      = "img4hugo.size"
	noxyswap        bool
	imgsizes        = []int{1024, 640, 320}
	imgsizescfg     = "img4hugo.thumbs"
	tohtmltemplates []*template.Template
	tohtmlcfg       = "img4hugo.tohtml"
	newDefaultSize  string
	newThumbsSizes  string
	caption         string
	class           string
	tplidx          int
	noerrors        bool
)

func main() {

	configure()

	var img4hugoRootCmd = &cobra.Command{
		Use:   "img4hugo",
		Short: "img4hugo is an application to simplify the embedding of images into hugo content.",
	}

	var defaultSizeCmd = &cobra.Command{
		Use:   "size image(s)",
		Short: "Resize the max. resolution image " + fmt.Sprint(stdsize),
		Run: func(cmd *cobra.Command, args []string) {
			defaultSize(args, stdsize, noxyswap)
		},
	}
	defaultSizeCmd.Flags().StringVarP(&newDefaultSize, "size", "s", "1920,1080", "specifiy new default image size x,y")
	defaultSizeCmd.Flags().BoolVarP(&noxyswap, "noxyswap", "n", false, "do not scale relative to longest side")

	var thumbsCmd = &cobra.Command{
		Use:   "thumbs image",
		Short: "Create thumbnails for the image with a standard set of image sizes " + fmt.Sprint(imgsizes),
		Run: func(cmd *cobra.Command, args []string) {
			thumbs(args, imgsizes)
		},
	}
	thumbsCmd.Flags().StringVarP(&newThumbsSizes, "size", "s", "1024,640,320", "specifiy new list of thumbnail image sizes")

	var tohtml = &cobra.Command{
		Use:   "tohtml image",
		Short: "Produce a short HTML fragment for inclusion into a hugo post",
		Run: func(cmd *cobra.Command, args []string) {
			tohtml(args, noerrors)
		},
	}
	tohtml.Flags().StringVarP(&caption, "caption", "c", "", "caption text for the image")
	tohtml.Flags().StringVarP(&class, "class", "l", "", "additional css class for the image")
	tohtml.Flags().IntVarP(&tplidx, "template", "t", 0, "# of template to use")
	tohtml.Flags().BoolVarP(&noerrors, "noerrors", "n", false, "do not warn about location")

	img4hugoRootCmd.AddCommand(defaultSizeCmd)
	img4hugoRootCmd.AddCommand(thumbsCmd)
	img4hugoRootCmd.AddCommand(tohtml)
	img4hugoRootCmd.Execute()
}

func configure() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../..")
	viper.AddConfigPath("../../..")
	viper.AddConfigPath("../../../..")
	err := viper.ReadInConfig()
	if err != nil {
		log.Print(fmt.Errorf("configuration error: %s\n", err))
	}
	log.Print("using config file " + viper.ConfigFileUsed() + "\n")

	if viper.IsSet(stdsizecfg) {
		vals := viper.GetStringSlice(stdsizecfg)
		if len(vals) == 1 {
			num, err := strconv.Atoi(strings.TrimSpace(vals[0]))
			if err != nil {
				log.Fatal(err)
			}
			stdsize[0] = num
			stdsize[1] = num
		}
		if len(vals) == 2 {
			num, err := strconv.Atoi(strings.TrimSpace(vals[0]))
			if err != nil {
				log.Fatal(err)
			}
			stdsize[0] = num

			num, err = strconv.Atoi(strings.TrimSpace(vals[1]))
			if err != nil {
				log.Fatal(err)
			}
			stdsize[1] = num
		}
	}
	if viper.IsSet(imgsizescfg) {
		vals := viper.GetStringSlice(imgsizescfg)
		imgsizes = make([]int, len(vals))
		for i := 0; i < len(vals); i++ {
			num, err := strconv.Atoi(strings.TrimSpace(vals[i]))
			if err != nil {
				log.Fatal(err)
			}
			imgsizes[i] = num
		}
	}
	if viper.IsSet("img4hugo.test") {
		log.Print("config file sets test " +
			fmt.Sprint(viper.GetStringSlice("img4hugo.test")) + "\n")
	}
	if viper.IsSet(tohtmlcfg) {
		log.Print("config file sets template " + viper.GetString(tohtmlcfg) + "\n")
	}
}

func defaultSize(args []string, stdsize []int, noxyswap bool) {
	if newDefaultSize != "" {
		vals := strings.Split(newDefaultSize, ",")
		for i := 0; i < len(vals); i++ {
			numstr := strings.TrimSpace(vals[i])
			if numstr == "" {
				continue
			}
			num, err := strconv.Atoi(numstr)
			if err != nil {
				log.Print(err)
				continue
			}
			stdsize[i] = num
		}
	}

	for i := 0; i < len(args); i++ {

		var img image.Image

		orgext := ".org"
		file := args[i]

		_, err := os.Stat(file + orgext)
		// err == nil means file is already present and has already
		// been resize in which case we abort.
		if err == nil {
			log.Print(file + orgext + " exists; has apparently already been resized")
			log.Print("using " + file + orgext + " as source")

			img, err = imaging.Open(file + orgext)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			_, err := os.Stat(file)
			if err != nil {
				log.Fatal("file " + file + " is not accessible")
			}

			img, err = imaging.Open(file)
			if err != nil {
				log.Fatal(err)
			}
		}

		var resized image.Image
		if noxyswap || (img.Bounds().Max.X >= img.Bounds().Max.Y) {
			resized = imaging.Resize(img, stdsize[0], 0, imaging.Lanczos)
		} else {
			resized = imaging.Resize(img, 0, stdsize[1], imaging.Lanczos)
		}
		_, err = os.Stat(file + orgext)
		// err == nil means file is already present and has already
		// been resize in which case we abort.
		if err != nil {
			err = os.Rename(file, file+orgext)
			if err != nil {
				log.Fatal(err)
			}
		}
		err = imaging.Save(resized, file)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func thumbs(args []string, imgsizes []int) {
	if newThumbsSizes != "" {
		vals := strings.Split(newDefaultSize, ",")
		imgsizes := make([]int, 5, 10)
		for i := 0; i < len(vals); i++ {
			numstr := strings.TrimSpace(vals[i])
			if numstr == "" {
				imgsizes[i] = 0
				continue
			}
			num, err := strconv.Atoi(numstr)
			if err != nil {
				log.Print(err)
				continue
			}
			imgsizes[i] = num
		}
	}

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
			if imgsizes[j] == 0 {
				continue
			}
			resized := imaging.Resize(img, imgsizes[j], 0, imaging.Lanczos)
			rect := resized.Bounds().Max
			out := fmt.Sprintf("%s_%dx%d%s",
				strings.TrimSuffix(file, ext), rect.X, rect.Y, ext)
			err = imaging.Save(resized, out)
			log.Println("saved " + out)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func tohtml(args []string, noerrors bool) {

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
		if strings.Contains(webfullpath, sep+staticSplit+sep) {
			webfullpath = strings.Split(webfullpath, sep+staticSplit+sep)[1]
		} else {
			if !noerrors {
				log.Print("not within your Hugo directory structure")
			}
		}
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

				var webpath string
				if strings.Contains(fullpath, sep+staticSplit+sep) {
					webpath = strings.Split(fullpath, sep+staticSplit+sep)[1]
				} else {
					if !noerrors {
						log.Print("not within your Hugo directory structure")
					}
					webpath = fullpath
				}
				// webpath := strings.Split(fullpath, sep+staticSplit+sep)[1]
				webpath = filepath.ToSlash(filepath.Clean("/" + webpath))

				fmt.Printf("{{< imgdiv class=\"%s\" href=\"%s\" alt=\"%s\"\n",
					class, webfullpath, caption)
				fmt.Printf("    src=\"%s\" width=\"%d\" height=\"%d\" >}}\n",
					webpath, width, height)
			}
		}
	}
}
