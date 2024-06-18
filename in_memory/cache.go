package in_memory

import "time"

//Interface definition for API method implementation
type Cache interface {
	Set(key string, value interface{}, expiration time.Duration)
	Get(key string) (interface{}, bool)
	GetAll() map[string]interface{}
	Delete(key string) bool
	DeleteAll() bool
}
