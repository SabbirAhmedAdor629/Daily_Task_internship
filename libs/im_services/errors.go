package im_services

import "fmt"

type RedisError struct {
	Type string
	Err  error
}

func (m *RedisError) Error() string { return fmt.Sprintf("%s: err %v", m.Type, m.Err) }
