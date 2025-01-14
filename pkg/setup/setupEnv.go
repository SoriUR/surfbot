package setup

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env string

const (
	EnvProd Env = "prod"
	EnvTest Env = "test"
)

func checkEnv() error {
	if len(os.Args) < 1 {
		return fmt.Errorf("no enviroment arg. expected \"test\" or \"prod\"")
	}

	env := Env(os.Args[1])

	if (env != EnvProd) && (env != EnvTest) {
		return fmt.Errorf("unknown first env arg: %s. expected \"test\" or \"prod\"", env)
	}

	log.Printf("Env: %s", env)
	return nil
}

func ApiToken() string {
	checkEnv()

	env := Env(os.Args[1])
	switch env {
	case EnvProd:
		return ReadEnv("PROD_API_TOKEN")
	case EnvTest:
		return ReadEnv("TEST_API_TOKEN")
	default:
		return fmt.Sprintf("unknown env: %s", env)
	}
}

func ReadEnv(name string) string {
	err := godotenv.Load()
	if err != nil {
		return ""
	}

	env := os.Getenv(name)

	log.Printf("%s env: %s", name, env)

	return env
}
