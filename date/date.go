package date

import "time"

const layout = "2006-01-02"

func ParseFromString(input string) (time.Time, error) {
	return time.Parse(layout, input)
}

func ParseToString(input time.Time) string {
	return input.Format(layout)
}
