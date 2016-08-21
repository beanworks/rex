package rabbit

import (
	"io"
	"io/ioutil"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFailToCreateLogger(t *testing.T) {
	var err error

	c := &Config{
		Logger: struct {
			Output, Formatter, Level string
			LogFile                  string `mapstructure:"log_file"`
		}{
			Output: "stdout",
		},
	}

	c.Logger.Formatter = "Bilbo Baggins"
	_, err = NewLogger(c)

	assert.NotNil(t, err)
	assert.Error(t, err, "Unknown logger formatter: Bilbo Baggins")

	c.Logger.Formatter = "text"
	c.Logger.Level = "Gandalf the White"
	_, err = NewLogger(c)

	assert.NotNil(t, err)
	assert.Error(t, err, "Unkown logger level: Gandolf the White")
}

func TestCreateAndCloseLogger(t *testing.T) {
	tf, err := ioutil.TempFile("", "")
	require.Nil(t, err)

	c := &Config{
		Logger: struct {
			Output, Formatter, Level string
			LogFile                  string `mapstructure:"log_file"`
		}{
			Output:    "both",
			Formatter: "json",
			Level:     "info",
			LogFile:   tf.Name(),
		},
	}

	l, err := NewLogger(c)
	require.Nil(t, err)

	assert.Equal(t, c, l.config)
	assert.Equal(t, tf.Name(), l.file.Name())
	assert.Implements(t, (*io.Writer)(nil), l.Out)
	assert.IsType(t, new(log.JSONFormatter), l.Formatter)
	assert.Equal(t, log.InfoLevel, l.Level)

	_, err = l.file.Stat()
	assert.Nil(t, err)

	l.Close()

	_, err = l.file.Stat()
	assert.NotNil(t, err)
}
