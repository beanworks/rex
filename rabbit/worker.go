package rabbit

type Worker struct {
	Config *Config
	Logger *Logger
}

func NewWorker(c *Config, l *Logger) (*Worker, error) {
	return &Worker{}, nil
}

func (w *Worker) Consume() {
}
