package rabbit

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCmd(t *testing.T) {
	msg := []byte("Test Script")

	cfg := &Config{}
	cfg.Consumer.Worker.Script = "echo"

	scr := Script{cfg}
	out, err := scr.ExecWith(msg)
	require.Nil(t, err)

	str := string(out)
	decoded, err := base64.StdEncoding.DecodeString(str)
	require.Nil(t, err)

	assert.Equal(t, msg, decoded)
}
