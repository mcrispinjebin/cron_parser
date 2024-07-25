package handler

import (
	"cron_parser/usecase"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func CliHandler(cronUsecase usecase.CronUsecase) {
	input := os.Args[1:]
	if len(input) == 0 {
		fmt.Printf("error: cron string is not given")
		return
	}

	pattern := `\s+`
	regex := regexp.MustCompile(pattern)
	input[0] = regex.ReplaceAllString(input[0], " ")

	args := strings.Split(input[0], " ")
	if err := cronUsecase.ValidateCronString(args); err != nil {
		fmt.Printf("error: %s", err.Error())
		return
	}

	commandInd := 5
	isYear := cronUsecase.IsYear(args[5])
	if isYear {
		commandInd = 6
	}

	response, err := cronUsecase.ParseCronString(args, isYear, commandInd)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		return
	}

	fmt.Printf("minute: %s\nhour: %s\nday of month: %s\nmonth: %s\nday of week: %s\nyear: %s\ncommand: %s",
		response.Minute,
		response.Hour,
		response.DayOfMonth,
		response.Month,
		response.DayOfWeek,
		response.Year,
		response.Command)
}
