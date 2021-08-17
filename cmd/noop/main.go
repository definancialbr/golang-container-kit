package main

import (
	"time"

	"github.com/definancialbr/golang-container-kit/pkg/configuration/viper"
	"github.com/definancialbr/golang-container-kit/pkg/container"
	"github.com/definancialbr/golang-container-kit/pkg/logging/zap"
	"github.com/definancialbr/golang-container-kit/pkg/metrics/prometheus"
	"github.com/definancialbr/golang-container-kit/pkg/probes/healthcheck"
	"github.com/definancialbr/golang-container-kit/pkg/signaler"
)

func terminationHandler(release func()) error {
	release()
	return nil
}

func main() {

	cont := container.NewContainer()

	cont.Configuration = viper.NewConfigurationService(
		viper.WithOptionalConfigurationFile(),
		viper.WithEnvPrefix("NOOP"),
		viper.WithFileName("noop.env"),
		viper.WithFileType("env"),
		viper.WithSearchPaths("."),
		viper.WithHomeSearchPath(),
		viper.WithConfiguration("somevar", "hello"),
	)

	cont.Logging = zap.NewLoggingService(
		zap.WithDevelopmentMode(),
		zap.WithName("noop"),
	)

	cont.Metrics = prometheus.NewMetricService()

	cont.Probes = healthcheck.NewProbeService(
		healthcheck.WithDNSResolveCheckForLiveness("google-is-resolvable", "google.com", 10*time.Second),
		healthcheck.WithTCPDialCheckForReadiness("google-is-reachable", "google.com", 10*time.Second),
	)

	hangupHandler := func(func()) error {
		cont.Logging.Debug("O grande problema que a nação está enfrentando hoje é a falta de amor!")
		return nil
	}

	terminationHandler := func(release func()) error {
		cont.Logging.Info("Gastamos R$ 700 e alguma coisa na campanha presidencial.")
		release()
		return nil
	}

	interruptHandler := func(release func()) error {
		cont.Logging.Warn("A democracia é uma delícia!")
		release()
		return nil
	}

	cont.Signaler = signaler.NewSignaler(
		signaler.WithOnHangup(hangupHandler),
		signaler.WithOnTermination(terminationHandler),
		signaler.WithOnInterrupt(interruptHandler),
	)

	cont.Open()

	config := cont.Configuration.Load()

	cont.Logging.Info("Glória a Deux!", "somevar", config.GetString("somevar"))

	cont.Signaler.WaitForSignal(func(err error) {
		cont.Logging.Error(err.Error())
	})

	cont.Close()

}
