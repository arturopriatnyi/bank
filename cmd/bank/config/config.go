package config

type Config struct {
	HTTPServer       HTTPServer       `env:",prefix=HTTP_SERVER_"`
	ExchangeRatesAPI ExchangeRatesAPI `env:",prefix=EXCHANGE_RATES_API_"`
	MongoDB          MongoDB          `env:",prefix=MONGODB_"`
	PostgreSQL       PostgreSQL       `env:",prefix=POSTGRESQL_"`
}

type HTTPServer struct {
	Addr string `env:"ADDR,default=localhost:9001"`
}

type ExchangeRatesAPI struct {
	AccessKey string `env:"ACCESS_KEY,default=67fbb39c33ab482b96f85ab232bf1c99"`
}

type MongoDB struct {
	URL      string `env:"URL,default=mongodb://localhost:27017"`
	Username string `env:"USERNAME,default=root"`
	Password string `env:"PASSWORD,default=password"`
	Database string `env:"DATABASE,default=bank"`
}

type PostgreSQL struct {
	Host     string `env:"HOST,default=localhost"`
	Port     int    `env:"PORT,default=5432"`
	User     string `env:"USER,default=postgres"`
	Password string `env:"PASSWORD,default=password"`
	Database string `env:"DATABASE,default=bank"`
	SSLMode  string `env:"SSL_MODE,default=disable"`
}
