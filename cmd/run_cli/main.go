package main

import (
	"fmt"
	"os/exec"
)

func main() {
	args := []string{"C:/Users/jaros/Downloads/cli-windows-amd64.exe/cli-windows-amd64.exe", "-d", "{\"events\":[{\"id\":\"\",\"createdAt\":\"2\",\"type\":\"2\",\"project\":\"2\",\"projectBaseDir\":\"C:/Users/jaros/GolandProjects/cli\",\"language\":\"golang\",\"target\":\"2\",\"branch\":\"\",\"timezone\":\"2\",\"params\":{\"count\":\"12\"}}]}", "-k", "qwerty12345"}
	output, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(output))
}
