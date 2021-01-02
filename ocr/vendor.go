package ocr

import (
	"regexp"
	"strings"
)

var v *regexp.Regexp
var v2 *regexp.Regexp

func init() {
	// company symbol at end
	v = regexp.MustCompile(`[^/,]+\s(AS|TÜ|UÜ|OÜ|As|Tü|Uü|Oü|OU|Ou|TU|Tu|UU|Uu|0Ü|0u|0U|0ü)`)
	// company symbol at front
	v2 = regexp.MustCompile(`(AS|TÜ|UÜ|OÜ|OÙ|As|Tü|Uü|Oü|OU|Ou|TU|Tu|UU|Uu|0Ü|0u|0U|0ü)\s[^/,]+$`)
}

type vendor struct{}

func (vendor) Find(r *Result, text []string) error {
	result := extractVendor(text)
	if result == "" {
		r.Errors = append(r.Errors, "no vendor found")
		return nil
	}

	r.Vendor = result

	return nil
}

func VendorRule() Rule {
	return vendor{}
}

func extractVendor(lines []string) string {
	regexes := []*regexp.Regexp{v, v2}
	for _, r := range regexes {
		out := extract(r, lines)
		if len(out) > 0 {
			return finalFixes(out)
		}
	}
	return ""
}

// extracts the vendor name from the receipt
func extract(r *regexp.Regexp, lines []string) string {
	for _, line := range lines {
		c := r.FindAllStringSubmatch(line, -1)
		if len(c) > 0 && len(c[0][0]) > 0 {
			return c[0][0]
		}
	}
	return ""
}

func finalFixes(s string) string {
	// fixes OCR mistakes like O->0 and Ü->Ù
	replacements := map[string]string{
		"Ù": "U",
		"0": "O",
	}
	for from, to := range replacements {
		s = strings.ReplaceAll(s, from, to)
	}
	return s
}
