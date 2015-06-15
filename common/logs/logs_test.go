package logs

import (
	"os"
	"testing"
)

var logger = NewLogger(os.Stdout)

func TestSetLevel(t *testing.T) {
	SetLevel("trace")
}

func TestTrace(t *testing.T) {
	logger.SetLevel("trace")
	logger.Trace("trace")
	logger.SetLevel("off")
	logger.Trace("trace2")
	logger.SetLevel("debug")
	logger.Trace("trace3")

}

func TestTraceEnabled(t *testing.T) {
	logger.SetLevel("trace")
	if !logger.IsTraceEnabled() {
		t.FailNow()
		return
	}
}
