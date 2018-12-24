package util

import "os"

const (
	// BootstrapServers : bootstrap servers list
	BootstrapServers string = "BOOTSTRAP_SERVERS"
	// Topic : topic
	Topic string = "TOPIC"
	// GroupID : consumer group
	GroupID string = "GROUP_ID"
)

// GetEnv : returns the environment variable value
func GetEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}
