package aws

import (
	"encoding/json"
	"os/exec"
)

func execAWS[T any](args []string, resp *T) ([]byte, error) {
	cmd := exec.Command("aws", args...)
	stdout, err := cmd.Output()
	if err != nil {
		return stdout, err
	}
	err = json.Unmarshal(stdout, resp)
	return stdout, err
}
