package config

const (
	DBUser     = "root"
	DBPassword = "secret"
	DBHost     = "db"
	DBName     = "e_uprava"
)

func GetDSN() string {
	return DBUser + ":" + DBPassword + "@tcp(" + DBHost + ":3306)/" + DBName
}
