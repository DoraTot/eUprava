package config

const (
	Auth0Domain = "dev-h1l4uuvj170yqf56.us.auth0.com"
	DBUser      = "root"
	DBPassword  = "secret"
	DBHost      = "db"
	DBName      = "e_uprava"
)

func GetDSN() string {
	return DBUser + ":" + DBPassword + "@tcp(" + DBHost + ":3306)/" + DBName
}
