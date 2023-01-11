package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
)

type MySelf struct {
	Name         string
	Parents_Info []parentsInfo
}

type parentsInfo struct {
	Type       string
	Name       string
	Occupation string
	Fav_colour []string
}

func convertJSONToCSV(source, destination string) error {
	// 2. Read the JSON file into the struct
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	// remember to close the file at the end of the function
	defer sourceFile.Close()

	// variable
	var self2 MySelf
	if err := json.NewDecoder(sourceFile).Decode(&self2); err != nil {
		return err
	}

	// 3. Create a new file to store CSV data
	outputFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// 4. Write the header of the CSV file and the successive rows by iterating through the JSON struct array
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	header := []string{"Name", "Parents_Info"}
	if err := writer.Write(header); err != nil {
		return err
	}

	var csvRow []string

	for _, r := range self2.Parents_Info {

		csvRow = append(csvRow, r.Type, r.Name, r.Occupation, r.Fav_colour[0], r.Fav_colour[1], r.Fav_colour[2])
		if err := writer.Write(csvRow); err != nil {
			return err
		}
	}
	return nil

}

func main() {
	if err := convertJSONToCSV("test.json", "csvdata1.csv"); err != nil {
		log.Fatal(err)
	}
}
