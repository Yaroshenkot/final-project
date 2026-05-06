package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

// NextDate вычисляет следующую дату по правилу repeat
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("правило повторения не указано")
	}

	startDate, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", errors.New("неверный формат даты")
	}

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", errors.New("неверный формат правила")
	}

	switch parts[0] {
	case "y":
		return nextDateYear(now, startDate)
	case "d":
		if len(parts) != 2 {
			return "", errors.New("неверный формат правила d")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("неверный формат числа дней (1-400)")
		}
		return nextDateDays(now, startDate, days)
	case "w":
		if len(parts) != 2 {
			return "", errors.New("неверный формат правила w")
		}
		return nextDateWeek(now, startDate, parts[1])
	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("неверный формат правила m")
		}
		return nextDateMonth(now, startDate, parts[1:])
	default:
		return "", errors.New("неподдерживаемый формат правила")
	}
}

// nextDateYear добавляет годы
func nextDateYear(now time.Time, date time.Time) (string, error) {
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			return date.Format(DateFormat), nil
		}
	}
}

// nextDateDays добавляет дни
func nextDateDays(now time.Time, date time.Time, days int) (string, error) {
	for {
		date = date.AddDate(0, 0, days)
		if afterNow(date, now) {
			return date.Format(DateFormat), nil
		}
	}
}

// nextDateWeek — правило w (дни недели)
func nextDateWeek(now time.Time, date time.Time, weekdaysStr string) (string, error) {
	// Парсим дни недели (1-7)
	parts := strings.Split(weekdaysStr, ",")
	var weekdays []int
	for _, p := range parts {
		d, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil || d < 1 || d > 7 {
			return "", errors.New("неверный день недели (1-7)")
		}
		weekdays = append(weekdays, d)
	}

	// Ищем следующую подходящую дату
	for {
		date = date.AddDate(0, 0, 1)
		if !afterNow(date, now) {
			continue
		}
		weekday := int(date.Weekday())
		if weekday == 0 {
			weekday = 7 // воскресенье = 7
		}
		for _, wd := range weekdays {
			if weekday == wd {
				return date.Format(DateFormat), nil
			}
		}
	}
}

// nextDateMonth — правило m (дни месяца)
func nextDateMonth(now time.Time, date time.Time, args []string) (string, error) {
	// Парсим дни месяца
	daysStr := args[0]
	dayParts := strings.Split(daysStr, ",")

	var targetDays []int // положительные = конкретные дни, отрицательные = с конца
	var months []int     // опциональные месяцы

	for _, p := range dayParts {
		p = strings.TrimSpace(p)
		d, err := strconv.Atoi(p)
		if err != nil {
			return "", errors.New("неверный формат дня месяца")
		}
		if d < -2 || d > 31 || d == 0 {
			return "", errors.New("день месяца должен быть от 1 до 31 или -1, -2")
		}
		targetDays = append(targetDays, d)
	}

	// Если указаны месяцы
	if len(args) == 2 {
		monthsStr := args[1]
		monthParts := strings.Split(monthsStr, ",")
		for _, p := range monthParts {
			p = strings.TrimSpace(p)
			m, err := strconv.Atoi(p)
			if err != nil || m < 1 || m > 12 {
				return "", errors.New("неверный месяц (1-12)")
			}
			months = append(months, m)
		}
	}

	// Ищем следующую подходящую дату
	for {
		date = date.AddDate(0, 0, 1)
		if !afterNow(date, now) {
			continue
		}

		// Проверяем месяц, если указан
		if len(months) > 0 {
			monthMatch := false
			currentMonth := int(date.Month())
			for _, m := range months {
				if currentMonth == m {
					monthMatch = true
					break
				}
			}
			if !monthMatch {
				continue
			}
		}

		// Проверяем день
		currentDay := date.Day()
		lastDayOfMonth := daysInMonth(date.Year(), int(date.Month()))

		for _, td := range targetDays {
			switch {

			case td > 0 && currentDay == td:
				return date.Format(DateFormat), nil
			case td == -1 && currentDay == lastDayOfMonth:
				return date.Format(DateFormat), nil
			case td == -2 && currentDay == lastDayOfMonth-1:
				return date.Format(DateFormat), nil
			}
		}
	}
}

// daysInMonth возвращает количество дней в месяце
func daysInMonth(year int, month int) int {
	return time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// afterNow проверяет, что date > now (без учёта времени)
func afterNow(date, now time.Time) bool {
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	nowOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return dateOnly.After(nowOnly)
}
