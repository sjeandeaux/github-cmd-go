package os

import goos "os"

//Getenv get the environmental value, otherwise default value.
func Getenv(name string, defaultValue string) string {
	envValue := goos.Getenv(name)
	if envValue == "" {
		return defaultValue
	}
	return envValue
}
