package settings

type (
	IConfiguration interface {
		IsDockerized() bool
		IsTestEnvironment() bool
		IsDevelopmentEnvironment() bool
		IsStagingEnvironment() bool
		IsProductionEnvironment() bool
		GetEnvironment() string
		GetServerConfiguration() IServerConfiguration
		GetInfluxConfiguration() IInfluxConfiguration
		GetMySQLConfiguration() IMySqlConfiguration
		GetClientsConfiguration() []IClientConfiguration
		GetPorts() (int, int, int)
	}

	IServerConfiguration interface {
		GetFQDN() string
		GetProtocol() string
		GetPortConfiguration() IPortConfiguration
		GetTLSConfiguration() ITLSConfiguration
		GetBuildNumber() int32
		SetBuildNumber(int32)
		GetHashKey() string
		GetBlockKey() string
	}

	IPortConfiguration interface {
		GetActive() int
		GetPassive() int
		GetDiagnostics() int
	}

	ITLSConfiguration interface {
		IsEnabled() bool
		GetKeyFile() string
		GetCertFile() string
	}

	IMySqlConfiguration interface {
		GetAddress() string
		GetDatabase() string
		GetUsername() string
		GetPassword() string
		IsPasswordSkipped() bool
	}

	IInfluxConfiguration interface {
		IsEnabled() bool
		GetAddress() string
		GetDatabase() string
		GetUsername() string
		GetPassword() string
		GetReplicas() []string
	}

	IClientConfiguration interface {
		GetId() string
		GetUrl() string
	}
)
