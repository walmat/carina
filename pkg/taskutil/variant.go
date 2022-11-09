package taskutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Variant struct {
	ID      string `json:"id"`
	Size    string `json:"size"`
	InStock bool   `json:"inStock"`
}

type Predicate func(size string, variants []Variant) (Variant, error)

func MatchVariants(sizes []string, variants []Variant, predicate Predicate) ([]Variant, error) {
	//fmt.Printf("[DEBUG] Matching desired sizes: %s\n", strings.Join(sizes, ","))

	var matches []Variant
	var wantsRandom = false
	for _, size := range sizes {
		// NOTE: Check early if they want a random size if not matched any specific
		if strings.Contains(strings.ToLower(size), "random") {
			wantsRandom = true
		}

		variant, err := predicate(size, variants)
		if err == nil {
			matches = append(matches, variant)
		}
	}

	//fmt.Printf("[DEBUG] Matched %d variants\n", len(matches))
	//fmt.Printf("[DEBUG] User wants random: %t\n", wantsRandom)

	if len(matches) == 0 && !wantsRandom {
		return nil, errors.New("no variants matched")
	}

	if len(matches) == 0 && wantsRandom {
		rand.Seed(time.Now().Unix())

		n := rand.Intn(len(variants))

		choice, err := json.Marshal(variants[n])
		if err != nil {
			return nil, errors.New("unable to match random variant")
		}

		fmt.Printf("[DEBUG] Matched random variant: %s\n", choice)

		var matched []Variant
		matched = append(matched, variants[n])
		return matched, nil
	}

	return matches, nil
}
