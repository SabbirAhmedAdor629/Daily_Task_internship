package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Mother struct {
	M_name       string
	M_fav_colour []string
}

type Father struct {
	F_name       string
	F_fav_colour []string
}

type Details struct {
	Name  string
	Id    int32
	EMAIL string
	Father
	Mother
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
	var details2 Details
	if err := json.NewDecoder(sourceFile).Decode(&details2); err != nil {
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

	header := []string{"Name","ID", "Email", "F_name", "F_fav_colour", "M_name", "M_fav_colour"}
	if err := writer.Write(header); err != nil {
		return err
	}

	var csvRow []string
	csvRow = append(csvRow, details2.Name, details2.EMAIL, details2.F_name, details2.F_fav_colour[], details2.M_name, details2.M_fav_colour[])
	if err := writer.Write(csvRow); err != nil {
		return err
	}

	return nil

}

func main() {
	if err := convertJSONToCSV("test.json", "csvdata1.csv"); err != nil {
		log.Fatal(err)
	}
}