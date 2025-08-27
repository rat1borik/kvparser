// package domain Мейн-логика
package domain

import (
	"kvparser/internal/utils"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type DoctorOptions struct {
	Subdivision  string
	Specialities []string
	Specialists  []string
}

type MatchResult struct {
	Matches []DoctorMatch
}

type DoctorMatch struct {
	Name        string
	Speciality  string
	Status      string
	Subdivision string
}

func FindMatches(html string, options DoctorOptions) (*MatchResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	nodes := make([]DoctorMatch, 0)

	doc.Find(".docsInLpuTableDetail").Children().Each(func(i int, s *goquery.Selection) {
		addr := s.Find(".docsInLpuTableDescr > .lpu-address").Text()

		s.Children().Filter(".doc-row-ha").Each(func(i int, s1 *goquery.Selection) {
			el := DoctorMatch{
				Name:        strings.TrimSpace(s1.Find(".doc-name-mo > a").Text()),
				Speciality:  strings.TrimSpace(s1.Find(".profile-mo").Text()),
				Status:      strings.TrimSpace(s1.Find(".nearest-record-mo").Text()),
				Subdivision: strings.TrimSpace(addr),
			}

			// Фильтрация
			if el.Status != "" && el.Status != "Запись через интернет недоступна" && !utils.IsEmptyOrWhitespace(el.Status) &&
				(len(options.Specialists) == 0 || slices.Contains(options.Specialists, el.Name)) &&
				(len(options.Specialities) == 0 || slices.Contains(options.Specialities, el.Speciality)) &&
				(options.Subdivision == "" || options.Subdivision == el.Subdivision) {
				nodes = append(nodes, el)
			}
		})
	})

	return &MatchResult{
		Matches: nodes,
	}, nil

}
