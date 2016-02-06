package rabbit

type Config struct {
	Connection struct {
		Host     string
		Username string
		Password string
		Vhost    string
		Port     int
	}
	Worker struct {
		Exchange struct {
			Name       string
			Type       string
			AutoDelete bool `mapstructure:"auto_delete"`
			Durable    bool
		}
		Prefetch struct {
			Count  int
			Global bool
		}
		Queue  string
		Script string
	}
	Logger struct {
		Appenders struct {
			File struct {
				Enabled bool
				Path    string
			}
			Echo struct {
				Enabled bool
			}
		}
	}
}
