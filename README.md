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

`img4hugo` supports the three subcommands `size`, `thumbs` and `tohtml`.

``` bash
$ img4hugo.exe help
img4hugo is an application to simplyfy the embedding of images into hugo content.

Usage:
  img4hugo [command] <image file(s) ...>

Available Commands:
  size        Resize the max. resolution image [1920 1080]
  thumbs      Create thumbnails for the image with a standard set of image sizes [1024 640 320]
  tohtml      Produce a short HTML fragment for inclusion into a hugo post

Use "img4hugo [command] --help" for more information about a command.
```

### The `size` subcommand

``` bash
$ img4hugo.exe size -h
Resize the max. resolution image [1920 1080]

Usage:
  img4hugo size <image file> [flags]

Flags:
  -n, --noxyswap[=false]: don't scale relative to longest side
  -s, --size="1920,1080": specifiy new default image size x,y
```

### The `thumbs` subcommand

``` bash
$ img4hugo.exe thumbs -h
Create thumbnails for the image with a standard set of image sizes [1024 640 320]

Usage:
  img4hugo thumbs <image file> [flags]

Flags:
  -s, --size="1024,640,320": specifiy new list of thumbnail image sizes
```

### The `tohtml` subcommand



``` bash
$ img4hugo.exe tohtml -h
Produce a short HTML fragment for inclusion into a hugo post

Usage:
  img4hugo tohtml image [flags]

Flags:
  -c, --caption="": caption text for the image
  -l, --class="": additional css class for the image
```

``` bash
$ img4hugo.exe tohtml -l "floatright" -c "Nice shoot." .\IMG_20150613_132225.jpg
{{< imgdiv class="floatright" href="/images/2015/11/IMG_20150613_132225.jpg" alt="Nice shoot."
    src="/images/2015/11/IMG_20150613_132225_1024x1159.jpg" width="1024" height="1159" >}}
{{< imgdiv class="floatright" href="/images/2015/11/IMG_20150613_132225.jpg" alt="Nice shoot."
    src="/images/2015/11/IMG_20150613_132225_320x362.jpg" width="320" height="362" >}}
{{< imgdiv class="floatright" href="/images/2015/11/IMG_20150613_132225.jpg" alt="Nice shoot."
    src="/images/2015/11/IMG_20150613_132225_640x725.jpg" width="640" height="725" >}}
```

``` html
<div class="post-pic {{ .Get "class"}}"
     data-src="{{ .Get "href" }}"
     {{if .Get "alt"}}data-html-sub="{{ .Get "alt"}}"{{end}}>
  <img {{if .Get "alt"}}alt="{{.Get "alt"}}"{{end}}
       width="{{ .Get "width"}}"
       height="{{ .Get "height"}}"
       src="{{ .Get "src"}}"/><br/>
{{if .Get "alt"}}<p><em>{{.Get "alt"}}</em></p>{{end}}
</div>
```
