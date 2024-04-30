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
	TaskDefinitionArn    string
	ContainerDefinitions []ContainerDefinition
}

type describeTaskDefinitionResponse struct {
	taskDefinition TaskDefinition
}

func DescribeTaskDefinition(taskDefinitionArn string) (TaskDefinition, error) {
	result := TaskDefinition{}
	var args []string
	args = append(args, "ecs", "describe-task-definition", "--output", "json", "--no-paginate", "--task-definition", taskDefinitionArn)
	log.Debug(args)
	if viper.GetBool("dummy") {
		result = TaskDefinition{
			TaskDefinitionArn: taskDefinitionArn,
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
		}
		return result, nil
	}

	var resp describeTaskDefinitionResponse
	_, err := execAWS(args, &resp)
	if err != nil {
		return result, err
	}

	result = resp.taskDefinition

	return result, nil
}

func ListPortMapping(taskDefinitionArn string) ([]ContainerPortMapping, error) {
	options := []ContainerPortMapping{}
	td, err := DescribeTaskDefinition(taskDefinitionArn)
	if err != nil {
		return options, err
	}
	for _, cd := range td.ContainerDefinitions {
		for _, pm := range cd.PortMappings {
			options = append(options, ContainerPortMapping{
				Name:        cd.Name,
				PortMapping: pm,
			})
		}
	}

	// @TODO

	return options, nil
}
