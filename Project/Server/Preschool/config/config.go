package config

const (
	Auth0Domain = "dev-h1l4uuvj170yqf56.us.auth0.com"
	DBUser      = "root"
	DBPassword  = "secret"
	DBHost      = "db"
	DBName      = "e_uprava"
	Audience    = "https://dev-h1l4uuvj170yqf56.us.auth0.com/api/v2/"
	Auth0Client = "xSMUVFfxYctMHEjtZ8SEKBGL2ad4grQp"
	Auth0Secret = "uOjR2S2wm4HczgzS8xkGKQzzL5lgEr03yG8C7tmYz6tmBdB4BMTSL07zP7u1lrSi"
)

func GetDSN() string {
	return DBUser + ":" + DBPassword + "@tcp(" + DBHost + ":3306)/" + DBName
}
