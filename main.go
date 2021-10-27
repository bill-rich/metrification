package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type adMetric struct {
	userMetric    map[string]metric
	sessionMetric map[string]metric
	dateMetric    map[string]metric
}

type metric struct {
	impressions float32
	conversions float32
}

func processLine(key []string, adMetrics map[string]adMetric, lineRaw string) {
	data := strings.Split(lineRaw, "\t")
	if len(data) != len(key) {
		log.Println("mismatch in number of fields. expected: %d, got: %d", len(key), len(data))
	}

	row := map[string]string{}
	for i, value := range data {
		row[key[i]] = value
	}

	ad, ok := adMetrics[row["ad_id"]]
	if !ok {
		ad = adMetric{
			userMetric:    map[string]metric{},
			sessionMetric: map[string]metric{},
			dateMetric:    map[string]metric{},
		}
	}

	ad.addUserMetric(row)
	ad.addSessionMetric(row)
	ad.addDateMetric(row)

	adMetrics[row["ad_id"]] = ad

}

func (ad adMetric) addMetric(row map[string]string, colKey string) {
	adMetric, ok := ad.userMetric[row[colKey]]
	if !ok {
		adMetric = metric{
			impressions: 0,
			conversions: 0,
		}
	}
	switch row["event"] {
	case "view":
		adMetric.impressions++
	case "click":
		adMetric.conversions++
	}
	ad.userMetric[row[colKey]] = adMetric
}

func (ad adMetric) addSessionMetric(row map[string]string) {
	return
}

func (ad adMetric) addDateMetric(row map[string]string) {
	return
}

func main() {
	processFile("Sim.tsv")
}

func processFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("unable to read file %s: %s", fileName, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	ad := map[string]adMetric{}

	firstLine := true
	var key []string
	for scanner.Scan() {
		if firstLine {
			key = generateKey(scanner.Text())
			firstLine = false
		}
		processLine(key, ad, scanner.Text())
	}
	print(ad)
}

func print(adMetrics map[string]adMetric) {
	for key, ad := range adMetrics {
		fmt.Printf("Ad ID: %s\n", key)
		fmt.Printf("-User Metrics\n")
		for userId, user := range ad.userMetric {
			fmt.Printf("--UserID: %s\n", userId)
			fmt.Printf("---Impressions: %f\n", user.impressions)
			fmt.Printf("---Conversions: %f\n", user.conversions)
			conversionRate := float32(1.0)
			if user.impressions > 0 {
				conversionRate = user.conversions / user.impressions
			}
			fmt.Printf("---Conversion Ratio: %f\n", conversionRate)

		}
	}
}

func generateKey(header string) []string {
	key := []string{}
	keys := strings.Split(header, "\t")
	for _, keyString := range keys {
		key = append(key, keyString)
	}
	return key
}
