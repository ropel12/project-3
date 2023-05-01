package config

import "github.com/spf13/viper"

type Server struct {
	Port string `mapstructure:"PORT"`
}
type DatabaseConfig struct {
	Host     string `mapstructure:"HOST"`
	Port     string `mapstructure:"PORT"`
	Username string `mapstructure:"USERNAME"`
	Password string `mapstructure:"PASSWORD"`
	Name     string `mapstructure:"NAME"`
}
type RedisConfig struct {
	Host     string `mapstructure:"HOST"`
	Port     string `mapstructure:"PORT"`
	Password string `mapstructure:"PASSWORD"`
	DB       int    `mapstructure:"DB"`
}

type GCPConfig struct {
	Credential string `mapstructure:"CREDEN"`
	PRJID      string `mapstructure:"PROJECTID"`
	BCKNM      string `mapstructure:"BUCKETNAME"`
	Path       string `mapstructure:"PATH"`
}

type MidtransConfig struct {
	ServerKey      string `mapstructure:"SERVERKEY"`
	ClientKey      string `mapstructure:"CLIENTKEY"`
	Env            int    `mapstructure:"ENV"`
	URLHandler     string `mapstructure:"URL"`
	ExpiryDuration int    `mapstructure:"EXP"`
	Unit           string `mapstructure:"UNIT"`
}
type NSQConfig struct {
	Host   string `mapstructure:"HOST"`
	Port   string `mapstructure:"PORT"`
	Topic  string `mapstructure:"TOPIC"`
	Topic2 string `mapstructure:"TOPIC2"`
	Topic3 string `mapstructure:"TOPIC3"`
	Topic4 string `mapstructure:"TOPIC4"`
}
type PusherConfig struct {
	AppId   string `mapstructure:"APPID"`
	Key     string `mapstructure:"KEY"`
	Secret  string `mapstructure:"SECRET"`
	Cluster string `mapstructure:"CLUSTER"`
	Secure  bool   `mapstructure:"SECURE"`
	Channel string `mapstructure:"CHANNEL"`
	Event   string `mapstructure:"EVENT"`
}
type Config struct {
	Server     Server         `mapstructure:"SERVER"`
	Database   DatabaseConfig `mapstructure:"DATABASE"`
	Midtrans   MidtransConfig `mapstructure:"MIDTRANS"`
	JwtSecret  string         `mapstructure:"JWTSECRET"`
	Redis      RedisConfig    `mapstructure:"REDIS"`
	CSRFLength int            `mapstructure:"CSRFLENGTH"`
	CSRFMode   string         `mapstructure:"CSRFMODE"`
	NSQ        NSQConfig      `mapstructure:"NSQ"`
	GCP        GCPConfig      `mapstructure:"GCP"`
	Pusher     PusherConfig   `mapstructure:"PUSHER"`
}

func InitConfiguration() (*Config, error) {
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	viper.AutomaticEnv()
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
