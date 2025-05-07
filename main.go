package main

import (
	"os"
	"time"
)

type AccessControl struct {
	Start time.Time
	End   time.Time
	Mode  os.FileMode
}

type AccessManager struct {
	Rules map[string]AccessControl
}

func NewAccessManager() *AccessManager {
	return &AccessManager{
		Rules: make(map[string]AccessControl),
	}
}

func (am *AccessManager) Add(file string, start, end time.Time, mode os.FileMode) {
	am.Rules[file] = AccessControl{
		Start: start,
		End:   end,
		Mode:  mode,
	}
}

func (am *AccessManager) Enforce() {
	now := time.Now()
	for file, rule := range am.Rules {
		if now.After(rule.Start) && now.Before(rule.End) {
			os.Chmod(file, rule.Mode)
		} else {
			os.Chmod(file, 0000) //no permission
		}
	}
}

func main() {

	testFile := "test.txt"

	os.WriteFile(testFile, []byte("secret"), 0644)

	manager := NewAccessManager()

	now := time.Now()
	manager.Add(testFile, now, now.Add(10*time.Second), 0600)

	for range 10 {
		manager.Enforce()
		time.Sleep(2 * time.Second)
	}

}
