package config

import (
	"errors"
	goflag "flag"
	"fmt"
	"os"
	"runtime"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const VERSION = "1.0.0"

// Config is the top level cofiguration object
type Config struct {
	Database struct {
		URI string `mapstructure:"uri" env:"RTCGW_DB" env-default:"postgres://postgres:postgres@localhost/rtcgw?sslmode=disable"`
	} `yaml:"database"`

	Server struct {
		Host                string `mapstructure:"host" env:"RTCGW_HOST" env-default:"localhost"`
		Port                string `mapstructure:"http_port" env:"RTCGW_SERVER_PORT" env-description:"Server port" env-default:"9292"`
		ProxyPort           string `mapstructure:"proxy_port" env:"RTCGW_PROXY_PORT" env-description:"Server port" env-default:"9191"`
		RedisAddress        string `mapstructure:"redis_address" env:"RTCGW_REDIS" env-description:"Redis address" env-default:"127.0.0.1:6379"`
		MigrationsDirectory string `mapstructure:"migrations_dir" env:"RTCGW_MIGRATTIONS_DIR" env-default:"file:///usr/share/rtcgw/db/migrations"`
	} `yaml:"server"`

	API struct {
		DHIS2BaseURL           string            `mapstructure:"dhis2_base_url" env:"DHIS2_BASE_URL" env-description:"The DHIS2 instance base API URL"`
		DHIS2User              string            `mapstructure:"dhis2_user"  env:"DHIS2_USER" env-description:"The DHIS2 username"`
		DHIS2Password          string            `mapstructure:"dhis2_password"  env:"DHIS2_PASSWORD" env-description:"The DHIS2  user password"`
		DHIS2PAT               string            `mapstructure:"dhis2_pat"  env:"DHIS2_PAT" env-description:"The DHIS2  Personal Access Token"`
		DHIS2AuthMethod        string            `mapstructure:"dhis2_auth_method"  env:"DHIS2_AUTH_METHOD" env-description:"The DHIS2 Authentication Method"`
		DHIS2TrackerProgram    string            `mapstructure:"dhis2_tracker_program"  env:"DHIS2_PROGRAM" env-description:"The DHIS2 tracker Program"`
		DHIS2TrackedEntityType string            `mapstructure:"dhis2_tracked_entity_type"  env:"DHIS2_TRACKED_ENTITY_TYPE" env-description:"The DHIS2 tracked entity type"`
		DHIS2SearchAttribute   string            `mapstructure:"dhis2_search_attribute" env:"DHIS_SEARCH_ATTRIBUTE" env-description:"The DHIS2 Search Attribute"`
		DHIS2Mapping           map[string]string `mapstructure:"dhis2_mapping" env:"DHIS_MAPPING" env-description:"The Request JSON keys mapping to DHIS2 Data Elements"`
	} `yaml:"api"`
}

var RTCGwConf Config
var ShowVersion *bool

func init() {
	var configFilePath, configDir string
	currentOS := runtime.GOOS
	switch currentOS {
	case "windows":
		configDir = "C:\\ProgramData\\Rtcgw"
		configFilePath = "C:\\ProgramData\\Rtcgw\\rtcgw.yml"
	case "darwin", "linux":
		configFilePath = "/etc/rtcgw/rtcgw.yml"
		configDir = "/etc/rtcgw/"
	default:
		fmt.Println("Unsupported operating system")
		return
	}

	configFile := flag.String("config-file", configFilePath,
		"The path to the configuration file of the application")

	ShowVersion = flag.Bool("version", false, "Display version of sukuma server")

	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()
	if *ShowVersion {
		fmt.Println("RTCGw: ", VERSION)
		os.Exit(1)
	}

	viper.SetConfigName("rtcgw")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	if len(*configFile) > 0 {
		viper.SetConfigFile(*configFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			panic(fmt.Errorf("Fatal Error %w \n", err))
		}
	}

	err := viper.Unmarshal(&RTCGwConf)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		err = viper.ReadInConfig()
		if err != nil {
			log.Fatalf("unable to reread configuration into global conf: %v", err)
		}
		_ = viper.Unmarshal(&RTCGwConf)
	})
	viper.WatchConfig()
}
