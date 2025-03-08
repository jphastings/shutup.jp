package mapping

import (
	"github.com/biter777/countries"
	"github.com/jphastings/dotpostcard/types"
)

type Countries map[string]struct{}

func (list Countries) Codes() []string {
	var codes []string
	for code := range list {
		codes = append(codes, code)
	}
	return codes
}

func (list Countries) Add(l types.Location) error {
	c := countries.ByName(l.CountryCode)
	list[c.Alpha2()] = struct{}{}
	return nil
}
