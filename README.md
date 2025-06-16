[![Go Version][gover-img]][gover] [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov] [![GoReport][rpt-img]][rpt]

# Profanity detection library

This project is inspired by [github.com/TwiN/go-away](https://github.com/TwiN/go-away) and
[github.com/finnbear/moderation](https://github.com/finnbear/moderation). It also uses
language data from the libs (with modifications).

This project is still in development and more tests are needed to ensure the accuracy.
However, you may use it in your work as it can produce good results.

## Highlights

- Fully supports Unicode
- Utilizes radix tree to improve performance

## Installation

```shell
go get github.com/tiendc/go-profanity-out
```

## How to use

```go
import (
    profanityout "github.com/tiendc/go-profanity-out"
    profanityDataEN "github.com/tiendc/go-profanity-out/data/en"
)

detector := profanityout.NewProfanityDetector().
    WithProfaneWords(profanityDataEN.DefaultProfanities).          // required
    WithFalsePositiveWords(profanityDataEN.DefaultFalsePositives). // required
    WithSuspectWords(profanityDataEN.DefaultSuspects).             // required
    WithLeetSpeakCharacters(profanityDataEN.LeetSpeakCharacters).  // required
    WithSpecialCharacters(profanityDataEN.SpecialCharacters).      // required
    WithWildcardCharacters(profanityDataEN.WildcardCharacters).    // required
    WithSanitizeLeetSpeak(true).                                   // default: true
    WithSanitizeSpecialCharacters(true).                           // default: true
    WithSanitizeSpaces(true).                                      // default: true
    WithMatchWholeWord(true).                                      // default: true
    WithSanitizeRepeatedCharacters(true).                          // default: true
    WithSanitizeWildcardCharacters(false).                         // default: false
    WithSanitizeAccents(true).                                     // default: true
    WithProcessInputAsHTML(false).                                 // default: false
    WithConfidenceCalculator(calculator).                          // default: built-in
    WithCensorCharacter('*')                                       // default: *

// Scan for at most one profanity (result may contain found suspect words and/or false positives)
matches := detector.ScanProfanity("fuck this $h!!t") // profane: true

// Scan for all profanities
matches := detector.ScanAllProfanities("fuck this $h!!t") // profane: true

// Censor the profanities
res, matches := detector.Censor("fuck this $h!!t") // res == "**** this *****"

// WithSanitizeLeetSpeak: true
ScanProfanity("$h!t") // profane: true
// WithSanitizeLeetSpeak: false
ScanProfanity("$h!t") // profane: false

// WithSanitizeSpecialCharacters: true
ScanProfanity("sh_it") // profane: true
// WithSanitizeSpecialCharacters: false
ScanProfanity("sh_it") // profane: false

// WithSanitizeSpaces: true
ScanProfanity("f u c k") // profane: true
// WithSanitizeSpaces: false
ScanProfanity("f u c k") // profane: false

// WithSanitizeRepeatedCharacters: true
ScanProfanity("fuuuuck") // profane: true
// WithSanitizeRepeatedCharacters: false
ScanProfanity("fuuuuck") // profane: false

// WithSanitizeWildcardCharacters: true
ScanProfanity("f**k") // profane: true
// WithSanitizeWildcardCharacters: false
ScanProfanity("f**k") // profane: false
// Suppose "f*ck" is in the profanity dictionary
WithProfaneWords([]string{"f*ck"}).ScanProfanity("fxck") // profane: true
// NOTE: With this option you can turn on matching portion of a word for specific words
// without turning off `WithMatchWholeWord(false)` by putting "*word*" in dictionary.

// WithSanitizeAccents: true
ScanProfanity("fúck") // profane: true
// WithSanitizeAccents: false
ScanProfanity("fúck") // profane: false

// WithMatchWholeWord: true
ScanProfanity("fuckyou") // profane: false
// WithMatchWholeWord: false (NOTE: this may reduce the accuracy significantly)
ScanProfanity("fuckyou") // profane: true
// NOTE: consider turning on WithSanitizeWildcardCharacters and putting "*word*" in dictionary
// to scan for non-whole word matching for specific words.

// WithProcessInputAsHTML: true
ScanProfanity("&lt;ock") // profane: true
// WithProcessInputAsHTML: false
ScanProfanity("&lt;ock") // profane: false
```

## Benchmarks

[Benchmark code](https://gist.github.com/tiendc/bd5a0655ad07251f626402d819786d84)

```
tiendc/go-profanity-out
tiendc/go-profanity-out-10         	   10000	    104278 ns/op	   43759 B/op	     300 allocs/op
TwiN/go-away
TwiN/go-away-10                    	    2745	    415685 ns/op	  444899 B/op	     498 allocs/op
finnbear/moderation
finnbear/moderation-10             	   15432	     77601 ns/op	    2496 B/op	      22 allocs/op
```

## Help wanted

- You are welcome to make pull requests for new functions and bug fixes.
- It's really nice if you can add more input data for English and other languages.

## License

- [MIT License](LICENSE)

[doc-img]: https://pkg.go.dev/badge/github.com/tiendc/go-profanity-out
[doc]: https://pkg.go.dev/github.com/tiendc/go-profanity-out
[gover-img]: https://img.shields.io/badge/Go-%3E%3D%201.20-blue
[gover]: https://img.shields.io/badge/Go-%3E%3D%201.20-blue
[ci-img]: https://github.com/tiendc/go-profanity-out/actions/workflows/go.yml/badge.svg
[ci]: https://github.com/tiendc/go-profanity-out/actions/workflows/go.yml
[cov-img]: https://codecov.io/gh/tiendc/go-profanity-out/branch/main/graph/badge.svg
[cov]: https://codecov.io/gh/tiendc/go-profanity-out
[rpt-img]: https://goreportcard.com/badge/github.com/tiendc/go-profanity-out
[rpt]: https://goreportcard.com/report/github.com/tiendc/go-profanity-out
