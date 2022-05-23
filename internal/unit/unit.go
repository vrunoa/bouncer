package unit

import (
	"errors"
	"math"
	"regexp"
	"strconv"
)

type Unit string

const (
	unsupported Unit = ""
	mega        Unit = "Mi"
	giga        Unit = "Gi"
)

func ParseSize(size string) (float64, Unit, error) {
	reg := regexp.MustCompile(`^(?P<Size>[0-9]+(\.[0-9])?)(?P<Unit>(Mi|Gi))`)
	if ok := reg.MatchString(size); !ok {
		return 0, unsupported, errors.New("invalid size")
	}
	match := reg.FindStringSubmatch(size)
	paramsMap := make(map[string]string)
	for i, name := range reg.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	sizeStr := paramsMap["Size"]
	s, err := strconv.ParseFloat(sizeStr, 64)
	if err != nil {
		return 0, unsupported, err
	}
	unitStr := paramsMap["Unit"]
	unt := unsupported
	if unitStr == string(giga) {
		unt = giga
	} else if unitStr == string(mega) {
		unt = mega
	}
	return s, unt, nil
}

func ToBytes(size float64, unit Unit) int64 {
	var bites int64 = 0
	if unit == mega {
		bites = int64(math.Round(size * 1012 * 1024))
	} else if unit == giga {
		bites = int64(math.Round(size * 1012 * 1024 * 1024))
	}
	return bites
}
