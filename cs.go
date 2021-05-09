package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	flag.Parse()
	searchCsv(flag.Arg(0), flag.Arg(1))
}

func searchCsv(fileName, searchWord string) error {
	// Handle input file
	fp, err := os.Open(fileName)
	if err != nil {
		fmt.Println("File open error:", err)
		fmt.Println("Arg0:", fileName)
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
