package config

import "github.com/spf13/viper"

type (
	Config struct {
		App        string `mapstructure:"APP"`
		Env        string `mapstructure:"ENV"`
		SvcVersion string `mapstructure:"SVC_VERSION"`
		Port       string `mapstructure:"PORT"`

		// Postgres
		PostgreHost         string `mapstructure:"POSTGRES_URL"`
		DBMaxOpenConnection int    `mapstructure:"DB_MAX_OPEN_CONN"`
		DBMaxIdleConnection int    `mapstructure:"DB_MAX_IDLE_CONN"`

		// Redis
		RedisDB       int      `mapstructure:"REDIS_DB"`
		RedisHost     []string `mapstructure:"REDIS_URL"`
		RedisUsername string   `mapstructure:"REDIS_USERNAME"`
		RedisPassword string   `mapstructure:"REDIS_PASSWORD"`

		// HTTP client
		HttpClientTimeout             int  `mapstructure:"HTTP_CLIENT_TIMEOUT"`
		HttpClientDisableKeepAlives   bool `mapstructure:"HTTP_CLIENT_DISABLE_KEEP_ALIVE"`
		HttpClientMaxIdleConns        int  `mapstructure:"HTTP_CLIENT_MAX_IDLE_CONNS"`
		HttpClientMaxConnsPerHost     int  `mapstructure:"HTTP_CLIENT_MAX_CONNS_PER_HOST"`
		HttpClientMaxIdleConnsPerHost int  `mapstructure:"HTTP_CLIENT_MAX_IDLE_CONNS_PER_HOST"`
		HttpClientIdleConnTimeout     int  `mapstructure:"HTTP_CLIENT_IDLE_CONN_TIMEOUT"`

		// Logging
		LogDirectory  string `mapstructure:"LOG_DIR"`
		LogFileName   string `mapstructure:"LOG_FILENAME"`
		LogConsole    bool   `mapstructure:"LOG_CONSOLE"`
		LogMaxSize    int    `mapstructure:"LOG_MAX_SIZE"`
		LogMaxAge     int    `mapstructure:"LOG_MAX_AGE"`
		LogMaxBackups int    `mapstructure:"LOG_MAX_BACKUP"`

		// Datadog
		DatadogAgentHost string `mapstructure:"DATADOG_AGENT_HOST"`
	}
)

func NewConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	redisHost := viper.GetStringSlice("REDIS_URL")
	cfg.RedisHost = redisHost
	return &cfg, err
}
