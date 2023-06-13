package settings

import (
	"os"
	"strconv"
	"strings"

	"github.com/jinzhu/configor"
	. "github.com/xeronith/diamante/contracts/settings"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Server struct {
	FQDN               string `yaml:"fqdn"`
	Protocol           string `yaml:"protocol"`
	Ports              *Ports `yaml:"ports"`
	TLS                *TLS   `yaml:"tls"`
	BuildNumber        int32  `yaml:"build_number"`
	JwtTokenKey        string `yaml:"jwt_token_key"`
	JwtTokenExpiration string `yaml:"jwt_token_expiration"`
	HashKey            string `yaml:"hash_key"`
	BlockKey           string `yaml:"block_key"`
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

func (server *Server) GetJwtTokenKey() string {
	return server.JwtTokenKey
}

func (server *Server) GetJwtTokenExpiration() string {
	if strings.TrimSpace(server.JwtTokenExpiration) == "" {
		server.JwtTokenExpiration = "10h"
	}

	return server.JwtTokenExpiration
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

type PostgreSQL struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (postgres *PostgreSQL) GetHost() string {
	if postgres.Host == "" {
		postgres.Host = "127.0.0.1"
	}

	return postgres.Host
}

func (postgres *PostgreSQL) SetHost(host string) {
	postgres.Host = host
}

func (postgres *PostgreSQL) GetPort() string {
	if postgres.Port == "" {
		postgres.Port = "5432"
	}

	return postgres.Port
}

func (postgres *PostgreSQL) SetPort(port string) {
	postgres.Port = port
}

func (postgres *PostgreSQL) GetDatabase() string {
	return postgres.Database
}

func (postgres *PostgreSQL) SetDatabase(database string) {
	postgres.Database = database
}

func (postgres *PostgreSQL) GetUsername() string {
	return postgres.Username
}

func (postgres *PostgreSQL) SetUsername(username string) {
	postgres.Username = username
}

func (postgres *PostgreSQL) GetPassword() string {
	return postgres.Password
}

func (postgres *PostgreSQL) SetPassword(password string) {
	postgres.Password = password
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

type Configuration struct {
	Dockerized     bool
	Environment    string      `yaml:"environment"`
	TrafficRecord  bool        `yaml:"traffic_record"`
	AllowedOrigins []string    `yaml:"allowed_origins"`
	Server         *Server     `yaml:"server"`
	Influx         *Influx     `yaml:"influx"`
	PostgreSQL     *PostgreSQL `yaml:"postgres"`
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

func (configuration *Configuration) IsTrafficRecordEnabled() bool {
	return configuration.TrafficRecord
}

func (configuration *Configuration) GetAllowedOrigins() []string {
	return configuration.AllowedOrigins
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
			Address:  "http://127.0.0.1:8086",
			Database: "",
			Username: "",
			Password: "",
		}
	}

	return configuration.Influx
}

func (configuration *Configuration) GetPostgreSQLConfiguration() IPostgreSQLConfiguration {
	if configuration.PostgreSQL == nil {
		configuration.PostgreSQL = &PostgreSQL{
			Host:     "127.0.0.1",
			Port:     "5432",
			Username: "",
			Password: "",
		}
	}

	return configuration.PostgreSQL
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
		Server:     &Server{},
	}

	if dockerized {
		conf.Server = &Server{
			FQDN:     "localhost",
			Protocol: "http",
			Ports: &Ports{
				Active:      7070,
				Passive:     7080,
				Diagnostics: 6061,
			},
		}

		conf.PostgreSQL = &PostgreSQL{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			Database: os.Getenv("POSTGRES_DATABASE"),
			Username: os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
		}

		conf.Influx = &Influx{
			Enabled:  os.Getenv("INFLUX_ENABLED") == "true",
			Address:  os.Getenv("INFLUX_ADDRESS"),
			Database: os.Getenv("INFLUX_DATABASE"),
			Username: os.Getenv("INFLUX_USER"),
			Password: os.Getenv("INFLUX_PASSWORD"),
		}

		if conf.Influx.Address == "" {
			conf.Influx.Address = "http://localhost:8086"
		}

		conf.Environment = os.Getenv("ENVIRONMENT")
	} else {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, err
		}

		if err := configor.Load(conf, path); err != nil {
			return nil, err
		}
	}

	if conf.Server != nil {
		if os.Getenv("FQDN") != "" {
			conf.Server.FQDN = os.Getenv("FQDN")
		}

		if os.Getenv("PROTOCOL") != "" {
			conf.Server.Protocol = os.Getenv("PROTOCOL")
		}

		if os.Getenv("PORT") != "" {
			port, err := strconv.Atoi(os.Getenv("PORT"))
			if err != nil {
				return nil, err
			}

			conf.Server.Ports = &Ports{
				Active:      7070,
				Passive:     port,
				Diagnostics: 6061,
			}
		}

		if os.Getenv("JWT_TOKEN_KEY") != "" {
			conf.Server.JwtTokenKey = os.Getenv("JWT_TOKEN_KEY")
		}

		if os.Getenv("JWT_TOKEN_EXP") != "" {
			conf.Server.JwtTokenExpiration = os.Getenv("JWT_TOKEN_EXP")
		}

		if os.Getenv("SECURE_COOKIE_HASH_KEY") != "" {
			conf.Server.HashKey = os.Getenv("SECURE_COOKIE_HASH_KEY")
		}

		if os.Getenv("SECURE_COOKIE_BLOCK_KEY") != "" {
			conf.Server.BlockKey = os.Getenv("SECURE_COOKIE_BLOCK_KEY")
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
			FQDN:               "localhost",
			Protocol:           "http",
			JwtTokenKey:        "",
			JwtTokenExpiration: "10h",
			HashKey:            "",
			BlockKey:           "",
		},
		Influx: &Influx{
			Enabled:  false,
			Address:  "http://127.0.0.1:8086",
			Database: "",
			Username: "",
			Password: "",
		},
		PostgreSQL: &PostgreSQL{
			Host:     "127.0.0.1",
			Port:     "5432",
			Username: "postgres",
			Password: "password",
		},
	}
}

func NewBenchmarkConfiguration() IConfiguration {
	return &Configuration{
		Environment: "test",
		Server: &Server{
			FQDN:               "localhost",
			Protocol:           "http",
			JwtTokenKey:        "",
			JwtTokenExpiration: "10h",
			HashKey:            "",
			BlockKey:           "",
		},
		Influx: &Influx{
			Enabled:  false,
			Address:  "http://127.0.0.1:8086",
			Database: "",
			Username: "",
			Password: "",
		},
		PostgreSQL: &PostgreSQL{
			Host:     "127.0.0.1",
			Port:     "5432",
			Username: "postgres",
			Password: "password",
		},
	}
}
