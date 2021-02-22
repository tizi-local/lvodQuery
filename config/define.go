package config

type Config struct {
	VodQueryRpc *VodUploadRpc `json:"vodQueryRpc"`
	PprofAddr   string        `json:"pprof"`
	Redis       *RedisConfig  `json:"redis"`
	DB          *DBConfig     `json:"db"`
	Logger      *LoggerConfig `json:"logger"`
	RpcConfig   *RpcConfig `json:"rpc"`
}

// mysql 配置
type DBConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	DbName   string `json:"dbName"`
	Addr     string `json:"addr"`
}

// redis配置
type RedisConfig struct {
	Password string
	Addr     string
}

// rpc配置
type VodUploadRpc struct {
	Addr string `json:"addr"`
}

// 日志配置
type LoggerConfig struct {
	Path  string `json:"path"`
	Level string `json:"level"`
}

// 模块依赖调用rpc
type RpcConfig struct {
	Auth string `json:"auth"`
}