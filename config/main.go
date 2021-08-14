package config

type (
	Configurator interface {
		Validate() error
	}

	ServiceConfigurator interface {
		LoadEnvVariables(commit, builtAt string)
		IsLoaded() bool
		ConfigureLogger() error
	}
)
