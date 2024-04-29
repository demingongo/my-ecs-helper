package aws

import "strconv"

type TargetGroup struct {
	Arn  string
	Name string
}

func DescribeTargetGroups() []TargetGroup {
	options := []TargetGroup{}
	for i := 1; i <= 10; i += 1 {
		arn := "arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/ecs-tg-" + strconv.Itoa(i) + "/73e2d6bc24d8a067"
		name := "ecs-project-" + strconv.Itoa(i)
		options = append(options, TargetGroup{Arn: arn, Name: name})
	}
	return options
}
