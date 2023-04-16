package kana

import (
	"strings"
	"unicode"
)

// ToRomaji converts hiragana and/or katakana to lowercase romaji. By default,
// the literal transliteration of づ　and ぢ are used, returnin du and di,
// respectively. Set vocalize to true to return the romaji in its correctly
// pronounced form - zu and ji.
func ToRomaji(s string, vocalize bool) string {
	s = ToRomajiCased(s, vocalize)
	s = strings.ToLower(s)

	return s
}

// ToRomajiCased converts hiragana and/or katakana to cased romaji, where
// hiragana and katakana are presented in lowercase and uppercase respectively.
func ToRomajiCased(s string, vocalize bool) string {
	s = moraicNRomaji.Replace(s)
	s = kanaToRomaji.Replace(s)
	s = parseRomajiDoubles([]rune(s))
	s = postRomaji.Replace(s)
	s = postRomajiSpecial.Replace(s)

	if vocalize {
		s = vocalizedRomaji.Replace(s)
	}
	return s
}

// ToHiragana converts wapuro-hepburn romaji into the equivalent hiragana.
func ToHiragana(s string) string {
	s = strings.ToLower(s)
	s = preHiragana.Replace(s)
	s = romajiToHiragana.Replace(s)
	s = strings.Map(KatakanaToHiragana, s)
	s = postHiragana.Replace(s)
	return postKanaSpecial.Replace(s)
}

// ToKatakana converts wapuro-hepburn romaji into the equivalent katakana.
func ToKatakana(s string) string {
	s = strings.ToUpper(s)
	s = preKatakana.Replace(s)
	s = romajiToKatakana.Replace(s)
	s = strings.Map(HiraganaToKatakana, s)
	s = postKatakana.Replace(s)
	return postKanaSpecial.Replace(s)
}

// ToKana converts wapuro-hepburn uppercase and lowercase romaji into
// katakana and hiragana respectively.
func ToKana(s string) string {
	s = preHiragana.Replace(s)
	s = preKatakana.Replace(s)
	s = romajiToHiragana.Replace(s)
	s = romajiToKatakana.Replace(s)
	s = postHiragana.Replace(s)
	s = postKatakana.Replace(s)
	return postKanaSpecial.Replace(s)
}

// HiraganaToKatakana replaces a single hiragana character with the
// unicode equivalent katakana character.
func HiraganaToKatakana(r rune) rune {
	if (r >= 'ぁ' && r <= 'ゖ') || (r >= 'ゝ' && r <= 'ゞ') {
		return r + 0x60
	}
	return r
}

// KatakanaToHiragana replaces a single katakana character with the
// unicode equivalent hiragana character.
func KatakanaToHiragana(r rune) rune {
	if (r >= 'ァ' && r <= 'ヶ') || (r >= 'ヽ' && r <= 'ヾ') {
		return r - 0x60
	}
	return r
}

// IsKatakana returns true if every element of a string is katakana, except
// for characters indicated in sanitizeIsChecks (spaces and dashes).
func IsKatakana(s string) bool {
	s = sanitizeIsChecks.Replace(s)
	if s == "" {
		return false
	}

	for _, r := range s {
		if !unicode.In(r, unicode.Katakana) {
			return false
		}
	}

	return true
}

// IsHiragana returns true if every element of a string is hiragana, except
// for characters indicated in sanitizeIsChecks (spaces and dashes).
func IsHiragana(s string) bool {
	s = sanitizeIsChecks.Replace(s)
	if s == "" {
		return false
	}

	for _, r := range s {
		if !unicode.In(r, unicode.Hiragana) {
			return false
		}
	}

	return true
}

// IsKanji returns true if every element of a string is a kanji character,
// except for characters indicated in sanitizeIsChecksKanji (spaces).
func IsKanji(s string) bool {
	s = sanitizeIsChecksKanji.Replace(s)
	if s == "" {
		return false
	}

	for _, r := range s {
		if !unicode.In(r, unicode.Han) {
			return false
		}
	}

	return true
}

// ContainsKatakana returns true if a string contains any katakana characters.
func ContainsKatakana(s string) bool {
	for _, r := range s {
		if unicode.In(r, unicode.Katakana) {
			return true
		}
	}

	return false
}

// ContainsHiragana returns true if a string contains any hiragana characters.
func ContainsHiragana(s string) bool {
	for _, r := range s {
		if unicode.In(r, unicode.Hiragana) {
			return true
		}
	}

	return false
}

// ContainsKanji returns true if a string contains any kanji characters.
func ContainsKanji(s string) bool {
	for _, r := range s {
		if unicode.In(r, unicode.Han) {
			return true
		}
	}

	return false
}

// ExtractKanji returns a slice containing all kanji characters found in a
// string, in the order in which they were found. If a kanji exists multiple
// times in a string, then each instance of the kanji will be returned.
func ExtractKanji(s string) []string {
	k := []string{}
	for _, r := range s {
		if unicode.In(r, unicode.Han) {
			k = append(k, string(r))
		}
	}
	return k
}
