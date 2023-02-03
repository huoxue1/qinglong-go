package log

import (
	log "github.com/sirupsen/logrus"
	log2 "xorm.io/xorm/log"
)

type MyLog struct {
}

func (m *MyLog) Debug(v ...interface{}) {
	log.Debug(v)
}

func (m *MyLog) Error(v ...interface{}) {
	log.Error(v)
}

func (m *MyLog) Info(v ...interface{}) {
	log.Info(v)
}

func (m *MyLog) Infof(format string, v ...interface{}) {
	log.Infof(format, v)
}

func (m *MyLog) Warn(v ...interface{}) {
	log.Warn(v)
}

func (m *MyLog) Level() log2.LogLevel {
	switch log.GetLevel() {
	case log.InfoLevel:
		return log2.LOG_INFO
	case log.DebugLevel:
		return log2.LOG_DEBUG
	case log.WarnLevel:
		return log2.LOG_WARNING
	case log.ErrorLevel:
		return log2.LOG_ERR
	default:
		return log2.LOG_UNKNOWN
	}
}

func (m *MyLog) SetLevel(l log2.LogLevel) {

}

func (m *MyLog) ShowSQL(show ...bool) {

}

func (m *MyLog) IsShowSQL() bool {
	return true
}

func (m MyLog) Errorf(format string, v ...interface{}) {
	log.Errorf(format, v)
}

func (m MyLog) Warnf(format string, v ...interface{}) {
	log.Warnf(format, v)
}

func (m MyLog) Debugf(format string, v ...interface{}) {
	log.Debugf(format, v)
}

type LogWriter struct {
}

func (l *LogWriter) Write(p []byte) (n int, err error) {
	log.Debugln(string(p))
	return len(p), nil
}

type CronLog struct {
}

func (c *CronLog) Info(msg string, keysAndValues ...interface{}) {
	log.Infoln(msg, "  ", keysAndValues)
}

func (c *CronLog) Error(err error, msg string, keysAndValues ...interface{}) {
	log.Errorln(err)
}
