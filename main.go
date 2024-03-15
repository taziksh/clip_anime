package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector(colly.MaxDepth(4), colly.URLFilters(
		regexp.MustCompile(`https://fancaps\.net/anime/showimages\.php\?4880-Dragon_Ball_Z`),
		regexp.MustCompile(`https://fancaps\.net/anime/episodeimages\.php\?\d+-Dragon_Ball_Z/Episode_\d{1}`),
		regexp.MustCompile(`https://fancaps\.net/anime/picture\.php\?/\d+`),
		// regexp.MustCompile(`https://ancdn\.fancaps.net/\d+\.jpg`),
		// regexp.MustCompile(`https://cdni\.fancaps\.net/file/fancaps-animeimages/\d+\.jpg`),
	))

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		url := r.Request.URL.String()
		if strings.Contains(url, "https://fancaps.net/anime/picture.php?/") {

			lastSlashIndex := strings.LastIndex(url, "/")
			imageId := url[lastSlashIndex+1:]

			var b strings.Builder

			b.WriteString("https://cdni.fancaps.net/file/fancaps-animeimages/")
			b.WriteString(imageId)
			b.WriteString(".jpg")

			imageUrl := b.String()

			img, _ := os.Create(r.FileName() + ".jpg")
			defer img.Close()

			resp, _ := http.Get(imageUrl)
			defer resp.Body.Close()

			body, _ := io.Copy(img, resp.Body)
			fmt.Println("Body: ", body)
		}
		// fmt.Println("Visited", r.Body)

		// file, err := os.Create(r.FileName())
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer file.Close()

		// file.ReadFrom(bytes.NewReader(r.Body))

	})

	c.Visit("https://fancaps.net/anime/showimages.php?4880-Dragon_Ball_Z")
}
