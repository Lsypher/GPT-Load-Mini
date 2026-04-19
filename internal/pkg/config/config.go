package config

import (
	"os"
	"sync"
)

type Config struct {
	mu            sync.RWMutex
	ServerPort    string
	AuthKey       string
	MySQLHost     string
	MySQLPort     string
	MySQLUser     string
	MySQLPassword string
	MySQLDatabase string
	RedisHost     string
	RedisPort     string
	EncryptionKey string
}

func Load() *Config {
	return &Config{
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		AuthKey:       getEnv("AUTH_KEY", ""),
		MySQLHost:     getEnv("MYSQL_HOST", "localhost"),
		MySQLPort:     getEnv("MYSQL_PORT", "3306"),
		MySQLUser:     getEnv("MYSQL_USER", "root"),
		MySQLPassword: getEnv("MYSQL_PASSWORD", ""),
		MySQLDatabase: getEnv("MYSQL_DATABASE", "gptload"),
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		EncryptionKey: getEnv("ENCRYPTION_KEY", ""),
	}
}

func (c *Config) Reload() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ServerPort = getEnv("SERVER_PORT", c.ServerPort)
	c.AuthKey = getEnv("AUTH_KEY", c.AuthKey)
	c.MySQLHost = getEnv("MYSQL_HOST", c.MySQLHost)
	c.MySQLPort = getEnv("MYSQL_PORT", c.MySQLPort)
	c.MySQLUser = getEnv("MYSQL_USER", c.MySQLUser)
	c.MySQLPassword = getEnv("MYSQL_PASSWORD", c.MySQLPassword)
	c.MySQLDatabase = getEnv("MYSQL_DATABASE", c.MySQLDatabase)
	c.RedisHost = getEnv("REDIS_HOST", c.RedisHost)
	c.RedisPort = getEnv("REDIS_PORT", c.RedisPort)
	c.EncryptionKey = getEnv("ENCRYPTION_KEY", c.EncryptionKey)
}

func (c *Config) GetServerPort() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ServerPort
}

func (c *Config) GetAuthKey() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.AuthKey
}

func (c *Config) GetMySQLHost() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.MySQLHost
}

func (c *Config) GetMySQLPort() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.MySQLPort
}

func (c *Config) GetMySQLUser() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.MySQLUser
}

func (c *Config) GetMySQLPassword() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.MySQLPassword
}

func (c *Config) GetMySQLDatabase() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.MySQLDatabase
}

func (c *Config) GetRedisHost() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.RedisHost
}

func (c *Config) GetRedisPort() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.RedisPort
}

func (c *Config) GetEncryptionKey() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.EncryptionKey
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
