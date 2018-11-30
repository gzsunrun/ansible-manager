package output_test

import (
	"testing"

	"github.com/gzsunrun/ansible-manager/core/config"
	"github.com/gzsunrun/ansible-manager/core/output"
)

func init() {
	config.SetLog("/var/log/test/log.FFF")
	config.NewConfig("/etc/ansible-manager/ansible-manager.conf")

}

func Test_Output(t *testing.T) {
	logger, err := output.NewLogOutput("ddddd")
	if err != nil {
		t.Error(err)
	}
	logger.Write("sddddddd")

}
