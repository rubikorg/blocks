package logger

import (
	"fmt"
	"time"

	"github.com/rubikorg/rubik"
)

const (
	// BlockName is the name of the block
	BlockName = "HTTPLogger"
)

// BlockLogger implements Block interface of rubik
type BlockLogger struct {
	// config blockLoggerConfig
}

// TODO: make this a simple logger for projects
// type blockLoggerConfig struct {
// 	Level  string `json:"level"`
// 	Format string `json:"format"`
// }

// OnAttach implementation of rubik block
func (bl BlockLogger) OnAttach(app *rubik.App) error {
	rubik.BeforeRequest(bl.beforeHook)
	rubik.AfterRequest(bl.afterHook)
	return nil
}

func (bl BlockLogger) beforeHook(rc *rubik.RequestContext) {
	rc.Ctx["from"] = time.Now()
}

func (bl BlockLogger) afterHook(rc *rubik.RequestContext) {
	fromTime, _ := rc.Ctx["from"].(time.Time)
	responseTime := time.Since(fromTime).Seconds() * 100000
	layout := "Mon, 2 Jan 2006 15:04:05 MST"
	logTime := time.Now().Format(layout)
	logMsg := fmt.Sprintf("[%s] Method:%s - [%s] [%d] - %fms", logTime, rc.Request.Method,
		rc.Request.URL.Path, rc.Status, responseTime)
	fmt.Println(logMsg)
}

func init() {
	rubik.Attach(BlockName, BlockLogger{})
}