package models

import "fmt"

type Site struct {
	Url string
	EntriesCount int
}

func (s *Site) ChangeEntries(NewCount int) {
	s.EntriesCount = NewCount
}
func (s *Site) PrintCounts() {
	fmt.Printf("Count for %s: %d\n", s.Url, s.EntriesCount)
}
