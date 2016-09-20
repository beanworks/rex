package rabbit

// Config defines the configuration for Rex. It helps transform the yml config
// to a Go struct, which will be used and referenced directly in Logger and Rex
// structs.
type Config struct {
	Connection struct {
		Host     string
		Username string
		Password string
		Vhost    string
		Port     int
	}
	Consumer struct {
		Exchange struct {
			Name       string
			Type       string
			Durable    bool
			AutoDelete bool `mapstructure:"auto_delete"`
		}
		Prefetch struct {
			Count  int
			Global bool
		}
		Queue struct {
			Name       string
			RoutingKey string `mapstructure:"routing_key"`
			Durable    bool
			AutoDelete bool `mapstructure:"auto_delete"`
		}
		Worker struct {
			Script        string
			Count         int
			RetryInterval int `mapstructure:"retry_interval"`
		}
	}
	Logger struct {
		Output    string
		Formatter string
		Level     string
		LogFile   string `mapstructure:"log_file"`
	}
}
