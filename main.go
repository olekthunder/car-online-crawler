package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

const (
	configFileName = "config.yaml"
	url            = "http://api.car-online.ru/v2"
	dateFormat     = "2006-01-02"
)

var client = &http.Client{}

type appConfig struct {
	ApiToken string `yaml:"api_token"`
	Filename string `yaml:"file_to_save"`
	DateFrom string `yaml:"date_from"`
	DateTo   string `yaml:"date_to"`
	Timezone string `yaml:"timezone"`
}

type parseResult struct {
	date    time.Time
	mileage float64
}

type response struct {
	Mileage float64 `json:"mileage"`
}

type timeStampMilli struct {
	Time time.Time
}

func (t timeStampMilli) daysBefore(days int) timeStampMilli {
	t.Time = t.Time.Add(time.Duration(-days*24) * time.Hour)
	return t
}

func (t timeStampMilli) daysAfter(days int) timeStampMilli {
	t.Time = t.Time.Add(time.Duration(days*24) * time.Hour)
	return t
}

func (t timeStampMilli) String() string {
	return strconv.Itoa(int(t.Time.UnixNano() / int64(time.Millisecond)))
}

func (t timeStampMilli) now() timeStampMilli {
	t.Time = time.Now()
	return t
}

func (t timeStampMilli) fromTime(newTime time.Time) timeStampMilli {
	t.Time = newTime
	return t
}

func (t timeStampMilli) toTime() time.Time {
	return t.Time
}

func newTimestampMilli() *timeStampMilli {
	t := new(timeStampMilli)
	t.Time = time.Now()
	return t
}

func setParam(r *http.Request, key string, value string) {
	query := r.URL.Query()
	query.Set(key, value)
	r.URL.RawQuery = query.Encode()
}

func getMileage(client *http.Client, apiToken string, begin timeStampMilli, end timeStampMilli, c chan parseResult) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	setParam(request, "skey", apiToken)
	setParam(request, "content", "json")
	setParam(request, "begin", begin.String())
	setParam(request, "end", end.String())
	setParam(request, "get", "telemetry")

	resp, err := client.Do(request)
	if err != nil {
		log.Printf("[ERROR] %s\n", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	jsonResp := new(response)
	err = json.Unmarshal(body, jsonResp)
	if err != nil {
		panic(err)
	}
	c <- parseResult{date: end.Time, mileage: jsonResp.Mileage / 1000}
	fmt.Printf("[DONE] %s\n", end.Time.Format(dateFormat))
}

func saveToXlsx(filename string, res *map[time.Time]float64) error {
	// sort map
	var keys []time.Time
	for k := range *res {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].Before(keys[j]) })

	xlsx := excelize.NewFile()

	for i, t := range keys {
		xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(i+1), t.Format(dateFormat))
		xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(i+1), (*res)[t])
	}
	err := xlsx.SaveAs(filename)
	return err
}

func main() {
	var err error

	config := new(appConfig)
	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Fatalf("Failed to open file %s", configFileName)
	}
	configData, _ := ioutil.ReadAll(configFile)
	err = yaml.Unmarshal(configData, config)
	if err != nil {
		log.Fatalln("Failed to read the config file")
	}

	loc, err := time.LoadLocation(config.Timezone)
	if err != nil {
		log.Fatalf("Unknown timezone: %s", config.Timezone)
	}
	startDate, err := time.Parse(dateFormat, config.DateFrom)
	if err != nil {
		log.Fatalln("Failed to parse start date")
	}
	endDate, err := time.Parse(dateFormat, config.DateTo)
	if err != nil {
		log.Fatalln("Failed to parse end date")
	}
	startDate = startDate.In(loc)
	endDate = endDate.In(loc)
	days := int(math.Ceil(endDate.Sub(startDate).Hours() / 24))

	begin := newTimestampMilli().fromTime(startDate)
	end := begin.daysAfter(1)
	c := make(chan parseResult)

	// Sigterm channel
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	for i := 1; i < days; i++ {
		go getMileage(client, config.ApiToken, begin, end, c)
		begin, end = begin.daysAfter(1), end.daysAfter(1)
	}

	result := make(map[time.Time]float64)
L:
	for i := 1; i < days; i++ {
		select {
		case pr := <-c:
			result[pr.date] = pr.mileage
		case <-stop:
			fmt.Println("\nStopping...")
			break L
		}
	}

	if len(result) != 0 {
		err = saveToXlsx(config.Filename, &result)
		if err != nil {
			log.Fatalf("Failed to save file: %s", config.Filename)
		}
	} else {
		fmt.Println("No data to save")
		os.Exit(1)
	}
}
