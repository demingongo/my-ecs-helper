package aws

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

type ContainerPortMapping struct {
	Name        string
	PortMapping PortMapping
}

type PortMapping struct {
	ContainerPort int
	HostPort      int
	Name          string
}

type ContainerDefinition struct {
	Name         string
	Image        string
	PortMappings []PortMapping
}

type TaskDefinition struct {
	Arn                  string
	ContainerDefinitions []ContainerDefinition
}

func DescribeTaskDefinition(taskDefinitionArn string) ([]TaskDefinition, error) {
	options := []TaskDefinition{}
	cmd := "aws ecs describe-task-definition --output json --no-paginate --task-definition " + taskDefinitionArn
	log.Debug(cmd)
	if viper.GetBool("dummy") {
		options = append(options, TaskDefinition{
			Arn: "arn:task-definition/taskdef-ci-dmz-web:5",
			ContainerDefinitions: []ContainerDefinition{
				{
					Name:  "dmz-web",
					Image: "xxx/repository-dmz-web:tag",
					PortMappings: []PortMapping{
						{
							ContainerPort: 8080,
							HostPort:      0,
							Name:          "http",
						},
					},
				},
			},
		})
		return options, nil
	}

	// @TODO

	return options, nil
}

func ListPortMapping(taskDefinitionArn string) ([]ContainerPortMapping, error) {
	options := []ContainerPortMapping{}
	v, err := DescribeTaskDefinition(taskDefinitionArn)
	if err != nil {
		return options, err
	}
	for _, td := range v {
		for _, cd := range td.ContainerDefinitions {
			for _, pm := range cd.PortMappings {
				options = append(options, ContainerPortMapping{
					Name:        cd.Name,
					PortMapping: pm,
				})
			}
		}
	}

	// @TODO

	return options, nil
}
