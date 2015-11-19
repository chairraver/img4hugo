package main // -*- coding: utf-8 -*-

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	imIdentifyCmd = "identify"
	imConvertCmd  = "convert"
	imgsize       [4]int
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
			resize(args)
		},
	}

	var toHtml = &cobra.Command{
		Use:   "tohtml image",
		Short: "Produce a short HTML fragment for inclusion into a hugo post",
		Run: func(cmd *cobra.Command, args []string) {
			toHtml(args)
		},
	}

	img4hugoRootCmd.AddCommand(resizeCmd)
	img4hugoRootCmd.AddCommand(toHtml)
	img4hugoRootCmd.Execute()
}

func resize(args []string) {
	for i := 0; i < len(args); i++ {

		_, err := os.Stat(args[i])
		if err != nil {
			log.Fatal("file " + args[i] + " is not accessible")
		}

		convertCmd := exec.Command(imConvertCmd, args[i])
	}
}

func toHtml(args []string) {
	for i := 0; i < len(args); i++ {

		_, err := os.Stat(args[i])
		if err != nil {
			log.Fatal("file " + args[i] + " is not accessible")
		}

		identifyCmd := exec.Command(imIdentifyCmd, args[i])

		identifyOut, err := identifyCmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(identifyOut))
	}

}
