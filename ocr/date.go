package ocr

import (
	"fmt"
	"regexp"
)

var d *regexp.Regexp

func init() {
	d = regexp.MustCompile(`(01|02|03|04|05|06|07|08|09|10|11|12|13|14|15|16|17|18|19|20|21|22|23|24|25|26|27|28|29|30|31)\s?\.?\,?\/?\-?\s?(01|02|03|04|05|06|07|08|09|10|11|12)\s?\.?\,?\/?\-?\s?(2020|2021|20|21)`)
}

type date struct{}

func (date) Find(r *Result, text []string) error {
	date := extractDate(text)
	if date == "" {
		r.Errors = append(r.Errors, "no date found")
		return nil
	}

	r.Date = date
	return nil
}

// finds all dates of the form ddmmyy dd.mm.yy dd.mm.yyyy ddmmyyyy
func extractDate(raw []string) string {
	for _, line := range raw {
		r := d.FindAllStringSubmatch(line, -1)
		if len(r) > 0 && len(r[0]) == 4 {
			if len(r[0][3]) == 2 {
				r[0][3] = fmt.Sprintf("20%s", r[0][3])
			}
			return fmt.Sprintf("%s/%s/%s", r[0][1], r[0][2], r[0][3])
		}
	}
	return ""
}

func DateRule() Rule {
	return date{}
}
