package providers

type DatabaseConfig struct {
	ConnectionString string
	MaxOpenConns     int
	MaxIdleConns     int
}
