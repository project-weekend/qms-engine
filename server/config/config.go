package config

// Config is the main configuration structure for qms-engine service
type Config struct {
	Name        string      `json:"name"`
	ServiceName string      `json:"serviceName"`
	Env         string      `json:"env"`
	Host        string      `json:"host"`
	Port        int         `json:"port"`
	OwnerInfo   OwnerInfo   `json:"ownerInfo"`
	Data        Data        `json:"data"`
	RedisConfig RedisConfig `json:"redisConfig"`
	Statsd      Statsd      `json:"statsd"`
	Trace       Trace       `json:"trace"`
	Logger      Logger      `json:"logger"`
	Kafka       KafkaConfig `json:"kafka"`
}

// OwnerInfo contains information about the service owner
type OwnerInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	URL   string `json:"url"`
}

// Data contains all data source configurations
type Data struct {
	MySQL MySQL `json:"mysql"`
}

// MySQL contains MySQL database configurations
type MySQL struct {
	Master               DBConfig       `json:"master"`
	Slave                DBConfig       `json:"slave"`
	MasterCircuitBreaker CircuitBreaker `json:"masterCircuitBreaker"`
	SlaveCircuitBreaker  CircuitBreaker `json:"slaveCircuitBreaker"`
}

// DBConfig contains database connection configuration
type DBConfig struct {
	DSN             string `json:"dsn"`
	MaxIdle         int    `json:"maxIdle"`
	MaxOpen         int    `json:"maxOpen"`
	ConnMaxLifetime string `json:"connMaxLifetime"`
}

// CircuitBreaker contains circuit breaker configuration
type CircuitBreaker struct {
	TimeoutInMs            int `json:"timeoutInMs"`
	MaxConcurrentReq       int `json:"maxConcurrentReq"`
	VolumePercentThreshold int `json:"volumePercentThreshold"`
}

// RedisConfig contains Redis configuration
type RedisConfig struct {
	Addr               string `json:"addr"`
	IdleTimeoutInSec   int    `json:"idleTimeoutInSec"`
	PoolSize           int    `json:"poolSize"`
	ReadOnlyFromSlaves bool   `json:"readOnlyFromSlaves"`
	ReadTimeoutInSec   int    `json:"readTimeoutInSec"`
	WriteTimeoutInSec  int    `json:"writeTimeoutInSec"`
	TLSEnabled         bool   `json:"tlsEnabled"`
}

// Statsd contains StatsD configuration
type Statsd struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// Trace contains tracing configuration
type Trace struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Disable bool   `json:"disable"`
}

// Logger contains logging configuration
type Logger struct {
	WorkerCount     int    `json:"workerCount"`
	BufferSize      int    `json:"bufferSize"`
	LogLevel        int    `json:"logLevel"`
	StacktraceLevel int    `json:"stacktraceLevel"`
	LogFormat       string `json:"logFormat"`
}

// KafkaConfig contains Kafka configuration
type KafkaConfig struct {
	BootstrapServers string `json:"bootstrap.servers"`
	GroupID          string `json:"group.id"`
	AutoOffsetReset  string `json:"auto.offset.reset"`
	ProducerEnabled  bool   `json:"producer.enabled"`
}

// LoadConfig loads configuration from viper and unmarshals into Config struct
func LoadConfig() *Config {
	v := NewViper()

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		panic(err)
	}

	return &config
}
