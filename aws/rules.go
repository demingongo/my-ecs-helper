package aws

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

func CreateRule(filepath string, targetGroupArn string) (string, error) {
	var args []string
	args = append(args, "elbv2", "create-rule", "--cli-input-json", fmt.Sprintf("$(cat %s)", filepath))
	if targetGroupArn != "" {
		args = append(args, "--action", fmt.Sprintf(`[{
			\"Type\": \"forward\",
			\"TargetGroupArn\": \"%s\",
			\"Order\": 1,
			\"ForwardConfig\": {
				\"TargetGroups\": [
					{
						\"TargetGroupArn\": \"%s\",
						\"Weight\": 1
					}
				],
				\"TargetGroupStickinessConfig\": {
					\"Enabled\": false,
					\"DurationSeconds\": 3600
				}
			}
		}]`, targetGroupArn, targetGroupArn))
	}
	log.Debug(args)
	if viper.GetBool("dummy") {
		return strings.Join(args, " "), nil
	}

	var resp any
	stdout, err := execAWS(args, &resp)

	return string(stdout), err
}
