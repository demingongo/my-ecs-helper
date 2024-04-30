package aws

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

type TargetGroup struct {
	Arn  string
	Name string
}

func DescribeTargetGroups() ([]TargetGroup, error) {
	options := []TargetGroup{}
	cmd := "aws elbv2 describe-target-groups --output json --no-paginate"
	log.Debug(cmd)
	if viper.GetBool("dummy") {
		for i := 1; i <= 10; i += 1 {
			arn := "arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/ecs-tg-" + strconv.Itoa(i) + "/73e2d6bc24d8a067"
			name := "ecs-project-" + strconv.Itoa(i)
			options = append(options, TargetGroup{Arn: arn, Name: name})
		}
		return options, nil
	}

	// @TODO

	return options, nil
}

func CreateTargetGroup(filepath string) (TargetGroup, error) {
	var r TargetGroup
	cmd := fmt.Sprintf("aws elbv2 create-target-group --output json --cli-input-json \"$(cat %s)\"", filepath)
	log.Debug(cmd)
	if viper.GetBool("dummy") {
		return TargetGroup{
			Arn:  "arn:dummy",
			Name: "dummy",
		}, nil
	}

	// @TODO

	return r, nil
}
