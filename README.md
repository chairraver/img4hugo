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
$ img4hugo help
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

The `size` subcommand is used to resize the original or master
image. The thumbnails are derived from this image. The original image
is renamed with the `.org` extension and the new rescaled image is
written to a file with the original image name. If a file already
exists with the `.org` extension, the program terminates. You have to
manually rename the backup file to the original file name to repeat
the rescaling process.

``` bash
$ img4hugo size -h
Resize the max. resolution image [1920 1080]

Usage:
  img4hugo size <image file> [flags]

Flags:
  -n, --noxyswap[=false]: do not scale relative to longest side
  -s, --size="1920,1080": specifiy new default image size x,y
```

By default, the program uses standard HD resolution 1920x1080 as the
default size for the rescaled image. Additionally the default
behaviour is to rescale vertically to 1080 if the vertical dimension
of the image is larger than the horizontal dimension. This behaviour
can be changed by supplying the `-n` or `-noxyswap` flag. Then the
image is always resized to the default horizontal dimension.

The `-s` or `--size` flag can be used to specify different horizontal
and vertical dimensions.



``` bash
$ dir IMG_20150531_124021*

Mode                LastWriteTime     Length Name
----                -------------     ------ ----
-a---        25.11.2015     13:58    1054685 IMG_20150531_124021.jpg
-a---        16.11.2015     09:19     952557 IMG_20150531_124021.jpg.org
```

### The `thumbs` subcommand

``` bash
$ img4hugo thumbs -h
Create thumbnails for the image with a standard set of image sizes [1024 640 320]

Usage:
  img4hugo thumbs <image file> [flags]

Flags:
  -s, --size="1024,640,320": specifiy new list of thumbnail image sizes
```

``` bash
$ dir IMG_20150531_124021*

Mode                LastWriteTime     Length Name
----                -------------     ------ ----
-a---        25.11.2015     13:58     349000 IMG_20150531_124021_1024x729.jpg
-a---        25.11.2015     13:58      44588 IMG_20150531_124021_320x228.jpg
-a---        25.11.2015     13:58     148356 IMG_20150531_124021_640x456.jpg
```

### The `tohtml` subcommand


``` bash
$ img4hugo tohtml -h
Produce a short HTML fragment for inclusion into a hugo post

Usage:
  img4hugo tohtml image [flags]

Flags:
  -c, --caption="": caption text for the image
  -l, --class="": additional css class for the image
```

``` bash
$ img4hugo tohtml -l "floatright" -c "Nice shoot." IMG_20150613_132225.jpg
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
