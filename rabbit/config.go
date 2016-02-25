package rabbit

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
		Script string
	}
	Logger struct {
		Output    string
		Formatter string
		Level     string
		LogFile   string `mapstructure:"log_file"`
	}
}
