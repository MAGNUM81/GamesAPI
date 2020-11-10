package integration

import (
	"os"
)

//This is a fake Env loading. .env file are not mean to be loaded in a test configuration.
func SimulateEnv(){
	os.Setenv("STEAMKEY", "9230546D5E965861D940A995413DB4C8")
}
