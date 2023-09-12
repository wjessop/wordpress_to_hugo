# A small wordpress.com to Hugo conversion program

1. To use, download the XML and rar files from wordpress.com and unpack/unzip them somewhere on your local machine
2. Clone this repo and install the deps
3. Edit the code to point to the output, XML file and source image directories
4. Run `go run main.go`

You should end up with a directory structure that looks like this containing markdown files and images per article:

```
BaseDir/
  content/
    posts/
      Some-post-title/
        index.md
        ImagesDir/
          post-image-1.jpg
          post-image-2.png
      A-different-post-title/
        index.md
        ImagesDir/
          post-image-for-this-other-arcicle.jpg
```

[I wrote more details, and caveats etc. on my blog](https://willj.net/posts/converting-a-wordpress.com_dump_to_hugo).
