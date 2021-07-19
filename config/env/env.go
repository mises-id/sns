package env

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

var Envs *Env

type Env struct {
	Port              int    `env:"PORT" envDefault:"8080"`
	AppEnv            string `env:"APP_ENV" envDefault:"development"`
	LogLevel          string `env:"LOG_LEVEL" envDefault:"INFO"`
	MongoURI          string `env:"MONGO_URI,required"`
	DBName            string `env:"DB_NAME" envDefault:"mises"`
	AssetLogoBasePath string `env:"ASSET_LOGO_BASE_PATH" envDefault:"http://localhost/assets"`
}

func init() {
	fmt.Println("env initializing...")
	_, b, _, _ := runtime.Caller(0)
	appEnv := os.Getenv("APP_ENV")
	projectRootPath := filepath.Dir(b) + "/../../"
	envPath := projectRootPath + ".env"
	appEnvPath := envPath + "." + appEnv
	localEnvPath := appEnvPath + ".local"
	_ = godotenv.Load(envPath, appEnvPath, localEnvPath)
	Envs = &Env{}
	err := env.Parse(Envs)
	if err != nil {
		panic(err)
	}
	fmt.Println("env loaded...")
}
