package model

import (
	"regexp"
	"strings"
)

// move it to a yaml/json file
var DictCommonWords = map[string]bool{
	"the":        true,
	"be":         true,
	"to":         true,
	"of":         true,
	"and":        true,
	"a":          true,
	"in":         true,
	"that":       true,
	"have":       true,
	"I":          true,
	"it":         true,
	"for":        true,
	"not":        true,
	"on":         true,
	"with":       true,
	"he":         true,
	"as":         true,
	"you":        true,
	"do":         true,
	"at":         true,
	"this":       true,
	"but":        true,
	"his":        true,
	"by":         true,
	"from":       true,
	"they":       true,
	"we":         true,
	"say":        true,
	"her":        true,
	"she":        true,
	"or":         true,
	"an":         true,
	"will":       true,
	"my":         true,
	"one":        true,
	"all":        true,
	"would":      true,
	"there":      true,
	"their":      true,
	"what":       true,
	"so":         true,
	"up":         true,
	"out":        true,
	"if":         true,
	"about":      true,
	"who":        true,
	"get":        true,
	"which":      true,
	"go":         true,
	"me":         true,
	"when":       true,
	"make":       true,
	"can":        true,
	"like":       true,
	"time":       true,
	"no":         true,
	"just":       true,
	"him":        true,
	"know":       true,
	"take":       true,
	"people":     true,
	"into":       true,
	"year":       true,
	"your":       true,
	"good":       true,
	"some":       true,
	"could":      true,
	"them":       true,
	"see":        true,
	"other":      true,
	"than":       true,
	"then":       true,
	"now":        true,
	"look":       true,
	"only":       true,
	"come":       true,
	"its":        true,
	"over":       true,
	"think":      true,
	"also":       true,
	"back":       true,
	"after":      true,
	"use":        true,
	"two":        true,
	"how":        true,
	"our":        true,
	"work":       true,
	"first":      true,
	"well":       true,
	"way":        true,
	"even":       true,
	"new":        true,
	"want":       true,
	"because":    true,
	"any":        true,
	"these":      true,
	"give":       true,
	"day":        true,
	"most":       true,
	"us":         true,
	"here":       true,
	"such":       true,
	"much":       true,
	"yet":        true,
	"very":       true,
	"every":      true,
	"many":       true,
	"is":         true,
	"am":         true,
	"got":        true,
	"are":        true,
	"more":       true,
	"online":     true,
	"best":       true,
	"why":        true,
	"while":      true,
	"without":    true,
	"try":        true,
	"everything": true,
	"been":       true,
	"true":       true,
	"actual":     true,
	"actually":   true,
	"wouldn":     true,
	"couldn":     true,
	"haven":      true,
	"hasn":       true,
	"ain":        true,
	"between":    true,
	"themselves": true,
	"thing":      true,
	"nothing":    true,
	"things":     true,
	"man":        true,
	"being":      true,
	"has":        true,
	"never":      true,
	"must":       true,
	"were":       true,
	"was":        true,
	"wasn":       true,
	"weren":      true,
	"regarding":  true,
	"around":     true,
	"either":     true,
	"itself":     true,
	"himself":    true,
	"herself":    true,
	"myself":     true,
	"yourself":   true,
	"within":     true,
	"same":       true,
	"cannot":     true,
	"apart":      true,
	"where":      true,
	"proper":     true,
	"properly":   true,
	"short":      true,
	"long":       true,
	"shortest":   true,
	"longest":    true,
	"hit":        true,
	"told":       true,
	"favorite":   true,
	"isn":        true,
	"wants":      true,
	"ever":       true,
	"hard":       true,
	"hardest":    true,
	"don":        true,
	"lot":        true,
}

func FilterWords(list []string, dictionaries ...map[string]bool) []string {
	result := []string{}

	for _, el := range list {
		el = strings.ToLower(el)

		if len(el) < 2 || DictCommonWords[el] {
			continue
		}

		skip := false
		for _, dict := range dictionaries {
			if dict[el] {
				skip = true
				break
			}
		}

		if skip {
			continue
		}

		result = append(result, el)
	}

	return result
}

func TokenizeContent(input string, dictionaries ...map[string]bool) []string {
	return FilterWords(regexp.MustCompile("[^\\w]").Split(input, -1), dictionaries...)
}
