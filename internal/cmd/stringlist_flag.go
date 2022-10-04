package cmd

import "strings"

// StringList is a slice of strings that satisfies the interface for a
// flags.Value.
type StringList []string

func (stringlist *StringList) String() string {
	return "stringlist"
}

func (stringlist *StringList) Set(s string) error {
	values := strings.Split(s, ",")
	*stringlist = append(*stringlist, values...)
	return nil
}
