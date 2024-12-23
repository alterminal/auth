package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/alterminal/auth/api"
	"github.com/alterminal/auth/repo"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	path := os.Getenv("config")
	if path == "" {
		path = "."
	}
	fmt.Println(path)
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	slog.SetLogLoggerLevel(slog.LevelDebug)
}

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.database"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	repo.Init(db)
	api.Run(db)
}
