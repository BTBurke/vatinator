package ocr

import (
	"fmt"
	"regexp"
	"strings"
)

var amt *regexp.Regexp
var gastype *regexp.Regexp

func init() {
	amt = regexp.MustCompile(`([0-9]+(\.|\,)[0-9]{1,3})\s?L`)
	gastype = regexp.MustCompile(`\s?(95|98)(\s|$)`)
}

func GasRule() Rule {
	return gas{}
}

type gas struct{}

func (gas) Find(r *Result, text []string) error {

	isGas, amount, gasType := detectGasReceipt(text)
	if isGas {
		r.Excise = &Excise{Type: fmt.Sprintf("Gasoline %s", gasType), Amount: amount}
	}

	return nil
}

func detectGasReceipt(text []string) (isGas bool, amount string, gasType string) {

	for _, line := range text {
		if strings.Contains(line, "EUR/L") && !isGas {
			isGas = true
		}
		c := amt.FindStringSubmatch(line)

		// results are like ["35.4 L", "35.4", "."]
		if len(c) >= 2 && amount == "" {
			// get rid of spaces and convert EU decimal to US decimal
			amount = strings.ReplaceAll(strings.ReplaceAll(c[1], " ", ""), ",", ".")
		}

		c2 := gastype.FindStringSubmatch(line)
		// results like ["Futura 95" "95"]
		if len(c2) >= 2 && gasType == "" {
			gasType = c2[1]
		}
	}

	// make sure that if data is there, isGas is true
	if amount != "" && gasType != "" {
		isGas = true
	}

	return
}
