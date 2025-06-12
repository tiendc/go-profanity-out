package profanityout

import "strconv"

var (
	htmlEntityMap = map[string]rune{
		"quot":   '"',
		"apos":   '\'',
		"amp":    '&',
		"lt":     '<',
		"gt":     '>',
		"nbsp":   ' ',
		"iexcl":  '¡',
		"cent":   '¢',
		"pound":  '£',
		"curren": '¤',
		"yen":    '¥',
		"brvbar": '¦',
		"sect":   '§',
		"uml":    '¨',
		"copy":   '©',
		"ordf":   'ª',
		"laquo":  '«',
		"not":    '¬',
		"shy":    '­',
		"reg":    '®',
		"macr":   '¯',
		"deg":    '°',
		"plusmn": '±',
		"sup2":   '²',
		"sup3":   '³',
		"acute":  '´',
		"micro":  'µ',
		"para":   '¶',
		"middot": '·',
		"cedil":  '¸',
		"sup1":   '¹',
		"ordm":   'º',
		"raquo":  '»',
		"frac14": '¼',
		"frac12": '½',
		"frac34": '¾',
		"iquest": '¿',
		"times":  '×',
		"divide": '÷',
		"Agrave": 'À',
		"Aacute": 'Á',
		"Acirc":  'Â',
		"Atilde": 'Ã',
		"Auml":   'Ä',
		"Aring":  'Å',
		"AElig":  'Æ',
		"Ccedil": 'Ç',
		"Egrave": 'È',
		"Eacute": 'É',
		"Ecirc":  'Ê',
		"Euml":   'Ë',
		"Igrave": 'Ì',
		"Iacute": 'Í',
		"Icirc":  'Î',
		"Iuml":   'Ï',
		"ETH":    'Ð',
		"Ntilde": 'Ñ',
		"Ograve": 'Ò',
		"Oacute": 'Ó',
		"Ocirc":  'Ô',
		"Otilde": 'Õ',
		"Ouml":   'Ö',
		"Oslash": 'Ø',
		"Ugrave": 'Ù',
		"Uacute": 'Ú',
		"Ucirc":  'Û',
		"Uuml":   'Ü',
		"Yacute": 'Ý',
		"THORN":  'Þ',
		"szlig":  'ß',
		"agrave": 'à',
		"aacute": 'á',
		"acirc":  'â',
		"atilde": 'ã',
		"auml":   'ä',
		"aring":  'å',
		"aelig":  'æ',
		"ccedil": 'ç',
		"egrave": 'è',
		"eacute": 'é',
		"ecirc":  'ê',
		"euml":   'ë',
		"igrave": 'ì',
		"iacute": 'í',
		"icirc":  'î',
		"iuml":   'ï',
		"eth":    'ð',
		"ntilde": 'ñ',
		"ograve": 'ò',
		"oacute": 'ó',
		"ocirc":  'ô',
		"otilde": 'õ',
		"ouml":   'ö',
		"oslash": 'ø',
		"ugrave": 'ù',
		"uacute": 'ú',
		"ucirc":  'û',
		"uuml":   'ü',
		"yacute": 'ý',
		"thorn":  'þ',
		"yuml":   'ÿ',
	}
)

func skipHTMLTag(content []rune, i int) (next int) {
	length := len(content)
	j := i + 1
	for j < length {
		// TODO: checks that no '<' appears again in between a tag
		if content[j] == '>' {
			return j + 1
		}
		j++
	}
	return i
}

func decodeHTMLEntityAt(content []rune, i int) (rune, int) {
	length := len(content)
	j := i + 1
	found := false
	for j < length {
		if content[j] == ';' {
			found = true
			break
		}
		j++
	}
	if !found || j-i > 10 { //nolint:mnd
		return 0, i
	}
	name := string(content[i+1 : j])
	j++
	if name[0] == '#' {
		code, err := strconv.ParseInt(name[1:], 10, 32)
		if err != nil {
			return 0, i
		}
		return rune(code), j
	}
	if ch, exists := htmlEntityMap[name]; exists {
		return ch, j
	}
	return 0, i
}
