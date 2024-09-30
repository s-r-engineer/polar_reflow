package logger

type GinWriter struct{}

func (gw GinWriter) Write(p []byte) (n int, err error) {
	Logger.Info(string(p))
	return len(p), nil
}

type GinErrWriter struct{}

func (gw GinErrWriter) Write(p []byte) (n int, err error) {
	Logger.Error(string(p))
	return len(p), nil
}
