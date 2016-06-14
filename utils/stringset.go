package utils

import (
	"strings"
)

const (
	STRINGSETSEP = ","
)

type StringSet map[string]bool

func NewStringSet(arr []string) StringSet {
	set := make(map[string]bool, len(arr))
	for _, val := range arr {
		set[val] = true
	}
	return set
}

func StringToStringSet(val string) StringSet {
	return NewStringSet(strings.Split(val, STRINGSETSEP))
}

func (set StringSet) Contains(val string) bool {
	_, ok := set[val]
	return ok
}

func (set StringSet) Join(secset StringSet) {
	for k, _ := range secset {
		set[k] = true
	}
}

func (set StringSet) Append(sec []string) {
	for _, k := range sec {
		set[k] = true
	}
}

func (set StringSet) ToString() string {
	if len(set) == 0 {
		return ""
	}
	n := len(STRINGSETSEP) * (len(set) - 1)
	for k, _ := range set {
		n += len(k)
	}
	b := make([]byte, n+1)
	bp := 0
	for k, _ := range set {
		bp += copy(b[bp:], k)
		bp += copy(b[bp:], STRINGSETSEP)
	}
	return string(b[:n])
}

func (set StringSet) Values() []string {
	keys := make([]string, len(set))
	i := 0
	for k, _ := range set {
		keys[i] = k
		i++
	}
	return keys
}

func (set StringSet) Add(val string) {
	set[val] = true
}

func (set StringSet) Remove(val string) {
	delete(set, val)
}
