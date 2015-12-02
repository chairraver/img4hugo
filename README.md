# img4hugo

`ìmg4hugo` is a companion program to my [Hugo](http://gohugo.io) theme
variation
[hugo-uno-chairraver](https://github.com/chairraver/hugo-uno-chairraver),
based on [hugo-uno](https://github.com/SenjinDarashiva/hugo-uno). It
is also a first step to learn the [Go](https://golang.org) programming
language.

As the name implies, it's a helper program for the handling of images
in connection with Hugo. The program serves 3 purposes:

* Resize the master image. You will rarely need the images in full
  resolution, except perhaps on a retina Mac. Therefore I'm typically
  resizing my images so that they will fill a typical HD 1080p display.
* Create a set of thumbnail images for the use in responsive image set.
* Create small text fragments for the individual images, which can
  then be pasted into the Hugo post or page, which is currently being
  written.

# Usage

``` bash
$ img4hugo.exe help
img4hugo is an application to simplyfy the embedding of images into hugo content.

Usage:
  img4hugo [command] <image file(s) ...>

Available Commands:
  thumbs      Create thumbnails for the image with a standard set of image sizes ([1024 640 320])
  size        Resize the max. resolution image ([1920 1080])
  tohtml      Produce a short HTML fragment for inclusion into a hugo post

Use "img4hugo [command] --help" for more information about a command.
```


