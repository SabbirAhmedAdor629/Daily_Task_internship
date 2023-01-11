package main

import (
	"influencemobile.com/logging"
)

func (d *mockPgClient) StatusUpdate(recordId int, status string) error {
	logger.
		NewEntry(log_msg_update_delete_status, LogMsgTxt).
		WithFields(logging.Fields{"recordId": recordId, "status": status}).
		Debug()
	return nil
}

func (d *mockPgClient) ZendeskUserIdUpdate(recordId int, zendeskUserId int, status string) error {
	logger.
		NewEntry(log_msg_update_delete_status, LogMsgTxt).
		WithFields(logging.Fields{"recordId": recordId, "zendeskUserId": zendeskUserId, "status": status}).
		Debug()
	return nil
}

func (d *mockPgClient) Ping() error { return nil }
