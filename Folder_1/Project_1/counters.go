package main

import (
	"sync"
)

type OperationalCounters struct {
	sync.RWMutex
	SQSMessageBatch             int `json:"sqs_message_batch"`
	LogChunkStatusUpdateError   int `json:"log_chunk_update_failed"`
	logChunkStatusUpdateSuccess int `json:"log_chunk_update_success"`
	OutOfUserTime               int `json:"out_of_user_time"`
	OutOfPrimeTime              int `json:"out_of_prime_time"`
	ScheduledFailedPlayerIds    int `json:"schedules_failed_player_ids"`
	ScheduledSuccessPlayerIds   int `json:"schedules_success_player_ids"`
	TotalPlayer                 int `json:"total_player"`
	InvalidPlayerIds            int `json:"invalid_player_ids"`
}

func (c *OperationalCounters) incSQSMessageBatch(num int) {
	c.Lock()
	defer c.Unlock()
	c.SQSMessageBatch = c.SQSMessageBatch + num
}

func (c *OperationalCounters) incLogChunkStatusUpdateError(num int) {
	c.Lock()
	defer c.Unlock()
	c.LogChunkStatusUpdateError = c.LogChunkStatusUpdateError + num
}

func (c *OperationalCounters) incLogChunkStatusUpdateSuccess(num int) {
	c.Lock()
	defer c.Unlock()
	c.logChunkStatusUpdateSuccess = c.logChunkStatusUpdateSuccess + num
}

func (c *OperationalCounters) incOutOfUserTime(num int) {
	c.Lock()
	defer c.Unlock()
	c.OutOfUserTime = c.OutOfUserTime + num
}

func (c *OperationalCounters) incOutOfPrimeTime(num int) {
	c.Lock()
	defer c.Unlock()
	c.OutOfPrimeTime = c.OutOfPrimeTime + num
}

func (c *OperationalCounters) incScheduledFailedPlayerIds(num int) {
	c.Lock()
	defer c.Unlock()
	c.ScheduledFailedPlayerIds = c.ScheduledFailedPlayerIds + num
}

func (c *OperationalCounters) incScheduledSuccessPlayerIds(num int) {
	c.Lock()
	defer c.Unlock()
	c.ScheduledSuccessPlayerIds = c.ScheduledSuccessPlayerIds + num
}

func (c *OperationalCounters) incTotalPlayer(num int) {
	c.Lock()
	defer c.Unlock()
	c.TotalPlayer = c.TotalPlayer + num
}

func (c *OperationalCounters) incInvalidPlayerIds(num int) {
	c.Lock()
	defer c.Unlock()
	c.InvalidPlayerIds = c.InvalidPlayerIds + num
}
