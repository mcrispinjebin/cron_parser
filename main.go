package main

import (
	"cron_parser/handler"
	"cron_parser/usecase"
)

func setUp() {
	cronUsecase := usecase.NewCronUsecase()
	handler.CliHandler(cronUsecase)
}

func main() {
	setUp()
}
