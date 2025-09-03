package conf

import (
	"fmt"
	"os"
	"path"
	"strconv"
)

func GetDataPath() string {
	env := os.Getenv("IRBANKMOCK_DATA_PATH")

	if env == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "."
		}
		return cwd
	}

	return env
}

func GetDbFileName() string {
	env := os.Getenv("IRBANKMOCK_DB_NAME")
	if env == "" {
		return "irbankmock.db"
	}
	return env
}

func GetDbPath() string {
	return path.Join(GetDataPath(), GetDbFileName()+"?_pragma=foreign_keys(1)")
}

func GetListenAddress() string {
	env := os.Getenv("IRBANKMOCK_SERVER_PORT")
	if env == "" {
		return ":3000"
	}
	return env
}

func ShouldAutoMigrate() bool {
	env := os.Getenv("IRBANKMOCK_AUTOMIGRATE")
	val, err := strconv.ParseBool(env)
	if err != nil {
		if env != "" {
			fmt.Println("invalid value for IRBANKMOCK_AUTOMIGRATE value, falling-back to default=true")
		}
		return true
	}
	return val
}

func IsGormLogDisabled() bool {
	env := os.Getenv("IRBANKMOCK_DISABLE_GORM_LOG")
	val, _ := strconv.ParseBool(env)
	return val
}

func GetWebAppPath() string {
	env := os.Getenv("IRBANKMOCK_WEBAPP_PATH")
	if env == "" {
		return "./web/app/out"
	}
	return env
}

func GetPublicHostname() string {
	env := os.Getenv("IRBANKMOCK_PUBLIC_HOSTNAME")
	if env == "" {
		return "misconfig.example.com"
	}
	return env
}
