package aws

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

type TargetGroup struct {
	TargetGroupArn  string
	TargetGroupName string
}

type targetGroupsResponse struct {
	TargetGroups []TargetGroup
}

func DescribeTargetGroups() ([]TargetGroup, error) {
	options := []TargetGroup{}
	var args []string
	args = append(args, "elbv2", "describe-target-groups", "--output", "json", "--no-paginate")
	log.Debug(args)
	if viper.GetBool("dummy") {
		for i := 1; i <= 10; i += 1 {
			arn := "arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/ecs-tg-" + strconv.Itoa(i) + "/73e2d6bc24d8a067"
			name := "ecs-project-" + strconv.Itoa(i)
			options = append(options, TargetGroup{TargetGroupArn: arn, TargetGroupName: name})
		}
		return options, nil
	}

	var resp targetGroupsResponse
	_, err := execAWS(args, &resp)
	if err != nil {
		return options, err
	}

	options = resp.TargetGroups

	return options, nil
}

func CreateTargetGroup(filepath string) (TargetGroup, error) {
	var r TargetGroup
	cmd := fmt.Sprintf("aws elbv2 create-target-group --output json --cli-input-json \"$(cat %s)\"", filepath)
	log.Debug(cmd)
	if viper.GetBool("dummy") {
		return TargetGroup{
			TargetGroupArn:  "arn:dummy",
			TargetGroupName: "dummy",
		}, nil
	}

	// @TODO

	return r, nil
}
