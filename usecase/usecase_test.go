package usecase

import (
	"cron_parser/models"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// */15 0 1,15 * 1-5
// 0 0 * * * /usr/bin/find
// */15 * * * * /usr/bin/find
// 0 9-17 * * 1-5 /usr/bin/find
// 0 0,12 1,15 * * /usr/bin/find
// 1-5,9-25/2 0,12 1-15/2 1-6/2 0,6 /usr/bin/find
// * * * * * /usr/bin/find
// 0 12 1 1 1 /usr/bin/find
// 0 * 1,15 * 1-5 /usr/bin/find

// a b c d e /usr/bin/find
// 1     1     1    *   * /usr/bin/find
// 60 * * * * /usr/bin/find
// -1 * * * *

func TestCronUsecase_ValidateCronString(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expError string
	}{
		//{"invalid length", []string{"*", "*", "*", "*", "/usr/find"}, "invalid length for cron string"},
		//{"negative number", []string{"-1", "*", "*", "*", "*", "/usr/find"}, "val outside of range for Minute"},
		{"exceeding range", []string{"60", "*", "*", "*", "*", "/usr/find"}, "val outside of range for Minute"},
		{"invalid character", []string{"$", "*", "*", "*", "*", "/usr/find"}, "invalid character exists for Minute"},
		{"happy path", []string{"$", "*", "*", "*", "*", "mv", "cat.txt", "."}, "invalid character exists for Minute"},
	}

	for _, test := range tests {
		cronService := NewCronUsecase()
		err := cronService.ValidateCronString(test.args)
		if test.expError != "" {
			if err == nil {
				assert.Fail(t, fmt.Sprintf("did not raise error for %s", test.name))
			}
			assert.Equal(t, test.expError, err.Error())
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestCronUsecase_IsYear(t *testing.T) {
	tests := []struct {
		name   string
		args   string
		expOut bool
	}{
		{"happy flow with year", "2023-2025", true},
		{"command instead of year", "/usr/bin/find", false},
		{"only asterisk", "*", true},
		{"with division", "*/5", true},
		{"with invalid division", "*/usr", false},
		{"with valid separation", "2023,250", true},
		{"with valid separation", "mv,250", false},

		//{"invalid range", "-2025", },

	}

	for _, test := range tests {
		cronService := NewCronUsecase()
		result := cronService.IsYear(test.args)
		if result != test.expOut {
			assert.Fail(t, "output did not match")
		}
	}
}

func TestCronUsecase_ParseCronString(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expError string
		expResp  models.CronResponse
	}{
		{
			"simple cron with range, separation and division",
			[]string{"*/15", "0", "1,15", "*", "1-5"},
			"",
			models.CronResponse{
				Minute:     "0 15 30 45",
				Hour:       "0",
				DayOfMonth: "1 15",
				Month:      "1 2 3 4 5 6 7 8 9 10 11 12",
				DayOfWeek:  "1 2 3 4 5",
			},
		},
		{
			"simple cron with single digit",
			[]string{"1", "2", "3", "4", "5"},
			"",
			models.CronResponse{
				Minute:     "1",
				Hour:       "2",
				DayOfMonth: "3",
				Month:      "4",
				DayOfWeek:  "5",
			},
		},
		{
			"complex cron with combination of wildcard chars",
			[]string{"1-5,9-25/2", "0,12", "1-15/2", "1-6/2", "0,6"},
			"",
			models.CronResponse{
				Minute:     "1 2 3 4 5 9 11 13 15 17 19 21 23 25",
				Hour:       "0 12",
				DayOfMonth: "1 3 5 7 9 11 13 15",
				Month:      "1 3 5",
				DayOfWeek:  "0 6",
			},
		},
		{
			"invalid cron string with divided by 0",
			[]string{"2", "0,12", "1-15/2", "5/0", "0,6"},
			"invalid step 0",
			models.CronResponse{},
		},
		{
			"invalid cron string with divided by 0 in range",
			[]string{"2", "0,12", "1-15/2", "1-6/0", "0,6"},
			"invalid step 0",
			models.CronResponse{},
		},
		{
			"invalid cron string with divided by number",
			[]string{"2", "0,12", "1-15/2", "5/2", "0,6"},
			"not a valid one",
			models.CronResponse{},
		},
	}

	for _, test := range tests {
		cronService := NewCronUsecase()
		resp, err := cronService.ParseCronString(test.args)
		if test.expError != "" {
			if err == nil {
				assert.Fail(t, fmt.Sprintf("did not raise error for %s", test.name))
			}
			assert.Equal(t, test.expError, err.Error())
		} else {
			assert.Nil(t, err)
			assert.Equal(t, test.expResp.Minute, resp.Minute)
			assert.Equal(t, test.expResp.Hour, resp.Hour)
			assert.Equal(t, test.expResp.DayOfMonth, resp.DayOfMonth)
			assert.Equal(t, test.expResp.Month, resp.Month)
			assert.Equal(t, test.expResp.DayOfWeek, resp.DayOfWeek)
		}
	}
}

func TestCronUsecase_ParseRangeChar(t *testing.T) {
	tests := []struct {
		name   string
		args   string
		expOut string
	}{
		{"happy flow with normal range", "2-5", "2 3 4 5"},
		{"happy flow with exclusion range", "5-2", "2 3 4 5"},

		//{"invalid range", "-2025", },

	}

	for _, test := range tests {
		cronService := NewCronUsecase()
		result := cronService.IsYear(test.args)
		if result != test.expOut {
			assert.Fail(t, "output did not match")
		}
	}
}
