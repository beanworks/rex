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
			Durable    bool
			AutoDelete bool `mapstructure:"auto_delete"`
		}
		Prefetch struct {
			Count  int
			Global bool
		}
		Queue struct {
			Name       string
			Durable    bool
			AutoDelete bool `mapstructure:"auto_delete"`
		}
		Script string
	}
	Logger struct {
		LogToStderr     bool   `mapstructure:"log_to_stderr"`
		AlsoLogToStderr bool   `mapstructure:"also_log_to_stderr"`
		StderrThreshold string `mapstructure:"stderr_threshold"`
		LogDir          string `mapstructure:"log_dir"`
	}
}
