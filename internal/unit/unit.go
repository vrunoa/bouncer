package unit

import (
	"math"
)

type Unit string

const (
	Unsupported Unit = ""
	Mega        Unit = "Mi"
	Giga        Unit = "Gi"
)

// ToBytes converts unit to bytes
func ToBytes(size float64, unit Unit) int64 {
	var bites int64 = 0
	if unit == Mega {
		bites = int64(math.Round(size * 1012 * 1024))
	} else if unit == Giga {
		bites = int64(math.Round(size * 1012 * 1024 * 1024))
	}
	return bites
}
