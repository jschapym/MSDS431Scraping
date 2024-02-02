package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"sync"

	"github.com/gocolly/colly"
)

const (
	pagesDirectory = "./Pages"
	outputFile     = "./WikiScrapeOutput.jl"
)

type JSONoutput struct {
	Url   string `json:"url"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

func Scraping(urlstr string) (JSONoutput, []byte, error) {
	var jo JSONoutput
	var htmlbody []byte

	urlstring := urlstr
	jo.Url = urlstring

	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)
	c.OnResponse(func(r *colly.Response) {
		log.Println("Visited", r.StatusCode)
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", 0, err)
	})
	c.OnResponse(func(r *colly.Response) {
		htmlbody = r.Body
	})
	c.OnHTML("title", func(e *colly.HTMLElement) {
		jo.Title = e.Text
	})
	c.OnHTML("body", func(e *colly.HTMLElement) {
		jo.Text = e.Text
	})
	err := c.Visit(urlstring)
	if err != nil {
		return jo, nil, err
	}

	return jo, htmlbody, nil
}

func createDirectoryIfNotExists(directory string) error {
	_, err := os.Stat(directory)
	if os.IsNotExist(err) {
		err := os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directory: %s", err)
		}
	}
	return nil
}

func writeHTMLToFile(fileName string, content []byte) error {
	err := os.WriteFile(fileName, content, 0644)
	if err != nil {
		return err
	}
	return nil
}

func writeJSONToFile(fileName string, jo JSONoutput) error {
	output, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer output.Close()

	jsonData, err := json.Marshal(jo)
	if err != nil {
		return err
	}

	_, err = output.WriteString(string(jsonData) + "\n")
	if err != nil {
		return err
	}
	return nil
}

func processURL(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Scraping %s\n", url)
	jo, hb, err := Scraping(url)
	if err != nil {
		log.Printf("Error scraping %s: %v\n", url, err)
		return
	}

	fileName := path.Base(url)
	htmlFilePath := fmt.Sprintf("%s/%s.html", pagesDirectory, fileName)

	err = createDirectoryIfNotExists(pagesDirectory)
	if err != nil {
		log.Printf("Error creating directory: %v\n", err)
		return
	}

	err = writeHTMLToFile(htmlFilePath, hb)
	if err != nil {
		log.Printf("Error writing HTML to file for %s: %v\n", url, err)
		return
	}

	err = writeJSONToFile(outputFile, jo)
	if err != nil {
		log.Printf("Error writing JSON to file for %s: %v\n", url, err)
		return
	}
}

func main() {
	var wg sync.WaitGroup

	urls := []string{ // Add the missing assignment statement
		"https://en.wikipedia.org/wiki/Robotics",
		"https://en.wikipedia.org/wiki/Robot",
		"https://en.wikipedia.org/wiki/Reinforcement_learning",
		"https://en.wikipedia.org/wiki/Robot_Operating_System",
		"https://en.wikipedia.org/wiki/Intelligent_agent",
		"https://en.wikipedia.org/wiki/Software_agent",
		"https://en.wikipedia.org/wiki/Robotic_process_automation",
		"https://en.wikipedia.org/wiki/Chatbot",
		"https://en.wikipedia.org/wiki/Applications_of_artificial_intelligence",
		"https://en.wikipedia.org/wiki/Android_(robot)",
	}

	for _, nexturl := range urls {
		wg.Add(1)
		go processURL(nexturl, &wg)
	}

	wg.Wait()
}
