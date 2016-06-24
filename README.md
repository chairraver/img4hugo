# img4hugo

`ìmg4hugo` is a companion program to my [Hugo](http://gohugo.io) theme
variation
[hugo-uno-chairraver](https://github.com/chairraver/hugo-uno-chairraver),
based on [hugo-uno](https://github.com/SenjinDarashiva/hugo-uno). It
is also a first step to learn the [Go](https://golang.org) programming
language. I've been using this program on Linux and Windows.

As the name implies, it's a helper program for the handling of images
in connection with Hugo. The program serves 3 purposes:

* Resize the master image. You will rarely need the images in full
  resolution, except perhaps on a retina Mac. Therefore I'm typically
  resizing my images so that they will fill a typical HD 1080p display.
* Create a set of thumbnail images for the use in responsive image set.
* Create small text fragments for the individual images, which can
  then be pasted into the Hugo post or page, which is currently being
  written.

At this point (2015-12-19) it appears to be complete enough. Therefore
I'll be using it for a bit with my Hugo site. At a later time I plan
to extend it with text expansions for responsive images (`picture` and
`srcset`).

# Usage

`img4hugo` supports the three subcommands `size`, `thumbs` and `tohtml`.

``` bash
$ img4hugo help
img4hugo is an application to simplify the embedding of images into hugo content.

Usage:
  img4hugo [command]

Available Commands:
  size        Resize the max. resolution image [1920 1080]
  thumbs      Create thumbnails for the image with a standard set of image sizes [1024 640 320]
  tohtml      Produce a short HTML fragment for inclusion into a hugo post

Use "img4hugo [command] --help" for more information about a command.
```

Besides the configuration through the command line options `img4hugo`
can be configured by including the appropriate parameters in the main
Hugo config file. `img4hugo` will look into the current directory as
well up to 4 directory level above the current directory.

The following lines show `img4hugo`s default configuration as it would
appear in the Hugo configuration file. The macros for the `tohtml`
option are implemented through standard Go templates.

``` bash
[img4hugo]
  size = [ "1920", "1080" ]
  thumbs = [ "1024", "640", "320" ]
  tohtml = ['''
{{`{{<`}} imgdiv class="{{.Class}}" href="{{.Fullresimg}}" alt="{{.Caption}}"
    src="{{.Thumbnailimg}}" width="{{.Width}}" height="{{.Height}}" {{`>}}`}}
''', '''
{{`{{<`}} img id="" class="{{.Class}}" href="{{.Fullresimg}}" alt="{{.Caption}}"
    src="{{.Thumbnailimg}}" width="{{.Width}}" height="{{.Height}}" {{`>}}`}}
'''
]
```

* **Class** is directly passed through from the `tohtml` `-l` option.
* **Caption** the same it true for the  `tohtml` `-c` option.
* **Fullresimg** is the web path to the full resolution image.
* **Thumbnailimg** is the web path to the thumbnail image.
* **Width** and **Height** are directly extracted from the thumbnail
  image. 

### The `size` subcommand

The `size` subcommand is used to resize the original or master
image. The thumbnails are derived from this image. The original image
is renamed with the `.org` extension and the new rescaled image is
written to a file with the original image name. If a file already
exists with the `.org` extension, the program assumes, that this was
renamed by an earlier execution and that this it is the image with the
original proportions. Therefore the image with the `.org` extension
will be used as the source for creating the newly resized image.

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

Interesting tidbit if you look at the above listing is, that the file
witout the `.org` is actually the file with the smaller
dimensions. The `imaging` library, which is used for the image
manipulation saves the JPG format with a hard coded quality of 95%.

### The `thumbs` subcommand

The `thumbs` subcommand creates a number of thumbnails for the image
argument. By default 3 images are created with the horizontal
dimensions 1024, 640 and 320 pixels. The newly created thumbnail
images are saved with file name reflecting their respective image
dimensions.

``` bash
$ img4hugo thumbs -h
Create thumbnails for the image with a standard set of image sizes [1024 640 320]

Usage:
  img4hugo thumbs <image file> [flags]

Flags:
  -s, --size="1024,640,320": specifiy new list of thumbnail image sizes
```

The files created from a sample run would look like the following.

``` bash
$ dir IMG_20150531_124021*

Mode                LastWriteTime     Length Name
----                -------------     ------ ----
-a---        25.11.2015     13:58     349000 IMG_20150531_124021_1024x729.jpg
-a---        25.11.2015     13:58      44588 IMG_20150531_124021_320x228.jpg
-a---        25.11.2015     13:58     148356 IMG_20150531_124021_640x456.jpg
```

The optional `-s` or `--size` flag might be used to specify a list of
different image sizes. The comma separated list of sizes must be
passed to the program as one argument. So you might want to enclose
the size list in quotes depending on your shell. `img4hugo` itself is
prepared to discard any whitespace characters around the comma.

### The `tohtml` subcommand

This subcommand outputs a short Go template Hugo shortcode for
each thumbnail belonging to the master image. At this point you would
have to copy and paste the shortcode for the particular thumbnail,
that you want to embed into you Hugo post. Additionally the output is
currently hardcoded into to the program code. Later I intend to
implement some more variability for the shortcode. 

``` bash
$ img4hugo tohtml -h
Produce a short HTML fragment for inclusion into a hugo post

Usage:
  img4hugo tohtml image [flags]

Flags:
  -c, --caption="": caption text for the image
  -l, --class="": additional css class for the image
  -n, --noerrors[=false]: do not warn about location
  -t, --template=0: # of template to use
```

The `-c` or `--caption` option allows the specification of a caption
for the image.

The `-l` or `--class` option allows the definition of an additional
CSS class for the `div`.

The `-n` or `--noerrors` is intended for the case, if the program is
executed outside of Hugo site or content area. `img4hugo` calculates
the absolute paths for the image `href`s by looking at the path
component `/static/`. Typically any static data (images, CSS, etc) for
a Hugo site is located below the `static` directory. An error message
is generated, when no `static` path component is found in the path of
the images. The `-n` option suppressed this error message.

With the `-t` option the corresponding template is
selected. `img4hugo` contains two templates by default (see the usage
section above the respective definitions).

``` bash
$ img4hugo tohtml -l "floatright" -c "Nice shoot." IMG_20150613_132225.jpg
{{< imgdiv class="floatright" href="/images/2015/11/IMG_20150613_132225.jpg" alt="Nice shoot."
    src="/images/2015/11/IMG_20150613_132225_1024x1159.jpg" width="1024" height="1159" >}}
{{< imgdiv class="floatright" href="/images/2015/11/IMG_20150613_132225.jpg" alt="Nice shoot."
    src="/images/2015/11/IMG_20150613_132225_320x362.jpg" width="320" height="362" >}}
{{< imgdiv class="floatright" href="/images/2015/11/IMG_20150613_132225.jpg" alt="Nice shoot."
    src="/images/2015/11/IMG_20150613_132225_640x725.jpg" width="640" height="725" >}}
```

The Hugo shortcode `imgdiv` from above is defined in my theme
variation
[hugo-uno-chairraver](https://github.com/chairraver/hugo-uno-chairraver)
and is intended to be used with
[jQuery lightgallery](https://sachinchoolur.github.io/lightGallery/)
(note the `data-src` and `data-html-sub` attributes for the `div`
tag).

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

## License

Licensed under the [MIT License](LICENSE.txt).
