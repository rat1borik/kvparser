package config

import (
	"kvparser/internal/utils"
	"log"
	"os"
	"sync"
)

type CookieLoader interface {
	Load() (string, error)
	Save(string) error
}

type cookieLoader struct {
	cookie string
	mu     sync.RWMutex
}

func NewCookieLoader() CookieLoader {
	return &cookieLoader{}
}

func (cl *cookieLoader) Load() (string, error) {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	if cl.cookie != "" {
		return cl.cookie, nil
	}

	path, _ := utils.ExecPath("cookie")
	val, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	cl.cookie = string(val)
	return string(val), nil
}

func (cl *cookieLoader) Save(val string) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	log.Println(val)

	if cl.cookie == val {
		return nil
	}

	path, _ := utils.ExecPath("cookie")

	err := os.WriteFile(path, []byte(val), 0644)
	if err != nil {
		return err
	}

	cl.cookie = val
	return nil
}
