package main

import (
"database/sql"
"fmt"
"reflect"
"strings"
"testing"
)

type Department struct {
	ID int
	Name string
}

func TestDepartment_Insert(t *testing.T) {
	tests := []struct {
			name string
			department Department
			expectedErr error
		}{
		{
			name: "valid department insertion",
			department: Department{
				ID: 1,
				Name: "Marketing",
			},
			expectedErr: nil,
		},
		{
			name: "invalid department insertion with missing ID",
			department: Department{
				Name: "HR",
			},
			expectedErr: fmt.Errorf("missing required field: ID"),
		},
		{
			name: "invalid department insertion with missing Name",
			department: Department{
				ID: 1,
			},
			expectedErr: fmt.Errorf("missing required field: Name"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, _ := sql.Open("test", "")
			defer db.Close()
			err := test.department.Insert(db)
			if err != test.expectedErr {
				t.Errorf("expected error %v, but got %v", test.expectedErr, err)
			}
		})
	}
}