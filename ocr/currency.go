package ocr

import (
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var curr2 *regexp.Regexp
var curr3 *regexp.Regexp

func init() {
	curr2 = regexp.MustCompile(`[0-9]+(\,|\.)\s?[0-9]{2}`)
	curr3 = regexp.MustCompile(`[0-9]+\,\s?[0-9]{3}\s?$?`)
}

type currency struct{}

func (currency) Find(r *Result, text []string) error {
	tax, total, _ := findTaxTotal(text)
	if tax == 0 && total == 0 {
		r.Errors = append(r.Errors, "no tax/total found")
		return nil
	}

	r.Total = total
	r.VAT = tax
	return nil
}

func CurrencyRule() Rule {
	return currency{}
}

// findTaxTotal returns the tax, total or 0,0 if not found
func findTaxTotal(text []string) (int, int, CurrencyPrecision) {
	currencies := extractCurrency3(text)
	currencies = append(currencies, extractCurrency2(text)...)
	precision := Currency2

	tax, total := extractTaxTotal(currencies)

	return tax, total, precision
}

// extracts all numbers of the form dd+,ddd and returns them as integers in unit values (x100) to a 2-digit precision
func extractCurrency3(raw []string) []int {
	out := make([]int, 0)
	for _, line := range raw {
		lineT := strings.Trim(line, "€*EUR eur")
		c := curr3.FindAllString(lineT, -1)
		for _, c1 := range c {
			cUnit := strings.Replace(c1, ",", "", -1)
			cUnit = strings.Replace(cUnit, ".", "", -1)
			cUnit = strings.Replace(cUnit, " ", "", -1)
			// lop off last digit of 3-digit currencies
			cAsInt, err := strconv.Atoi(cUnit[0 : len(cUnit)-2])
			if err != nil {
				continue
			}
			out = append(out, cAsInt)
		}
	}
	return out
}

// extracts all numbers of the form dd+,dd and returns them as integers in unit values (x100)
func extractCurrency2(raw []string) []int {
	out := make([]int, 0)
	for _, line := range raw {
		lineT := strings.Trim(line, "€*EUR eur")
		c := curr2.FindAllString(lineT, -1)
		for _, c1 := range c {
			cUnit := strings.Replace(c1, ",", "", -1)
			cUnit = strings.Replace(cUnit, ".", "", -1)
			cUnit = strings.Replace(cUnit, " ", "", -1)
			cAsInt, err := strconv.Atoi(cUnit)
			if err != nil {
				continue
			}
			out = append(out, cAsInt)
		}
	}
	return out
}

// determine tax and total by checking for 20% tax for every number on receipt
// only works because the values are sorted and it starts looking at the number most likely to be total
// TODO: doesn't handle the 9% or 10% tax brackets
func extractTaxTotal(in []int) (tax int, total int) {
	sort.Ints(in)
	for i := len(in)/2 - 1; i >= 0; i-- {
		opp := len(in) - 1 - i
		in[i], in[opp] = in[opp], in[i]
	}

	maxCost := 0
	if len(in) > 0 {
		maxCost = in[0]
	}
	for _, i := range in {
		total = i
		expectedTax := total - int(float64(total)/1.20)
		for _, j := range in {
			if j >= expectedTax-1 && j <= expectedTax+1 {
				tax = j
				if math.Abs(float64(maxCost-total)) <= 10 {
					total = maxCost
				}
				return
			}
		}
	}
	return 0, 0
}
