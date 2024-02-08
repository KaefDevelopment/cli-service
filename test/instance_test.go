package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(integrationSuite))
}

type integrationSuite struct {
	suite.Suite

	binaryPath string
}

func (s *integrationSuite) SetupSuite() {
	s.binaryPath = filepath.Join(os.TempDir(), "cli")

	s.Require().NoError(os.Setenv("CGO_ENABLED", "1"))

	cmd := exec.Command("go", []string{
		"build",
		"-cover",
		"-race",
		"-o", s.binaryPath,
		"../cmd/cli/main.go",
	}...)

	s.Require().NoError(cmd.Run())
}

func (s *integrationSuite) TestInstances() {
	var (
		wg   sync.WaitGroup
		args = []string{
			"event",
			"-d",
			"{\"events\":[{\"pluginId\":\"346d7f75-4b20-4166-8577-e656cdf3caec\",\"id\":\"\",\"createdAt\":\"3\",\"type\":\"2\",\"project\":\"2\",\"projectBaseDir\":\"/mnt/c/Users/jaros/GolandProjects/tts\",\"language\":\"golang\",\"target\":\"2\",\"branch\":\"\",\"timezone\":\"2\",\"params\":{\"count\":\"12\"}}]}",
			"-k",
			"346d7f75-4b20-4166-8577-e656cdf3caec",
			"-s",
			"http://localhost:8181/events",
		}
	)

	for i := 0; i < 10; i++ {
		cmd := exec.Command(s.binaryPath, args...)
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			s.Require().NoError(cmd.Run())
		}(i)
	}

	wg.Wait()
}
