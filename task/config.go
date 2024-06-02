package task

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

type Config struct {
	Name                     string
	AttachStdin              bool
	AttachStdout             bool
	AttachStderr             bool
	ExposedPorts             nat.PortSet
	Cmd                      []string
	Image                    string
	CPU                      float64
	Memory                   int64
	Disk                     int64
	Env                      []string
	RestartPolicy            container.RestartPolicyMode
	RestartMaximumRetryCount int
}
