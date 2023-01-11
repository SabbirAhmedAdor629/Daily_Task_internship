package main

import "fmt"

// SQSEventError

type SQSEventError struct {
	Type string
	Err  error
}

func (m *SQSEventError) Error() string { return fmt.Sprintf("%s: err: %v", m.Type, m.Err) }

type DatabaseError struct {
	Type string
	Err  error
}

func (m *DatabaseError) Error() string { return fmt.Sprintf("%s: err %v", m.Type, m.Err) }
