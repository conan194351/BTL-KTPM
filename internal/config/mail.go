package config

type Mail struct {
	Host        string `env:"HOST" envDefault:"localhost"`
	Port        int    `env:"PORT" envDefault:"1025"`
	SenderEmail string `env:"SENDER_EMAIL" envDefault:""`
	Password    string `env:"PASSWORD" envDefault:""`
}
