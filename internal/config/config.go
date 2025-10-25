package config

import (
	"flag"
	"os"
)

func LookupVar(envKey string, flagKey string) (string, bool) {
	env, exists := os.LookupEnv(envKey)
	if exists {
		return env, true
	}
	f := flag.Lookup(flagKey)
	if f != nil {
		return f.Value.String(), true
	}
	return "", false
}
