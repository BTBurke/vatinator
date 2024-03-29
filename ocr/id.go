package ocr

import (
	"regexp"
)

var kviitung *regexp.Regexp
var arve *regexp.Regexp
var hash *regexp.Regexp
var hash2 *regexp.Regexp
var nr *regexp.Regexp
var kvarve *regexp.Regexp
var tseki *regexp.Regexp
var boltUUID *regexp.Regexp
var wolt *regexp.Regexp
var telia *regexp.Regexp
var telia2 *regexp.Regexp
var alexela *regexp.Regexp

func init() {
	kviitung = regexp.MustCompile(`kviitung[^0-9]+([0-9]*\/?[0-9]*)?`)
	arve = regexp.MustCompile(`arve[^0-9]+([0-9]*)`)
	hash = regexp.MustCompile(`#([0-9]*)`)
	// for # that looks like h instead
	hash2 = regexp.MustCompile(`h([0-9]*)`)
	nr = regexp.MustCompile(`nr[^0-9]+([0-9]*)`)
	kvarve = regexp.MustCompile(`kv-arve[^0-9]+([0-9]*)`)
	tseki = regexp.MustCompile(`tšek[^0-9]+([0-9]*)`)
	// bolt uses UUIDv4, truncated to first two sections
	boltUUID = regexp.MustCompile(`document\sno\.\s([0-9a-f]{8}-[0-9a-f]{4})`)
	wolt = regexp.MustCompile(`order id\:?\s?([0-9a-f]+)`)
	// telia always starts with year
	telia = regexp.MustCompile(`invoice (2023[0-9]{10}).?`)
	telia2 = regexp.MustCompile(`(2023[0-9]{10})`)
	alexela = regexp.MustCompile(`.*arve ([0-9]+-[0-9]+)`)
}

type id struct{}

func (id) Find(r *Result, text []string) error {
	result := extractID(text)
	if result == "" {
		r.Errors = append(r.Errors, "no receipt number found")
		return nil
	}

	r.ID = result
	return nil
}

func IDRule() Rule {
	return id{}
}

// extracts the receipt id number, looking for either kviitung or arve
func extractID(lines []string) string {
	regexes := []*regexp.Regexp{
		alexela,
		kviitung,
		arve,
		hash,
		hash2,
		nr,
		kvarve,
		tseki,
		boltUUID,
		wolt,
		telia2,
		telia,
	}
	// this tries first with standard OCR output, then tries joining lines
	// again to see if you get any difference
	tries := [][]string{lines, joinFollowing(lines)}
	for _, try := range tries {
		for _, r := range regexes {
			if k := idFinder(r, try); k != "" {
				return k
			}
		}
	}
	return ""
}

// subroutine for executing a substring match for given regex
func idFinder(r *regexp.Regexp, lines []string) string {
	for _, line := range lines {
		k := r.FindAllStringSubmatch(line, -1)
		if len(k) > 0 && len(k[0]) == 2 {
			if len(k[0][1]) > 0 {
				return k[0][1]
			}
		}
	}
	return ""
}
