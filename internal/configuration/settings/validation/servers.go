package validation

import (
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

func sortedInsert(ss []string, s string) []string {
	i := sort.SearchStrings(ss, s)
	ss = append(ss, "")
	copy(ss[i+1:], ss[i:])
	ss[i] = s
	return ss
}

func ExtractCountries(servers []models.Server) (values []string) {
	seen := make(map[string]struct{}, len(servers))
	values = make([]string, 0, len(servers))
	for _, server := range servers {
		value := server.Country
		_, alreadySeen := seen[value]
		if alreadySeen {
			continue
		}
		seen[value] = struct{}{}

		values = sortedInsert(values, value)
	}
	return values
}

func ExtractRegions(servers []models.Server) (values []string) {
	seen := make(map[string]struct{}, len(servers))
	values = make([]string, 0, len(servers))
	for _, server := range servers {
		value := server.Region
		_, alreadySeen := seen[value]
		if alreadySeen {
			continue
		}
		seen[value] = struct{}{}

		values = sortedInsert(values, value)
	}
	return values
}

func ExtractCities(servers []models.Server) (values []string) {
	seen := make(map[string]struct{}, len(servers))
	values = make([]string, 0, len(servers))
	for _, server := range servers {
		value := server.City
		_, alreadySeen := seen[value]
		if alreadySeen {
			continue
		}
		seen[value] = struct{}{}

		values = sortedInsert(values, value)
	}
	return values
}

func ExtractISPs(servers []models.Server) (values []string) {
	seen := make(map[string]struct{}, len(servers))
	values = make([]string, 0, len(servers))
	for _, server := range servers {
		value := server.ISP
		_, alreadySeen := seen[value]
		if alreadySeen {
			continue
		}
		seen[value] = struct{}{}

		values = sortedInsert(values, value)
	}
	return values
}

func ExtractServerNames(servers []models.Server) (values []string) {
	seen := make(map[string]struct{}, len(servers))
	values = make([]string, 0, len(servers))
	for _, server := range servers {
		value := server.ServerName
		_, alreadySeen := seen[value]
		if alreadySeen {
			continue
		}
		seen[value] = struct{}{}

		values = sortedInsert(values, value)
	}
	return values
}

func ExtractHostnames(servers []models.Server) (values []string) {
	seen := make(map[string]struct{}, len(servers))
	values = make([]string, 0, len(servers))
	for _, server := range servers {
		value := server.Hostname
		_, alreadySeen := seen[value]
		if alreadySeen {
			continue
		}
		seen[value] = struct{}{}

		values = sortedInsert(values, value)
	}
	return values
}
