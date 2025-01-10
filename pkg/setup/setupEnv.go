package setup

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// type Env string

// const (
// 	EnvProd Env = "prod"
// 	EnvTest Env = "test"
// )

// func CheckEnv() error {
// 	if len(os.Args) < 1 {
// 		return fmt.Errorf("no enviroment arg")
// 	}

// 	env := Env(os.Args[1])

// 	if (env != EnvProd) && (env != EnvTest) {
// 		return fmt.Errorf("unknown Env arg: %s", env)
// 	}

// 	log.Printf("Env: %s", env)
// 	return nil
// }

// func apiToken() string {
// 	env := Env(os.Args[1])
// 	switch env {
// 	case EnvProd:
// 		return ReadEnv("PROD_API_TOKEN")
// 	case EnvTest:
// 		return ReadEnv("TEST_API_TOKEN")
// 	default:
// 		return fmt.Sprintf("unknown env: %s", env)
// 	}
// }

func ReadEnv(name string) string {
	err := godotenv.Load()
	if err != nil {
		return ""
	}

	env := os.Getenv(name)

	log.Printf("%s env: %s", name, env)

	return env
}
