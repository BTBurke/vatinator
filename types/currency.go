package types

import (
	"fmt"
	"strconv"
)

type Currency int

func (c Currency) String() string {
	cS := strconv.Itoa(int(c))
	switch {
	case c >= 0 && c < 10:
		return fmt.Sprintf("0.0%d", c)
	case c >= 10 && c < 100:
		return fmt.Sprintf("0.%d", c)
	default:
		return fmt.Sprintf("%s.%s", cS[0:len(cS)-2], cS[len(cS)-2:])
	}
}
