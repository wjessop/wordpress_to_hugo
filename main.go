package main

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/beevik/etree"
)

const (
	/*
		BaseDir is where you want the resulting articles to go. The data will be in a structure like this (subdir names
		based on the consts defined):

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

	*/
	BaseDir    = "/some/path/to/the/output/"
	ContentDir = "content"
	PostsDir   = "posts"
	ImagesDir  = "images"

	// The path to the exported XML file containing all posts
	WordPressXMLFile = "/path/to/yoursite.wordpress.com-2023-09-11-23_09_49/willwrites.wordpress.2023-09-11.000.xml"

	// The path to the images export dir
	LocalImageDir = "/path/to/media-export-31097282-from-0-to-1426"
)

func main() {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(WordPressXMLFile); err != nil {
		panic(err)
	}

	for _, item := range doc.FindElements("//item") {
		postType := item.SelectElement("wp:post_type")
		if postType.Text() == "attachment" {
			continue
		}

		title := item.SelectElement("title").Text()
		content := item.SelectElement("content:encoded").Text()
		dateStr := item.SelectElement("pubDate").Text()
		tags := extractTags(item)

		postDir := filepath.Join(BaseDir, ContentDir, PostsDir, strings.ReplaceAll(title, " ", "-"))
		_ = os.MkdirAll(filepath.Join(postDir, ImagesDir), os.ModePerm)

		converter := md.NewConverter("", true, nil)

		markdown, err := converter.ConvertString(content)
		if err != nil {
			panic(err)
		}

		frontMatter := fmt.Sprintf("---\ntitle: %s\ndate: %s\ndraft: true\ntags: [%s]\nsummary: \ncategory: \"\"\ntype: Post\n---\n", title, formatDate(dateStr), strings.Join(tags, ", "))
		contentString := frontMatter + string(markdown)
		err = os.WriteFile(filepath.Join(postDir, "index.md"), []byte(contentString), 0644)
		if err != nil {
			panic(err)
		}

		// Copy images
		for _, imgURL := range extractImageURLs(content) {
			parsedImgURL, _ := url.Parse(imgURL)
			pathParts := strings.Split(parsedImgURL.Path, "/")
			if len(pathParts) >= 4 {
				year, month, imgName := pathParts[len(pathParts)-3], pathParts[len(pathParts)-2], pathParts[len(pathParts)-1]
				srcPath := filepath.Join(LocalImageDir, year, month, imgName)
				destPath := filepath.Join(postDir, ImagesDir, imgName)
				err := copyFile(srcPath, destPath)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func extractTags(item *etree.Element) []string {
	var tags []string
	for _, category := range item.SelectElements("category") {
		if domain := category.SelectAttrValue("domain", ""); domain == "post_tag" {
			tags = append(tags, category.Text())
		}
	}
	return tags
}

func extractImageURLs(content string) []string {
	var urls []string
	decoder := xml.NewDecoder(strings.NewReader(content))
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "img" {
				for _, attr := range se.Attr {
					if attr.Name.Local == "src" {
						urls = append(urls, attr.Value)
					}
				}
			}
		}
	}
	return urls
}

func copyFile(src, dest string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	err = os.WriteFile(dest, input, 0644)
	if err != nil {
		return err
	}
	return nil
}

func formatDate(date string) string {
	t, err := time.Parse(time.RFC1123Z, date)
	if err != nil {
		panic(err)
	}
	return t.Format(time.RFC3339)
}
