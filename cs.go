package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	flag.Parse()
	searchCsv(flag.Arg(0), flag.Arg(1))
}

func searchCsv(uri, searchWord string) error {
	if strings.HasPrefix(uri, "http") {
		return searchCsvFromHttp(uri, searchWord)
	} else {
		return searchCsvFromFile(uri, searchWord)
	}
}

func searchCsvFromHttp(uri, searchWord string) error {
	// Handle input file
	response, err := http.Get(uri)
	if err != nil {
		fmt.Println("http.Get() error:", err)
		return err
	}
	defer response.Body.Close()

	// Search CSV and output hit lines to stdout
	reader := csv.NewReader(response.Body)
	reader.FieldsPerRecord = -1 // Accepts reading irregular csv format
	for {
		words, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("File read error:", err)
			return err
		}
		for _, w := range words {
			if strings.Contains(w, "404: Not Found") {
				fmt.Println(w)
				return errors.New("404: Not Found error")
			} else if strings.Contains(w, searchWord) {
				fmt.Println(words)
				break
			}
		}
	}
	return nil
}

func searchCsvFromFile(uri, searchWord string) error {
	// Handle input file
	fp, err := os.Open(uri)
	if err != nil {
		fmt.Println("File open error:", err)
		fmt.Println("Arg0:", uri)
		fmt.Println("Arg1:", searchWord)
		return err
	}
	defer fp.Close()

	// Search CSV and output hit lines to stdout
	reader := csv.NewReader(fp)
	reader.FieldsPerRecord = -1 // Accepts reading irregular csv format
	for {
		words, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("File read error:", err)
			return err
		}
		for _, w := range words {
			if strings.Contains(w, searchWord) {
				fmt.Println(words)
				break
			}
		}
	}
	return nil
}
