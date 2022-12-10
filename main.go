package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type GraphData struct {
	DataDown[] string `yaml:"datadown"`
	DataUp[] string `yaml:"dataup"`
	TotalCount int64 `yaml:"total_count"`
	BandwidthLimit int64 `yaml:"bwlimit"`
	Label string `yaml:"label"`
	Title string `yaml:"title"`
}

const (
	SUBSTR_START = "graph_data = {"
	SUBSTR_END = "}"
)

func asJSON(graphData GraphData) {
	_, err: = json.Marshal(graphData)
	if err != nil {
		fmt.Println(err)
		return
	}

	pretty, err: = json.MarshalIndent(graphData, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(pretty))
}

func asCSV(graphData GraphData) {
	writer: = csv.NewWriter(os.Stdout)
	writer.Write([] string {
		"Title", "Label", "Time", "DataDown", "DataUp", "TotalCount", "BandwidthLimit"
	})
	totalCountString: = strconv.FormatInt(graphData.TotalCount, 10)
	bandwidthLimitString: = strconv.FormatInt(graphData.BandwidthLimit, 10)
	timer: = time.Time {}
	for index: = 0;index < len(graphData.DataDown);index++{
		writer.Write([] string {
			graphData.Title,
				graphData.Label,
				timer.Format("15:04"),
				graphData.DataDown[index],
				graphData.DataUp[index],
				totalCountString,
				bandwidthLimitString,
		})
		timer = timer.Add(time.Minute)
	}

	writer.Flush()
}

func scrape(dataSet string, dateString string, formatString string) {
	url: = fmt.Sprintf("https://residential.launtel.net.au/traffic-graphs/%s/%s", dataSet, dateString)
	res,
	err: = http.Get(url)
	if err != nil {
		fmt.Println("Unable to fetch data from Launtel, sorry. Are you online?")
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Printf("HTTP error when fetching data: %d %s\n", res.StatusCode, res.Status)
		return
	}

	docBuffer: = new(strings.Builder)
	io.Copy(docBuffer, res.Body)

	rx: = regexp.MustCompile(`(?s)` +
		regexp.QuoteMeta(SUBSTR_START) + `(.*?)` +
		regexp.QuoteMeta(SUBSTR_END))
	match: = rx.FindStringSubmatch(docBuffer.String())
	result: = match[1]

		result = strings.ReplaceAll(result, "\t", " ")
	result = strings.ReplaceAll(result, SUBSTR_START, "")
	result = strings.ReplaceAll(result, ",\n", "\n")

	var graphData GraphData
	if err: = yaml.Unmarshal([] byte(result), & graphData);err != nil {
		fmt.Println("The data we got from Launtel was malformed. They may have changed the format, or your data set + data combination aren't working.")
		return
	}

	if formatString == "json" {
		asJSON(graphData)
	} else {
		asCSV(graphData)
	}
}

func main() {
	now: = time.Now()
	var dataSet string
	var dateString string
	var formatString string

	flag.StringVar( & dataSet, "s", "", "Data-set to fetch (eg CVC000000623965)")
	flag.StringVar( & dateString, "d", now.Format("2006-01-02"), "Date to pull data for the given data set on (eg 2022-12-11)")
	flag.StringVar( & formatString, "f", "csv", "Format the output in 'json' or 'csv'")
	flag.Parse()

	dataSet = strings.TrimSpace(dataSet)
	dateString = strings.TrimSpace(dateString)
	formatString = strings.TrimSpace(formatString)

	if dataSet == "" {
		fmt.Println("You must specify a data set with -s (such as CVC000000623965). Use -h for help.")
		return
	}

	if formatString != "csv" && formatString != "json" {
		fmt.Println("-f (format) must be either 'csv' or 'json'. Use -h for help.")
		return
	}

	_,
	err: = time.Parse("2006-01-02", dateString)
	if err != nil {
		fmt.Println("Date (-d) is invalid. Try again. Use -h for help.")
		return
	}

	scrape(dataSet, dateString, formatString)
}
