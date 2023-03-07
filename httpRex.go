package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
	"github.com/spf13/pflag"
	"github.com/fatih/color"
	"bufio"
	"math/rand"
)

// import flag "github.com/spf13/pflag"

type Result struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
}

func getRandomColour() (*color.Color){

	rand.Seed(time.Now().UnixNano())
	colors := []*color.Color{
		color.New(color.FgRed),
		color.New(color.FgGreen),
		color.New(color.FgYellow),
		color.New(color.FgBlue),
		color.New(color.FgMagenta),
		color.New(color.FgCyan),
		color.New(color.FgWhite),
		color.New(color.FgHiRed),
		color.New(color.FgHiGreen),
		color.New(color.FgHiYellow),
		color.New(color.FgHiBlue),
		color.New(color.FgHiMagenta),
		color.New(color.FgHiCyan),
		color.New(color.FgHiWhite),
		color.New(color.Bold),
		color.New(color.Faint),
		color.New(color.Italic),
		color.New(color.Underline),
		color.New(color.BlinkSlow),
		color.New(color.BlinkRapid),
		color.New(color.ReverseVideo),
		color.New(color.Concealed),
		color.New(color.CrossedOut),
		color.New(color.FgHiBlack),
	}
	randomColor := colors[rand.Intn(len(colors))]
	return randomColor
}


func makeBanner() {

	col := getRandomColour()
	col.Println(`
██╗  ██╗████████╗████████╗██████╗ ██████╗ ███████╗██╗  ██╗
██║  ██║╚══██╔══╝╚══██╔══╝██╔══██╗██╔══██╗██╔════╝╚██╗██╔╝
███████║   ██║      ██║   ██████╔╝██████╔╝█████╗   ╚███╔╝ 
██╔══██║   ██║      ██║   ██╔═══╝ ██╔══██╗██╔══╝   ██╔██╗ 
██║  ██║   ██║      ██║   ██║     ██║  ██║███████╗██╔╝ ██╗
	`)

	col = getRandomColour()
	col.Println("\nAuthor: avik_saikat")
	col.Println("Github: https://github.com/aviksaikat")
	col.Println("Gitlab: https://gitlab.com/aviksaikat")
	fmt.Println("\n----------------------------------------------------------------\n")
}


func getStatusCode(url string) int {
	resp, err := http.Get(url)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

func checkUrls(urls []string, printStatusCode bool, saveFile string) {
	var wg sync.WaitGroup
	results := make([]Result, 0)
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			statusCode := getStatusCode(url)
			if statusCode != 0 {
				if printStatusCode {
					col := color.New(color.FgGreen)
					if statusCode >= 400 && statusCode < 500 {
						col = color.New(color.FgRed)
					} else if statusCode >= 300 && statusCode < 400 {
						col = color.New(color.FgYellow)
					}
					colouredStatusString := col.Sprintf("%d", statusCode)
					fmt.Printf("%s [%s]\n", url, colouredStatusString)
					// col.Printf(" [%d]\n", colouredStatusString)
					results = append(results, Result{URL: url, StatusCode: statusCode})
				} else if statusCode == 200 {
					fmt.Println(url)
					results = append(results, Result{URL: url, StatusCode: statusCode})
				}
			}
		}(url)
	}
	wg.Wait()

	if saveFile != "" {
		file, err := os.Create(saveFile)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		if ext := getFileExtension(saveFile); ext == "json" {
			jsonData, err := json.MarshalIndent(results, "", "  ")
			if err != nil {
				fmt.Println("Error encoding to JSON:", err)
				return
			}
			_, err = file.Write(jsonData)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		} else {
			for _, result := range results {
				_, err := file.WriteString(result.URL + "\n")
				if err != nil {
					fmt.Println("Error writing to file:", err)
					return
				}
			}
		}
	}
}

func getFileExtension(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i+1:]
		}
	}
	return ""
}

func main() {
	var urls []string
	var url string
	var printStatusCode bool
	var saveFile string
	var outputFileType string
	var urlFile string
	var banner bool
	var showHelp bool


	pflag.StringSliceVar(&urls, "l", nil, "URLs comma, space separated")
	pflag.StringVar(&url, "u", "", "Single URL to check")
	pflag.StringVar(&urlFile, "f", "", "File containing URLs to check")
	pflag.BoolVar(&printStatusCode, "status-code", false, "Print status code of each URL")
	pflag.StringVar(&saveFile, "o", "", "Save output to file")
	pflag.StringVar(&outputFileType, "file-type", "text", "Output file type (text or json)")
	pflag.BoolVar(&banner, "banner", false, "Print banner")
	pflag.BoolVar(&showHelp, "h", false, "show help")


	pflag.Parse()

	if showHelp {
		makeBanner()
		fmt.Println("Options:")
		pflag.PrintDefaults()
		os.Exit(0)
	}

	if banner != false {
		makeBanner()
		os.Exit(0)
	}

	makeBanner()

	if len(urls) == 0 && url == "" && urlFile == ""{
		fmt.Println("Please provide at least one URL to check")
		return
	}

	if url != "" {
		urls = append(urls, url)
	}
	if urlFile != "" {
		file, err := os.Open(urlFile)
		if err != nil {
			fmt.Printf("Error opening file %s: %s\n", urlFile, err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			urls = append(urls, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading file %s: %s\n", urlFile, err)
			return
		}
	}

	checkUrls(urls, printStatusCode, saveFile)

	if saveFile != "" {
		col := color.New(color.FgGreen)
		col.Printf("Results saved to %s\n", saveFile)
	}
}
