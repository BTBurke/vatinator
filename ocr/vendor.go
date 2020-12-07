package ocr

import (
	"regexp"
)

var v *regexp.Regexp

func init() {
	v = regexp.MustCompile(`[^/,]+\s(AS|TÜ|UÜ|OÜ|As|Tü|Uü|Oü|OU|Ou|TU|Tu|UU|Uu|0Ü|0u|0U|0ü)`)
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

// extracts the vendor name from the receipt
func extractVendor(lines []string) string {
	for _, line := range lines {
		c := v.FindAllStringSubmatch(line, -1)
		if len(c) > 0 && len(c[0][0]) > 0 {
			return c[0][0]
		}
	}
	return ""
}
