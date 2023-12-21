/*Question-01*/
package main

type Config struct {
	DBHost     string `envconfig:"DB_HOST" default:"localhost"`
	DBPort     int    `envconfig:"DB_PORT" default:"5432"`
	DBUser     string `envconfig:"DB_USER" default:"postgres"`
	DBPassword string `envconfig:"DB_PASSWORD" default:"password"`
	DBName     string `envconfig:"DB_NAME" default:"yourdb"`
}
