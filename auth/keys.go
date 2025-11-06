package auth

import "github.com/franzego/distributed_task_queue/authutil"

// KeyGenerator delegates to authutil.KeyGenerator to keep a stable auth API
func KeyGenerator(prefix string) (string, error) {
	return authutil.KeyGenerator(prefix)
}

// HashApiKeys delegates to authutil.HashApiKeys
func HashApiKeys(key string) string {
	return authutil.HashApiKeys(key)
}
