package aws

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

type ServiceLoadBalancer struct {
	TargetGroupArn string
	ContainerName  string
	ContainerPort  int
}

func CreateService(filepath string, loadBalancer ServiceLoadBalancer) (string, error) {
	cmd := fmt.Sprintf("aws ecs create-service --output json --cli-input-json \"$(cat %s)\"", filepath)
	if loadBalancer.TargetGroupArn != "" && loadBalancer.ContainerName != "" {
		action := fmt.Sprintf(
			"--load-balancers targetGroupArn=%s,containerName=%s,containerPort=%d",
			loadBalancer.TargetGroupArn, loadBalancer.ContainerName, loadBalancer.ContainerPort,
		)
		cmd += " " + action
	}
	log.Debug(cmd)
	if viper.GetBool("dummy") {
		return cmd, nil
	}

	// @TODO

	return cmd, nil
}
