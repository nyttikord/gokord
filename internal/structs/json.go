package structs

import (
	"fmt"
	"strconv"
)

func SnowflakesToUints(sn []string) []uint64 {
	u := make([]uint64, len(sn))
	for i, v := range sn {
		var err error
		u[i], err = strconv.ParseUint(v, 10, 64)
		if err != nil {
			panic(err)
		}
	}
	return u
}

func UintsToSnowflakes(u []uint64) []string {
	sn := make([]string, len(u))
	for i, v := range u {
		sn[i] = fmt.Sprintf("%d", v)
	}
	return sn
}
