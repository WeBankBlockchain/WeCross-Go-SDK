package wecrosslog

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/wecrosslog"
)

func TestLoggerV2Severity(t *testing.T) {
	buffers := []*bytes.Buffer{new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)}
	SetLoggerV1(NewLogger(buffers[infoLog], buffers[warningLog], buffers[errorLog]))

	wecrosslog.Logger.Info(severityName[infoLog])
	wecrosslog.Logger.Warning(severityName[warningLog])
	wecrosslog.Logger.Error(severityName[errorLog])

	for i := 0; i < fatalLog; i++ {
		buf := buffers[i]
		// The content of info buffer should be something like:
		// INFO: 2022/04/16 16:19:42 INFO
		// WARNING: 2022/04/16 16:19:42 WARNING
		// ERROR: 2022/04/16 16:19:42 ERROR
		for j := i; j < fatalLog; j++ {
			b, err := buf.ReadBytes('\n')
			if err != nil {
				t.Fatal(err)
			}
			if err = checkLogForSeverity(j, b); err != nil {
				t.Fatal(err)
			}
		}
	}
}

// check if b is in the format of:
// 2022/04/16 16:19:42 WARNING: WARNING
func checkLogForSeverity(s int, b []byte) error {
	expected := regexp.MustCompile(fmt.Sprintf(`^[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} %s: %s\n$`, severityName[s], severityName[s]))
	if m := expected.Match(b); !m {
		return fmt.Errorf("got: %v, want string in format of: %v", string(b), severityName[s]+": 2016/10/05 17:09:26 "+severityName[s])
	}
	return nil
}
