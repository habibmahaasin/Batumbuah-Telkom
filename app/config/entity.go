package config

type App struct {
	Name        string `env:"APP_NAME"`
	Port        string `env:"APP_PORT"`
	Mode        string `env:"APP_MODE"`
	Url         string `env:"APP_URL"`
	Secret_key  string `env:"APP_SECRET"`
}
type Db struct {
	Host string `env:"DB_HOST"`
	Name string `env:"DB_NAME"`
	User string `env:"DB_USER"`
	Pass string `env:"DB_PASSWORD"`
	Port string `env:"DB_PORT"`
}

type Storage struct {
	Access_Key string `env:"STORAGE_ACCESS_KEY"`
	Private_Key string `env:"STORAGE_SECRET_KEY"`
	Region     string `env:"STORAGE_REGION"`
	Bucket     string `env:"STORAGE_BUCKET"`
	Endpoint  string `env:"STORAGE_ENDPOINT"`
}
type BasicAuth struct {
	Username string `env:"BASIC_AUTH_USER"`
	Password string `env:"BASIC_AUTH_PASSWORD"`
}
type Conf struct {
	App       App
	Db        Db
	BasicAuth BasicAuth
	Storage   Storage
}
