package settings

type (
	IConfiguration interface {
		IsDockerized() bool
		IsTestEnvironment() bool
		IsDevelopmentEnvironment() bool
		IsStagingEnvironment() bool
		IsProductionEnvironment() bool
		IsTrafficRecordEnabled() bool
		IsRequestLogEnabled() bool
		GetAllowedOrigins() []string
		GetEnvironment() string
		GetServerConfiguration() IServerConfiguration
		GetInfluxConfiguration() IInfluxConfiguration
		GetPostgreSQLConfiguration() IPostgreSQLConfiguration
		GetMastodonApplication(string) IMastodonApplication
		GetPorts() (int, int, int)
	}

	IServerConfiguration interface {
		GetFQDN() string
		GetProtocol() string
		GetPortConfiguration() IPortConfiguration
		GetTLSConfiguration() ITLSConfiguration
		GetBuildNumber() int32
		SetBuildNumber(int32)
		GetJwtTokenKey() string
		GetJwtTokenExpiration() string
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

	IPostgreSQLConfiguration interface {
		GetHost() string
		SetHost(string)
		GetPort() string
		SetPort(string)
		GetDatabase() string
		SetDatabase(string)
		GetUsername() string
		SetUsername(string)
		GetPassword() string
		SetPassword(string)
	}

	IMastodonApplication interface {
		GetName() string
		GetServer() string
		GetClientID() string
		GetClientSecret() string
		GetToken() string
	}

	IInfluxConfiguration interface {
		IsEnabled() bool
		GetAddress() string
		GetDatabase() string
		GetUsername() string
		GetPassword() string
		GetReplicas() []string
	}
)
