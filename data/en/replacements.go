package en

var (
	LeetSpeakCharacters = map[rune]rune{
		// Common
		'!': 'i',
		'@': 'a',
		'4': 'a',
		'8': 'b',
		'6': 'b',
		'(': 'c',
		'<': 'c',
		'3': 'e',
		'9': 'g',
		'#': 'h',
		'1': 'i',
		'0': 'o',
		'5': 's',
		'$': 's',
		'+': 't',
		'7': 'l',
		'2': 'z',

		// Greek letters
		'α': 'a',
		'β': 'b',
		'γ': 'y',
		'∆': 'a',
		'δ': 'd',
		'ε': 'e',
		'ζ': 'z',
		'η': 'n', // 'η': 'hn',
		'θ': 'o',
		'ι': 'i',
		'κ': 'k',
		'λ': 'l',
		'μ': 'u', // 'μ': 'mu',
		'ν': 'v', // 'ν': 'nv',
		'ο': 'o',
		'ρ': 'p',
		'ς': 's',
		'τ': 't',
		'υ': 'u',
		'φ': 'p',
		'χ': 'x',
		'ψ': 't',
		'Ω': 'o',
		'ω': 'w',

		// Math symbols
		'⊗': 'o',
		'⊕': 'o',
		'σ': 'o',
		'∩': 'n',
		'∪': 'u',
		'⊂': 'c',
		'⊆': 'c',
		'⊄': 'c',
		'∈': 'e',
		'⊖': 'o',
		'Ø': 'o',
		'∨': 'v',
		'∄': 'a', // '∄': 'ab',
		'∫': 'l',

		// Letterlike
		'ℂ': 'c',
		'℃': 'c',
		'℄': 'c',
		'ℇ': 'e',
		'℉': 'f',
		'ℊ': 'g',
		'ℋ': 'h',
		'ℌ': 'h',
		'ℍ': 'h',
		'ℎ': 'h',
		'ℏ': 'h',
		'ℐ': 'j',
		'ℑ': 'j',
		'ℒ': 'l',
		'ℓ': 'l',
		'℔': 'b',
		'ℕ': 'n',
		'№': 'n',
		'℗': 'p',
		'℘': 'p',
		'ℙ': 'p',
		'ℚ': 'q',
		'ℛ': 'r',
		'ℜ': 'r',
		'ℝ': 'r',
		'℟': 'r',
		'℣': 'v',
		'ℤ': 'z',
		'℧': 'o',
		'℩': 'i',
		'K': 'k',
		'Å': 'a',
		'ℬ': 'b',
		'ℭ': 'c',
		'℮': 'e',
		'ℰ': 'e',
		'ℱ': 'f',
		'ℳ': 'm',
		'ℴ': 'o',
		'ℵ': 'n',
		'ℹ': 'i',
		'℺': 'o',
		'ℼ': 'n',
		'ℽ': 'v',
		'ℿ': 'n',
		'⅀': 'e',
		'⅁': 'g',
		'⅄': 'l',
		'ⅅ': 'd',
		'ⅆ': 'd',
		'ⅇ': 'e',
		'ⅈ': 'i',
		'ⅉ': 'j', // 'ⅉ': 'ji',
		'ⓟ': 'p',
		'ʉ': 'u',
		'ȿ': 's',
		'ⓢ': 's',
		'ⓨ': 'y',
		'ż': 'z',
		'ž': 'z',

		// Confusable: http://www.unicode.org/reports/tr36/confusables.txt
		'е': 'e',
		'о': 'o',
		'ѕ': 's',
		'х': 'x',
		'і': 'i',
		'ј': 'j',
		'р': 'p',
		'с': 'c',
		'у': 'y',
		'ѵ': 'v',
		'ɑ': 'a',
		'ɡ': 'g',
		'ɩ': 'i', // 'ɩ': 'li',
		'ɒ': 'o',
		'г': 'r',
		'π': 'n',
		'ո': 'n',
		'հ': 'h',
		'ս': 'u',
		'ց': 'g',
		'ք': 'p', // 'ք': 'fp',
		'ყ': 'y',
		'୦': 'o',
		'০': 'o',
		'੦': 'o',
		'౦': 'o',
		'೦': 'o',
		'๐': 'o',
		'໐': 'o',
		'᠐': 'o',
		'〇': 'o',
		'օ': 'o',
		'б': 'b',
		'৪': 'b',
		'৭': 'g',
		'੧': 'g',
		'୨': 'g',
	}

	SpecialCharacters = map[rune]rune{
		'-': ' ',
		'_': ' ',
		'|': ' ',
		'.': ' ',
		',': ' ',
		'(': ' ',
		')': ' ',
		'<': ' ',
		'>': ' ',
		'"': ' ',
		'`': ' ',
		'~': ' ',
		'*': ' ',
		'&': ' ',
		'%': ' ',
		'$': ' ',
		'#': ' ',
		'@': ' ',
		'!': ' ',
		'?': ' ',
		'+': ' ',
	}

	WildcardCharacters = map[rune]rune{
		'*': '*',
		'?': '*',
	}
)
