package main

import (
"testing"
"fmt"
)

type testCase struct {
	description string
	input interface{}
	id int
	expected error
}

func TestInsert(t *testing.T) {
	db := &DB{}
	err := db.Connect()
	if err != nil {
	t.Fatal(err)
	}
	defer db.Close()
	testCases := []testCase{
		{
			description: "Insert new department",
			input: &Department{
				ID:       6075,
				Dept_Name: "SWE",
				DeptCode: "",
			},
			expected: nil,
		},
		{
			description: "Insert new student",
			input: &Student{
				Id:    2505464,
				Name:  "Ador",
				Email: "sdfasfaaf",
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := db.Insert(tc.input)
			if err != tc.expected {
				t.Errorf("Expected: %v, Got: %v", tc.expected, err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	db := &DB{}
	err := db.Connect()
	if err != nil {
	t.Fatal(err)
	}
	defer db.Close()
	testCases := []testCase{
		{
			description: "Update department name",
			input: &Department{
				ID:       6075,
				Dept_Name: "Computer Science",
				DeptCode: "CS",
			},
			id:       6075,
			expected: nil,
		},
		{
			description: "Update non-existent department",
			input: &Department{
				ID:       6076,
				Dept_Name: "Computer Science",
				DeptCode: "CS",
			},
			id:       6076,
			expected: fmt.Errorf("record not found"),
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := db.Update(tc.input, tc.id)
			if err != tc.expected {
				t.Errorf("Expected: %v, Got: %v", tc.expected, err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	db := &DB{}
	err := db.Connect()
	if err != nil {
	t.Fatal(err)
	}
	defer db.Close()
	testCases := []testCase{
		{
			description: "Delete existing department",
			input:       &Department{},
			id:          6075,
			expected:    nil,
		},
		{
			description: "Delete non-existent department",
			input:       &Department{},
			id:          6076,
			expected:    fmt.Errorf("record not found"),
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := db.Delete(tc.input, tc.id)
			if err != tc.expected {
				t.Errorf("Expected: %v, Got: %v", tc.expected, err)
			}
		})
	}
}	