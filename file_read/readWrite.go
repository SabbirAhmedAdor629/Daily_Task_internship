package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	// "strconv"
)

type Mother struct {
	M_name       string
	M_occupation string
}

type Father struct {
	F_name       string
	F_occupation string
	F_MOBILE     string
}

type FamilyDetails struct {
	Myself string
	EMAIL  string
	Father
	Mother
}

func main() {
	// Details 1
	family_details_1 := FamilyDetails{}
	family_details_1.Myself = "Sabbir"

	family_details_1.EMAIL = "fadd@gmail.com"
	// FATHER DETAILS
	family_details_1.F_name = "Ahmed"
	family_details_1.F_occupation = "Businessman"
	family_details_1.F_MOBILE = "01788988568"
	// MOTHER DETAILS
	family_details_1.M_name = "Jahanara Begum"
	family_details_1.M_occupation = "Teacher"

	// Details 4
	family_details_4 := FamilyDetails{}
	family_details_4.Myself = "Ador"
	family_details_4.EMAIL = "fadd@gmail.com"
	// FATHER DETAILS
	family_details_4.F_name = "Ahmed"
	family_details_4.F_occupation = "Businessman"
	family_details_4.F_MOBILE = "01788988568"
	// MOTHER DETAILS
	family_details_4.M_name = "Jahanara Begum"
	family_details_4.M_occupation = "Teacher"

	// Writing Data in JSON
	file, _ := json.MarshalIndent(family_details_4, "", " ")
	_ = ioutil.WriteFile("test.json", file, 0777)

	// Reading the JSON file
	jsonFile, _ := os.Open("test.json") // Openning our json file

	defer jsonFile.Close() // Closing of our json file to parse it later on

	byteValue, _ := ioutil.ReadAll(jsonFile) // read our opened xml file as bytevalue

	family_details_2 := FamilyDetails{} // another object to storing and printing out value

	_ = json.Unmarshal(byteValue, &family_details_2) // Unmarshalling

	fmt.Println(family_details_2)
}
