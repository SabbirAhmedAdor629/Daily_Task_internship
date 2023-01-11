package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
)

type FamilyDetails struct {
	Myself       string
	EMAIL        string
	F_name       string
	F_occupation string
	F_MOBILE     string
	M_name       string
	M_occupation string
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
	var fm3 FamilyDetails
	if err := json.NewDecoder(sourceFile).Decode(&fm3); err != nil {
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

	header := []string{"Myself", "Email", "F_name", "F_occupation", "F_mobile", "M_name", "M_occupation"}
	if err := writer.Write(header); err != nil {
		return err
	}

	var csvRow []string
	csvRow = append(csvRow, fm3.Myself, fm3.EMAIL, fm3.F_name, fm3.F_occupation, fm3.F_MOBILE, fm3.M_name, fm3.M_occupation)
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
