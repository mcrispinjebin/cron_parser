package usecase

import (
	"cron_parser/models"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/golang-collections/collections/set"
)

type CronUsecase interface {
	ValidateCronString(cronArgs []string) (err error)
	ParseCronString(cronStr []string, isYear bool, commandIndex int) (models.CronResponse, error)
	ParseEachField(fieldStr string, min, max int) (string, error)
	IsYear(cronStr string) bool
}

type cronUsecase struct{}

func (c cronUsecase) ValidateCronString(cronArgsArr []string) (err error) {
	if len(cronArgsArr) < 6 {
		err = fmt.Errorf("invalid length for cron string")
		return err
	}

	for ind, val := range cronArgsArr[:5] {
		config := models.CronConfig[ind]
		if c.isNumeric(val) {
			intVal, _ := strconv.Atoi(val)
			if intVal < config.Min || intVal > config.Max {
				err = fmt.Errorf("val outside of range for %s", config.Name)
				return
			}
			continue
		}
		splitFunc := func(r rune) bool {
			return r == ',' || r == '*' || r == '/' || r == '-'
		}
		tokensArr := strings.FieldsFunc(val, splitFunc)
		for _, token := range tokensArr {
			if !c.isNumeric(token) {
				err = fmt.Errorf("invalid character exists for %s", config.Name)
				return
			}

			intVal, _ := strconv.Atoi(token)
			if intVal < config.Min || intVal > config.Max {
				err = fmt.Errorf("val outside of range for %s", config.Name)
				return
			}
		}
	}
	return nil
}

func (c cronUsecase) IsYear(cronStr string) bool {
	splitFunc := func(r rune) bool {
		return r == ',' || r == '*' || r == '/' || r == '-'
	}
	tokensArr := strings.FieldsFunc(cronStr, splitFunc)

	for _, val := range tokensArr {
		if !c.isNumeric(val) {
			return false
		}
	}
	return true
}

func (c cronUsecase) ParseCronString(cronArgsArr []string, isYear bool, commandIndex int) (models.CronResponse, error) {
	// orchestrator
	// month day check, invalid check
	response := models.CronResponse{}
	v := reflect.ValueOf(&response).Elem()
	//cronArgsArr := strings.Split(cronStr, " ")[:5]

	for i := models.DayOfWeek; i >= models.Minute; i-- {
		config, _ := models.CronConfig[i]
		res, err := c.ParseEachField(cronArgsArr[i], config.Min, config.Max)
		if err != nil {
			return response, err
		}
		v.Field(i).SetString(res)
	}

	if isYear {
		config, _ := models.CronConfig[5]
		responseYear, err := c.ParseEachField(cronArgsArr[5], config.Min, config.Max)
		if err != nil {
			return response, err
		}

		response.Year = responseYear
	}

	response.Command = strings.Join(cronArgsArr[commandIndex:], " ")
	return response, nil
}

func (c cronUsecase) ParseEachField(fieldStr string, min, max int) (result string, err error) {
	//* , - /
	resultSet := set.New()
	if fieldStr == "*" {
		for i := min; i <= max; i++ {
			resultSet.Insert(fmt.Sprintf("%d", i))
		}
		result = strings.Join(setToSliceSorted(resultSet), " ")
		return
	}

	separationArr := strings.Split(fieldStr, ",")
	for _, fieldItem := range separationArr {
		if strings.Contains(fieldItem, "-") {
			//1-5 1-9/2
			rangeItemSet, err := c.ParseRangeChar(fieldItem, min, max)
			if err != nil {
				return "", err
			}
			resultSet = resultSet.Union(rangeItemSet)
		} else if strings.Contains(fieldItem, "/") || strings.Contains(fieldItem, "*") {
			// */2, *
			divisionSplitArr := strings.Split(fieldItem, "/")
			step := 1
			if len(divisionSplitArr) > 1 {
				step, _ = strconv.Atoi(divisionSplitArr[1])
				if step <= 0 {
					return "", fmt.Errorf("invalid step %d", step)
				}
			}

			if divisionSplitArr[0] == "*" {
				for i := min; i <= max; {
					resultSet.Insert(fmt.Sprintf("%d", i))
					i += step
				}

			} else {
				//1/10
				return result, fmt.Errorf("not a valid one")
			}
		} else if c.isNumeric(fieldItem) {
			// 1 2 -> just single digit
			resultSet.Insert(fieldItem)
		}
	}
	result = strings.Join(setToSliceSorted(resultSet), " ")
	return
}

func (c cronUsecase) ParseRangeChar(fieldStr string, configMin, configMax int) (result *set.Set, err error) {
	resultSet := set.New()

	isExclusion := false

	divisionSplitArr := strings.Split(fieldStr, "/")
	if !strings.Contains(divisionSplitArr[0], "-") {
		return nil, fmt.Errorf("invalid range")
	}
	rangeArr := strings.Split(divisionSplitArr[0], "-")
	minRange, _ := strconv.Atoi(rangeArr[0]) // 2   5
	maxRange, _ := strconv.Atoi(rangeArr[1]) // 5   2
	step := 1
	if len(divisionSplitArr) > 1 {
		step, _ = strconv.Atoi(divisionSplitArr[1])
		if step <= 0 {
			return nil, fmt.Errorf("invalid step %d", step)
		}
	}
	if minRange > maxRange {
		for i := configMin; i < configMax; i++ {
			exclusionSet.Insert(fmt.Sprintf("%d", i))
		}

	}

	for i := minRange; i <= maxRange; {
		resultSet.Insert(fmt.Sprintf("%d", i))
		i += step
	}

	// 2-5
	// 3-4
	if isExclusion && maxRange-minRange > 1 {
		exclusionSet := set.New()
		for i := minRange + 1; i < maxRange; i++ {
			exclusionSet.Insert(fmt.Sprintf("%d", i))
		}

		resultSet = resultSet.Difference(exclusionSet)
	}
	return resultSet, nil
}

func NewCronUsecase() CronUsecase {
	return &cronUsecase{}
}

func (c cronUsecase) isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func setToSliceSorted(s *set.Set) []string {
	result := make([]string, 0)
	s.Do(func(elem interface{}) {
		result = append(result, elem.(string))
	})
	sort.Slice(result, func(i, j int) bool {
		val1, _ := strconv.Atoi(result[i])
		val2, _ := strconv.Atoi(result[j])
		return val1 < val2
	})
	return result

}
