package wgo

type Config struct {
	Service struct {
		Name      string `json:"name"`
		Version   string `json:"version"`
		DeployEnv string `json:"deploy_env"`
		Verbose   bool   `json:"verbose"`
	} `json:"service"`
	Logger struct {
		Level      string `json:"level"`
		Filename   string `json:"filename"`
		MaxSize    int    `json:"max_size"`
		MaxBackups int    `json:"max_backups"`
		MaxAge     int    `json:"max_age"`
		Compress   bool   `json:"compress"`
	} `json:"logger"`
}
