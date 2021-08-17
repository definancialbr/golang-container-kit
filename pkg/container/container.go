package container

import (
	"github.com/definancialbr/golang-container-kit/pkg/configuration"
	"github.com/definancialbr/golang-container-kit/pkg/logging"
	"github.com/definancialbr/golang-container-kit/pkg/metrics"
	"github.com/definancialbr/golang-container-kit/pkg/probes"
	"github.com/definancialbr/golang-container-kit/pkg/signaler"
)

type ContainerServiceState int

const (
	Closed ContainerServiceState = iota
	Open
)

type Container struct {
	Configuration configuration.ConfigurationService
	Logging       logging.LoggingService
	Metrics       metrics.MetricService
	Signaler      signaler.SignalerService
	Probes        probes.ProbeService

	configurationState ContainerServiceState
	loggingState       ContainerServiceState
}

func NewContainer() *Container {
	return &Container{
		configurationState: Closed,
		loggingState:       Closed,
	}
}

func (c *Container) Open() {

	if c.configurationState == Closed && c.Configuration != nil {

		if err := c.Configuration.Read(); err != nil {
			panic(err)
		}

		c.configurationState = Open

	}

	if c.loggingState == Closed && c.Logging != nil {

		if err := c.Logging.Open(); err != nil {
			panic(err)
		}

		c.loggingState = Open

	}

}

func (c *Container) Close() {

	if c.loggingState == Open {

		if err := c.Logging.Close(); err != nil {
			panic(err)
		}

		c.loggingState = Closed

	}

}
