package im_badger

import "github.com/honeybadger-io/honeybadger-go"

type (
	BadgerClient     struct{ BadgerIface }
	MockBadgerClient struct{ BadgerIface }
)

type BadgerIface interface {
	Configure(c honeybadger.Configuration)
	Monitor()
	Notify(err interface{}, extra ...interface{}) (string, error)
	Flush()
}

func (bc *BadgerClient) Flush()                                { honeybadger.Flush() }
func (bc *BadgerClient) Monitor()                              { honeybadger.Monitor() }
func (bc *BadgerClient) Configure(c honeybadger.Configuration) { honeybadger.Configure(c) }
func (bc *BadgerClient) Notify(err interface{}, extra ...interface{}) (string, error) {
	return honeybadger.Notify(err, extra...)
}

func (bc *MockBadgerClient) Flush()                                {}
func (bc *MockBadgerClient) Monitor()                              {}
func (bc *MockBadgerClient) Configure(c honeybadger.Configuration) {}
func (bc *MockBadgerClient) Notify(err interface{}, extra ...interface{}) (string, error) {
	return "", nil
}
