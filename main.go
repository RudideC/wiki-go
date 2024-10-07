package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/RudideC/wiki-go/utils"
	"github.com/orus-dev/osui"
	"github.com/orus-dev/osui/ui"
)

const baseURL = "https://www.wikipedia.org"

type LinkData struct {
	Title string
	URL   string
}

func main() {
	if len(os.Args) < 2 {
		utils.HelpMessage()
		os.Exit(0)
	}
	switch strings.ToLower(os.Args[1]) {
	case "-s":
		fallthrough
	case "--search":
		searchTerm := strings.Join(os.Args[2:], "_")
		url := baseURL + "/wiki/" + searchTerm
		if isList(url) {
			readSite(ListResults(url)[0].URL, getPageTitle(ListResults(url)[0].URL))
		} else {
			readSite(url, getPageTitle(url))
		}
	case "-h":
		fallthrough
	case "--help":
		utils.HelpMessage()
		os.Exit(0)
	case "-v":
		fallthrough
	case "--version":
		utils.VersionMessage()
		os.Exit(0)
	default:
		searchTerm := strings.Join(os.Args[1:], "_")
		url := baseURL + "/wiki/" + searchTerm
		if isList(url) {
			resultMenu(ListResults(url))
		} else {
			readSite(url, getPageTitle(url))
		}
	}
}

func getPageTitle(url string) string {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("Error: Status code %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	title := doc.Find("title").Text()
	cleanTitle := strings.TrimSuffix(title, " - Wikipedia")
	return cleanTitle
}

func isList(url string) bool {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("Error: Status code %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	found := false
	content := doc.Find("div.mw-content-ltr.mw-parser-output")
	content.Find("p").Each(func(index int, element *goquery.Selection) {
		text := element.Text()
		if strings.Contains(text, "may refer to") {
			found = true
		}
	})
	return found
}

func ListResults(url string) []LinkData {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("Error: Status code %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var linkDataList []LinkData
	content := doc.Find("div.mw-content-ltr.mw-parser-output")

	content.Find("ul").Each(func(i int, ul *goquery.Selection) {
		ul.Find("li").Each(func(j int, li *goquery.Selection) {
			li.Find("a").Each(func(k int, a *goquery.Selection) {
				href, exists := a.Attr("href")
				if exists {
					fullURL := baseURL + href
					parts := strings.Split(href, "/")
					title := parts[len(parts)-1]
					linkDataList = append(linkDataList, LinkData{
						Title: title,
						URL:   fullURL,
					})
				}
			})
		})
	})
	return linkDataList
}

func readSite(url string, title string) {
	utils.ResetStyles()
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Error: Status code %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var paragraphs string
	paragraphs = title
	content := doc.Find("div.mw-content-ltr.mw-parser-output")
	content.Find("p").Each(func(index int, element *goquery.Selection) {
		text := element.Text()
		paragraphs += fmt.Sprintf("%s\n", text)
	})
	utils.Clear()
	cmd := exec.Command("less")
	cmd.Stdout = os.Stdout

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	_, err = stdin.Write([]byte(paragraphs))
	if err != nil {
		log.Fatal(err)
	}
	stdin.Close()

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

}

func resultMenu(linkDataList []LinkData) {
	count := len(linkDataList)
	if count > 10 {
		count = 10
	}
	items := make([]string, count)
	for i := 0; i < count; i++ {
		items[i] = linkDataList[i].Title
	}
	app := ui.Menu(items...).Params(ui.MenuParams{
		OnSelected: func(mc *ui.MenuComponent, b bool) {
			item := mc.Items[mc.SelectedItem]
			readSite(baseURL+"/wiki/"+item, item)
		},
	})
	screen := osui.NewScreen(app)
	screen.Run()
}
