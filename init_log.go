package keychain

import (
	"gitlab.com/swarmfund/keychain/log"
)

// initLog initialized the logging subsystem, attaching app.log and
// app.logMetrics.  It also configured the logger's level using Config.LogLevel.
func initLog(app *App) {
	log.DefaultLogger.Logger.Level = app.Config().Log().Level
}

func init() {
	appInit.Add("log", initLog)
}
