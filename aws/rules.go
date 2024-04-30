package aws

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

func CreateRule(filepath string, targetGroupArn string) (string, error) {
	cmd := fmt.Sprintf("aws elbv2 create-rule --cli-input-json \"$(cat %s)\"", filepath)
	if targetGroupArn != "" {
		action := fmt.Sprintf("--actions Type=forward,TargetGroupArn=%s", targetGroupArn)
		cmd += " " + action
	}
	log.Debug(cmd)
	if viper.GetBool("dummy") {
		return cmd, nil
	}

	// @TODO

	return cmd, nil
}
