package viesparser

import (
	"regexp"
	"slices"
	"strings"
)

type ParsedAddress struct {
	Street string
	City   string
	Zip    string
}

type Config struct {
	IgnoreGreek bool
}

const (
	cz = "CZ"
	sk = "SK"
	nl = "NL"
	be = "BE"
	fr = "FR"
	pt = "PT"
	it = "IT"
	fi = "FI"
	ro = "RO"
	si = "SI"
	at = "AT"
	pl = "PL"
	hr = "HR"
	el = "EL"
	dk = "DK"
	ee = "EE"
)

var (
	SupportedCountryCodes = []string{
		cz, sk, nl, be, fr, pt, it, fi, ro, si, at, pl, hr, el, dk, ee,
	}
	greekExpressions = map[*regexp.Regexp]string{
		regexp.MustCompile("/[αΑ][ιίΙΊ]/u"):                             "e",
		regexp.MustCompile("/[οΟΕε][ιίΙΊ]/u"):                           "i",
		regexp.MustCompile(`/[αΑ][υύΥΎ]([θΘκΚξΞπΠσςΣτTφΡχΧψΨ]|\s|$)/u`): "af$1",
		regexp.MustCompile("/[αΑ][υύΥΎ]/u"):                             "av",
		regexp.MustCompile(`/[εΕ][υύΥΎ]([θΘκΚξΞπΠσςΣτTφΡχΧψΨ]|\s|$)/u`): "ef$1",
		regexp.MustCompile("/[εΕ][υύΥΎ]/u"):                             "ev",
		regexp.MustCompile("/[οΟ][υύΥΎ]/u"):                             "ou",
		regexp.MustCompile(`/(^|\s)[μΜ][πΠ]/u`):                         "$1b",
		regexp.MustCompile(`/[μΜ][πΠ](\s|$)/u`):                         "b$1",
		regexp.MustCompile(`/[μΜ][πΠ]/u`):                               "mp",
		regexp.MustCompile(`/[νΝ][τΤ]/u`):                               "nt",
		regexp.MustCompile(`/[τΤ][σΣ]/u`):                               "ts",
		regexp.MustCompile(`/[τΤ][ζΖ]/u`):                               "tz",
		regexp.MustCompile(`/[γΓ][γΓ]/u`):                               "ng",
		regexp.MustCompile(`/[γΓ][κΚ]/u`):                               "gk",
		regexp.MustCompile(`/[ηΗ][υΥ]([θΘκΚξΞπΠσςΣτTφΡχΧψΨ]|\s|$)/u`):   "if$1",
		regexp.MustCompile(`/[ηΗ][υΥ]/u`):                               "iu",
		regexp.MustCompile(`/[θΘ]/u`):                                   "th",
		regexp.MustCompile(`/[χΧ]/u`):                                   "ch",
		regexp.MustCompile(`/[ψΨ]/u`):                                   "ps",
		regexp.MustCompile(`/[αά]/u`):                                   "a",
		regexp.MustCompile(`/[βΒ]/u`):                                   "v",
		regexp.MustCompile(`/[γΓ]/u`):                                   "g",
		regexp.MustCompile(`/[δΔ]/u`):                                   "d",
		regexp.MustCompile(`/[εέΕΈ]/u`):                                 "e",
		regexp.MustCompile(`/[ζΖ]/u`):                                   "z",
		regexp.MustCompile(`/[ηήΗΉ]/u`):                                 "i",
		regexp.MustCompile(`/[ιίϊΙΊΪ]/u`):                               "i",
		regexp.MustCompile(`/[κΚ]/u`):                                   "k",
		regexp.MustCompile(`/[λΛ]/u`):                                   "l",
		regexp.MustCompile(`/[μΜ]/u`):                                   "m",
		regexp.MustCompile(`/[νΝ]/u`):                                   "n",
		regexp.MustCompile(`/[ξΞ]/u`):                                   "x",
		regexp.MustCompile(`/[οόΟΌ]/u`):                                 "o",
		regexp.MustCompile(`/[πΠ]/u`):                                   "p",
		regexp.MustCompile(`/[ρΡ]/u`):                                   "r",
		regexp.MustCompile(`/[σςΣ]/u`):                                  "s",
		regexp.MustCompile(`/[τΤ]/u`):                                   "t",
		regexp.MustCompile(`/[υύϋΥΎΫ]/u`):                               "i",
		regexp.MustCompile(`/[φΦ]/iu`):                                  "f",
		regexp.MustCompile(`/[ωώ]/iu`):                                  "o",
		regexp.MustCompile(`/[Α]/iu`):                                   "a",
	}
)

func ParseAddress(countryCode, address string, config ...Config) (ParsedAddress, error) {
	if len(countryCode) == 0 {
		return ParsedAddress{}, ErrorMissingCountryCode
	}
	if len(address) == 0 {
		return ParsedAddress{}, ErrorMissingAddress
	}
	countryCode = strings.TrimSpace(countryCode)
	address = strings.TrimSpace(address)
	newlinesCount := strings.Count(address, "\n")
	if !slices.Contains(SupportedCountryCodes, countryCode) {
		return ParsedAddress{}, ErrorUnsupportedCountryCode
	}
	if newlinesCount == 1 && slices.Contains([]string{nl, be, fr, fi, at, pl, dk}, countryCode) {
		parts := strings.Split(address, "\n")
		locationParts := strings.Split(parts[1], " ")
		return ParsedAddress{
			Street: strings.TrimSpace(parts[0]),
			Zip:    strings.TrimSpace(locationParts[0]),
			City:   strings.TrimSpace(locationParts[1]),
		}, nil
	}
	if newlinesCount == 0 && slices.Contains([]string{si, hr}, countryCode) {
		parts := strings.Split(address, ",")
		street := strings.TrimSpace(parts[0])
		if len(parts) == 3 {
			street = street + ", " + strings.TrimSpace(parts[1])
		}
		locationParts := strings.Split(parts[len(parts)-1], " ")
		return ParsedAddress{
			Street: street,
			Zip:    locationParts[0],
			City:   locationParts[1],
		}, nil
	}
	if countryCode == sk {
		if newlinesCount == 1 {
			var city, zip string
			parts := strings.Split(address, "\n")
			street := strings.TrimSpace(parts[0])
			if parts[1] != "Slovensko" {
				locationParts := strings.Split(parts[len(parts)-1], " ")
				zip = locationParts[0]
				city = locationParts[1]
			}
			if parts[1] == "Slovensko" {
				locationParts := strings.Split(parts[0], " ")
				zip = locationParts[0]
				city = locationParts[1]
				street = ""
			}
			city = strings.Replace(city, "mestská časť ", "", 1)
			city = strings.Replace(city, "m. č. ", "", 1)
			return ParsedAddress{
				Street: strings.TrimSpace(street),
				City:   strings.TrimSpace(city),
				Zip:    strings.TrimSpace(zip),
			}, nil
		}
		if newlinesCount == 2 {
			var city, zip string
			parts := strings.Split(address, "\n")
			street := strings.TrimSpace(parts[0])
			locationParts := strings.Split(parts[len(parts)-1], " ")
			zip = locationParts[0]
			city = locationParts[1]
			city = strings.Replace(city, "mestská časť ", "", 1)
			city = strings.Replace(city, "m. č. ", "", 1)
			return ParsedAddress{
				Street: strings.TrimSpace(street),
				City:   strings.TrimSpace(city),
				Zip:    strings.TrimSpace(zip),
			}, nil
		}
	}
	if countryCode == cz {
		if newlinesCount == 1 {
			parts := strings.Split(address, "\n")
			street := strings.TrimSpace(parts[0])
			lastParts := strings.Split(strings.TrimSpace(parts[len(parts)-1]), " ")
			return ParsedAddress{
				Street: street,
				City:   strings.TrimSpace(strings.Join(lastParts[len(lastParts)-2:len(lastParts)-1], "")),
				Zip:    strings.TrimSpace(strings.Join(lastParts[:len(lastParts)-2], "")),
			}, nil
		}
		if newlinesCount == 2 {
			parts := strings.Split(address, "\n")
			lastParts := strings.Split(strings.TrimSpace(parts[len(parts)-1]), " ")
			return ParsedAddress{
				Street: strings.TrimSpace(parts[0]),
				City:   strings.TrimSpace(parts[1]),
				Zip:    strings.TrimSpace(strings.Join(lastParts[:len(lastParts)-2], "")),
			}, nil
		}
		return ParsedAddress{}, ErrorInvalidOption
	}
	return ParsedAddress{}, nil
}

func MustParseAddress(countryCode, address string) ParsedAddress {
	parsed, err := ParseAddress(countryCode, address)
	if err != nil {
		panic(err)
	}
	return parsed
}
