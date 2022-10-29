package settings

import (
	"fmt"
	"os"
	"strings"

	. "github.com/xeronith/diamante/contracts/settings"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

type Server struct {
	FQDN        string `yaml:"fqdn"`
	Protocol    string `yaml:"protocol"`
	Ports       *Ports `yaml:"ports"`
	TLS         *TLS   `yaml:"tls"`
	BuildNumber int32  `yaml:"build_number"`
	HashKey     string `yaml:"hash_key"`
	BlockKey    string `yaml:"block_key"`
}

func (server *Server) GetFQDN() string {
	if server.FQDN == "" {
		return "localhost"
	}

	return server.FQDN
}

func (server *Server) GetProtocol() string {
	switch server.Protocol {
	case "http", "https":
		return server.Protocol
	default:
		return "http"
	}
}

func (server *Server) GetPortConfiguration() IPortConfiguration {
	if server.Ports == nil {
		server.Ports = &Ports{
			Active:      0,
			Passive:     0,
			Diagnostics: 0,
		}
	}

	return server.Ports
}

func (server *Server) GetTLSConfiguration() ITLSConfiguration {
	if server.TLS == nil {
		server.TLS = &TLS{
			KeyFile:  "",
			CertFile: "",
		}
	}

	return server.TLS
}

func (server *Server) GetBuildNumber() int32 {
	return server.BuildNumber
}

func (server *Server) SetBuildNumber(value int32) {
	server.BuildNumber = value
}

func (server *Server) GetHashKey() string {
	return server.HashKey
}

func (server *Server) GetBlockKey() string {
	return server.BlockKey
}

//------------------------------------------------------------------------------------------------------------

type Ports struct {
	Active      int `yaml:"active"`
	Passive     int `yaml:"passive"`
	Diagnostics int `yaml:"diagnostics"`
}

func (ports *Ports) GetActive() int {
	return ports.Active
}

func (ports *Ports) GetPassive() int {
	return ports.Passive
}

func (ports *Ports) GetDiagnostics() int {
	return ports.Diagnostics
}

//------------------------------------------------------------------------------------------------------------

type TLS struct {
	KeyFile  string `yaml:"key_file"`
	CertFile string `yaml:"cert_file"`
}

func (tls *TLS) IsEnabled() bool {
	return tls.CertFile != "" && tls.KeyFile != ""
}

func (tls *TLS) GetKeyFile() string {
	return tls.KeyFile
}

func (tls *TLS) GetCertFile() string {
	return tls.CertFile
}

//------------------------------------------------------------------------------------------------------------

type MySQL struct {
	Address      string `yaml:"address"`
	Database     string `yaml:"database"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	SkipPassword bool   `yaml:"skip_password"`
}

func (mysql *MySQL) GetAddress() string {
	if mysql.Address == "" {
		mysql.Address = "localhost:3306"
	}

	return mysql.Address
}

func (mysql *MySQL) GetDatabase() string {
	return mysql.Database
}

func (mysql *MySQL) GetUsername() string {
	return mysql.Username
}

func (mysql *MySQL) GetPassword() string {
	return mysql.Password
}

func (mysql *MySQL) IsPasswordSkipped() bool {
	return mysql.SkipPassword
}

//------------------------------------------------------------------------------------------------------------

type Influx struct {
	Enabled  bool     `yaml:"enabled"`
	Address  string   `yaml:"address"`
	Database string   `yaml:"database"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Replicas []string `yaml:"replicas"`
}

func (influx *Influx) GetAddress() string {
	return influx.Address
}

func (influx *Influx) IsEnabled() bool {
	return influx.Enabled
}

func (influx *Influx) GetDatabase() string {
	return influx.Database
}

func (influx *Influx) GetUsername() string {
	return influx.Username
}

func (influx *Influx) GetPassword() string {
	return influx.Password
}

func (influx *Influx) GetReplicas() []string {
	return influx.Replicas
}

//------------------------------------------------------------------------------------------------------------

type Client struct {
	Id  string `yaml:"id"`
	Url string `yaml:"url"`
}

func (client *Client) GetId() string {
	return client.Id
}

func (client *Client) GetUrl() string {
	return client.Url
}

//------------------------------------------------------------------------------------------------------------

type Configuration struct {
	Dockerized  bool
	Environment string    `yaml:"environment"`
	Server      *Server   `yaml:"server"`
	Influx      *Influx   `yaml:"influx"`
	MySQL       *MySQL    `yaml:"mysql"`
	Clients     []*Client `yaml:"clients"`
}

func (configuration *Configuration) IsDockerized() bool {
	return configuration.Dockerized
}

func (configuration *Configuration) IsTestEnvironment() bool {
	return strings.ToLower(configuration.Environment) == "test"
}

func (configuration *Configuration) IsDevelopmentEnvironment() bool {
	return strings.ToLower(configuration.Environment) == "development"
}

func (configuration *Configuration) IsStagingEnvironment() bool {
	return strings.ToLower(configuration.Environment) == "staging"
}

func (configuration *Configuration) IsProductionEnvironment() bool {
	return strings.ToLower(configuration.Environment) == "production"
}

func (configuration *Configuration) GetEnvironment() string {
	return cases.Title(language.English, cases.Compact).String(strings.ToLower(configuration.Environment))
}

func (configuration *Configuration) GetServerConfiguration() IServerConfiguration {
	if configuration.Server == nil {
		configuration.Server = &Server{
			FQDN:     "localhost",
			Protocol: "http",
			Ports: &Ports{
				Active:      0,
				Passive:     0,
				Diagnostics: 0,
			},
			TLS: &TLS{
				KeyFile:  "",
				CertFile: "",
			},
		}
	}

	return configuration.Server
}

func (configuration *Configuration) GetInfluxConfiguration() IInfluxConfiguration {
	if configuration.Influx == nil {
		configuration.Influx = &Influx{
			Enabled:  false,
			Address:  "http://localhost:8086",
			Database: "",
			Username: "",
			Password: "",
		}
	}

	return configuration.Influx
}

func (configuration *Configuration) GetMySQLConfiguration() IMySqlConfiguration {
	if configuration.MySQL == nil {
		configuration.MySQL = &MySQL{
			Address:      "localhost:3306",
			Username:     "",
			Password:     "",
			SkipPassword: false,
		}
	}

	return configuration.MySQL
}

func (configuration *Configuration) GetClientsConfiguration() []IClientConfiguration {
	if configuration.Clients == nil {
		configuration.Clients = make([]*Client, 0)
	}

	result := make([]IClientConfiguration, 0)
	for _, client := range configuration.Clients {
		result = append(result, client)
	}

	return result
}

func (configuration *Configuration) GetPorts() (int, int, int) {
	ports := configuration.GetServerConfiguration().GetPortConfiguration()

	activePort := 7070
	if ports.GetActive() > 0 {
		activePort = ports.GetActive()
	}

	passivePort := 7080
	if ports.GetPassive() > 0 {
		passivePort = ports.GetPassive()
	}

	diagnostics := 6061
	if ports.GetDiagnostics() > 0 {
		diagnostics = ports.GetDiagnostics()
	}

	return activePort, passivePort, diagnostics
}

func NewConfiguration(path string, dockerized bool) (IConfiguration, error) {
	conf := &Configuration{
		Dockerized: dockerized,
	}

	if dockerized {
		conf.MySQL = &MySQL{
			Address:      fmt.Sprintf("%s:%s", os.Getenv("MYSQL_ADDRESS"), os.Getenv("MYSQL_PORT")),
			Database:     os.Getenv("MYSQL_DATABASE"),
			Username:     os.Getenv("MYSQL_USER"),
			Password:     os.Getenv("MYSQL_PASSWORD"),
			SkipPassword: false,
		}

		conf.Influx = &Influx{
			Enabled:  os.Getenv("INFLUX_ENABLED") == "true",
			Address:  os.Getenv("INFLUX_ADDRESS"),
			Database: os.Getenv("INFLUX_DATABASE"),
			Username: os.Getenv("INFLUX_USER"),
			Password: os.Getenv("INFLUX_PASSWORD"),
		}

		conf.Environment = os.Getenv("ENVIRONMENT")
	} else {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, err
		}

		configFile, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		if err := yaml.NewDecoder(configFile).Decode(conf); err != nil {
			return nil, err
		}
	}

	if conf.Server != nil {
		if os.Getenv("DIAMANTE_FQDN") != "" {
			conf.Server.FQDN = os.Getenv("DIAMANTE_FQDN")
		}

		if os.Getenv("DIAMANTE_PROTOCOL") != "" {
			conf.Server.Protocol = os.Getenv("DIAMANTE_PROTOCOL")
		}
	}

	switch conf.Environment {
	case "test", "development", "staging", "production":
	default:
		conf.Environment = "development"
	}

	return conf, nil
}

func NewTestConfiguration() IConfiguration {
	return &Configuration{
		Environment: "test",
		Server: &Server{
			FQDN:     "localhost",
			Protocol: "http",
			HashKey:  "OKq2gLmDCYJXnPweKrM=l7dFCDxp5Ff5EupcQCU",
			BlockKey: "v1K1s+S3vWudrypR",
		},
		Influx: &Influx{
			Enabled:  false,
			Address:  "http://localhost:8086",
			Database: "",
			Username: "",
			Password: "",
		},
		MySQL: &MySQL{
			Address:      "localhost:3306",
			Username:     "root",
			Password:     "password",
			SkipPassword: false,
		},
	}
}

func NewBenchmarkConfiguration() IConfiguration {
	return &Configuration{
		Environment: "test",
		Server: &Server{
			FQDN:     "localhost",
			Protocol: "http",
			HashKey:  "OKq2gLmDCYJXnPweKrM=l7dFCDxp5Ff5EupcQCU",
			BlockKey: "v1K1s+S3vWudrypR",
		},
		Influx: &Influx{
			Enabled:  false,
			Address:  "http://localhost:8086",
			Database: "",
			Username: "",
			Password: "",
		},
		MySQL: &MySQL{
			Address:      "localhost:3306",
			Username:     "root",
			Password:     "password",
			SkipPassword: false,
		},
	}
}
