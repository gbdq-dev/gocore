package money

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/shopspring/decimal"
)

// ErrConversion indicates error during amount conversion.
var ErrConversion = errors.New("conversion error")

const (
	decimalPlaces = 2
	triadBase     = 1000
	hundredBase   = 100
	tenBase       = 10

	thousandsIndex = 1
	millionsIndex  = 2
	billionsIndex  = 3
)

// ToRussian converts money amount to Russian words.
func ToRussian(money decimal.Decimal) string {
	return convert(money, "ru")
}

// ToKazakh converts money amount to Kazakh words.
func ToKazakh(money decimal.Decimal) string {
	return convert(money, "kk")
}

// Internal implementation.
func convert(money decimal.Decimal, lang string) string {
	intPart, fracPart, err := parseMoney(money)
	if err != nil {
		return err.Error()
	}

	if intPart == 0 {
		return formatZero(lang, fracPart)
	}

	triads := splitIntoTriads(intPart)
	wordsSlice := processTriads(triads, lang)
	wordsStr := strings.Join(wordsSlice, " ")

	return formatResult(wordsStr, fracPart, lang)
}

func parseMoney(money decimal.Decimal) (int64, string, error) {
	str := money.StringFixed(decimalPlaces)
	parts := strings.Split(str, ".")

	var intPart int64
	if _, err := fmt.Sscan(parts[0], &intPart); err != nil {
		return 0, "", ErrConversion
	}

	return intPart, parts[1], nil
}

func splitIntoTriads(n int64) []int {
	var triads []int
	for n > 0 {
		triads = append(triads, int(n%triadBase))
		n /= triadBase
	}

	return triads
}

func processTriads(triads []int, lang string) []string {
	var words []string
	for triadIndex := len(triads) - 1; triadIndex >= 0; triadIndex-- {
		if triads[triadIndex] == 0 {
			continue
		}

		triadWords := convertTriad(triads[triadIndex], triadIndex, lang)
		if triadWords == "" {
			continue
		}

		words = append(words, triadWords)
		if unitForm := getUnitForm(triads[triadIndex], triadIndex, lang); unitForm != "" {
			words = append(words, unitForm)
		}
	}

	currencyWord := "тенге"
	words = append(words, currencyWord)

	return words
}

func convertTriad(num, triadIndex int, lang string) string {
	if num == 0 {
		return ""
	}

	hundreds := num / hundredBase
	tens := (num % hundredBase) / tenBase
	units := num % tenBase

	var parts []string
	addHundreds(&parts, hundreds, lang)
	addTensAndUnits(&parts, tens, units, triadIndex, lang)

	return strings.Join(parts, " ")
}

func addHundreds(parts *[]string, hundreds int, lang string) {
	if hundreds == 0 {
		return
	}

	switch lang {
	case "ru":
		ruHundreds := []string{
			"", "сто", "двести", "триста", "четыреста",
			"пятьсот", "шестьсот", "семьсот", "восемьсот", "девятьсот",
		}
		*parts = append(*parts, ruHundreds[hundreds])
	case "kk":
		kkHundreds := []string{
			"", "жүз", "екі жүз", "үш жүз", "төрт жүз",
			"бес жүз", "алты жүз", "жеті жүз", "сегіз жүз", "тоғыз жүз",
		}
		*parts = append(*parts, kkHundreds[hundreds])
	}
}

func addTensAndUnits(parts *[]string, tens, units, triadIndex int, lang string) {
	const (
		one = 1
		two = 2
	)

	switch lang {
	case "ru":
		ruUnits := []string{"", "один", "два", "три", "четыре", "пять", "шесть", "семь", "восемь", "девять"}
		ruTens := []string{"", "десять", "двадцать", "тридцать", "сорок", "пятьдесят", "шестьдесят", "семьдесят", "восемьдесят", "девяносто"}
		ruTeens := []string{"", "одиннадцать", "двенадцать", "тринадцать", "четырнадцать", "пятнадцать", "шестнадцать", "семнадцать", "восемнадцать", "девятнадцать"}

		if tens == one && units > 0 {
			*parts = append(*parts, ruTeens[units])

			return
		}

		if tens > 0 {
			*parts = append(*parts, ruTens[tens])
		}

		if units > 0 {
			isThousands := triadIndex == thousandsIndex
			if isThousands {
				switch units {
				case one:
					*parts = append(*parts, "одна")
				case two:
					*parts = append(*parts, "две")
				default:
					*parts = append(*parts, ruUnits[units])
				}
			} else {
				*parts = append(*parts, ruUnits[units])
			}
		}

	case "kk":
		kkUnits := []string{"", "бір", "екі", "үш", "төрт", "бес", "алты", "жеті", "сегіз", "тоғыз"}
		kkTens := []string{"", "он", "жиырма", "отыз", "қырық", "елу", "алпыс", "жетпіс", "сексен", "тоқсан"}

		if tens == one && units > 0 {
			*parts = append(*parts, kkTens[tens]+" "+kkUnits[units])

			return
		}

		if tens > 0 {
			*parts = append(*parts, kkTens[tens])
		}

		if units > 0 {
			*parts = append(*parts, kkUnits[units])
		}
	}
}

func getUnitForm(num, triadIndex int, lang string) string {
	switch lang {
	case "ru":
		switch triadIndex {
		case thousandsIndex:
			return pluralizeRu(num, []string{"тысяча", "тысячи", "тысяч"})
		case millionsIndex:
			return pluralizeRu(num, []string{"миллион", "миллиона", "миллионов"})
		case billionsIndex:
			return pluralizeRu(num, []string{"миллиард", "миллиарда", "миллиардов"})
		default:
			return ""
		}
	case "kk":
		switch triadIndex {
		case thousandsIndex:
			return "мың"
		case millionsIndex:
			return "миллион"
		case billionsIndex:
			return "миллиард"
		default:
			return ""
		}
	}

	return ""
}

//nolint:varnamelen,mnd
func pluralizeRu(n int, forms []string) string {
	if n%100 >= 11 && n%100 <= 19 {
		return forms[2]
	}

	const (
		oneForm   = 1
		twoToFour = 2
		three     = 3
		four      = 4
	)

	switch n % 10 {
	case oneForm:
		return forms[0]
	case twoToFour, three, four:
		return forms[1]
	default:
		return forms[2]
	}
}

func formatZero(lang, fracPart string) string {
	switch lang {
	case "ru":
		return "Ноль тенге " + fracPart + " тиын"
	case "kk":
		return "Нөл тенге " + fracPart + " тиын"
	default:
		return ""
	}
}

func formatResult(words, fracPart, lang string) string {
	if words == "" {
		return formatZero(lang, fracPart)
	}

	runes := []rune(words)
	if len(runes) > 0 {
		runes[0] = unicode.ToUpper(runes[0])
		words = string(runes)
	}

	return words + " " + fracPart + " тиын"
}
