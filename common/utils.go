package common

import "os"

func Env(env string) string {
	return os.Getenv(env)
}


