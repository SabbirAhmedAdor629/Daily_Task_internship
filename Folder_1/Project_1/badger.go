package main

import "github.com/honeybadger-io/honeybadger-go"

type (
	badgerClient     struct{}
	mockBadgerClient struct{}
)

type badgerIface interface {
	Configure(c honeybadger.Configuration)
	Monitor()
	Notify(err interface{}, extra ...interface{}) (string, error)
	Flush()
}

func (bc *badgerClient) Flush()                                { honeybadger.Flush() }
func (bc *badgerClient) Monitor()                              { honeybadger.Monitor() }
func (bc *badgerClient) Configure(c honeybadger.Configuration) { honeybadger.Configure(c) }
func (bc *badgerClient) Notify(err interface{}, extra ...interface{}) (string, error) {
	return honeybadger.Notify(err, extra...)
}

func (bc *mockBadgerClient) Flush()                                {}
func (bc *mockBadgerClient) Monitor()                              {}
func (bc *mockBadgerClient) Configure(c honeybadger.Configuration) {}
func (bc *mockBadgerClient) Notify(err interface{}, extra ...interface{}) (string, error) {
	return "", nil
}
