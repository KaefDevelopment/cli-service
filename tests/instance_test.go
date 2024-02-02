package tests

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
)

var wg sync.WaitGroup

func TestInstances(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Error(err)
	}

	ch := make(chan int, 10)

	path := filepath.Join(homeDir, ".nau", "main") // main - бинарник (надо сбилдить отдельно от теста)
	t.Log(path)
	args := []string{
		"event",
		"-d",
		"{\"events\":[{\"pluginId\":\"346d7f75-4b20-4166-8577-e656cdf3caec\",\"id\":\"\",\"createdAt\":\"3\",\"type\":\"2\",\"project\":\"2\",\"projectBaseDir\":\"/mnt/c/Users/jaros/GolandProjects/tts\",\"language\":\"golang\",\"target\":\"2\",\"branch\":\"\",\"timezone\":\"2\",\"params\":{\"count\":\"12\"}}]}",
		"-k",
		"346d7f75-4b20-4166-8577-e656cdf3caec",
		"-s",
		"http://localhost:8181/events",
	}

	cmd := exec.Command(path, args...)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ch <- i
			if err := cmd.Run(); err != nil {
				t.Error(err)
			}
			fmt.Printf("i:%d - success done", i)
		}(i)
	}

	wg.Wait()
}
