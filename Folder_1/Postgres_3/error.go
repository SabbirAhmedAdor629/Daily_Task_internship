package main

import "fmt"

type S3EventError struct {
	Type string
	Err  error
}

func (m *S3EventError) Error() string { return fmt.Sprintf("%s: err: %v", m.Type, m.Err) }

type MetadataError struct {
	Type string
	Err  error
}

func (m *MetadataError) Error() string { return fmt.Sprintf("%s: err %v", m.Type, m.Err) }

type TransactionIdError struct {
	Type string
	Err  error
}

func (m *TransactionIdError) Error() string { return fmt.Sprintf("%s: err %v", m.Type, m.Err) }

type AppIdError struct {
	Type string
	Err  error
}

func (m *AppIdError) Error() string { return fmt.Sprintf("%s: err %v", m.Type, m.Err) }

type PutObjectError struct {
	Type string
	Err  error
}

func (m *PutObjectError) Error() string { return fmt.Sprintf("%s: err %v", m.Type, m.Err) }

type DatabaseError struct {
	Type string
	Err  error
}

func (m *DatabaseError) Error() string { return fmt.Sprintf("%s: err %v", m.Type, m.Err) }

type FileError struct {
	Type string
	Err  error
}

func (m *FileError) Error() string { return fmt.Sprintf("%s: err %v", m.Type, m.Err) }

//
// SQSEventError
//
type SQSEventError struct {
	Type string
	Err  error
}

func (m *SQSEventError) Error() string { return fmt.Sprintf("%s: err: %v", m.Type, m.Err) }
