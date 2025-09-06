package discord

// Locale represents the accepted languages for Discord.
// https://discord.com/developers/docs/reference#locales
type Locale string

func (l Locale) String() string {
	if name, ok := Locales[l]; ok {
		return name
	}
	return LocaleUnknown.String()
}

const (
	LocaleEnglishUS    Locale = "en-US"
	LocaleEnglishGB    Locale = "en-GB"
	LocaleBulgarian    Locale = "bg"
	LocaleChineseCN    Locale = "zh-CN"
	LocaleChineseTW    Locale = "zh-TW"
	LocaleCroatian     Locale = "hr"
	LocaleCzech        Locale = "cs"
	LocaleDanish       Locale = "da"
	LocaleDutch        Locale = "nl"
	LocaleFinnish      Locale = "fi"
	LocaleFrench       Locale = "fr"
	LocaleGerman       Locale = "de"
	LocaleGreek        Locale = "el"
	LocaleHindi        Locale = "hi"
	LocaleHungarian    Locale = "hu"
	LocaleItalian      Locale = "it"
	LocaleJapanese     Locale = "ja"
	LocaleKorean       Locale = "ko"
	LocaleLithuanian   Locale = "lt"
	LocaleNorwegian    Locale = "no"
	LocalePolish       Locale = "pl"
	LocalePortugueseBR Locale = "pt-BR"
	LocaleRomanian     Locale = "ro"
	LocaleRussian      Locale = "ru"
	LocaleSpanishES    Locale = "es-ES"
	LocaleSpanishLATAM Locale = "es-419"
	LocaleSwedish      Locale = "sv-SE"
	LocaleThai         Locale = "th"
	LocaleTurkish      Locale = "tr"
	LocaleUkrainian    Locale = "uk"
	LocaleVietnamese   Locale = "vi"
	LocaleUnknown      Locale = ""
)

// Locales is a map of all the languages codes to their names.
var Locales = map[Locale]string{
	LocaleEnglishUS:    "English (United States)",
	LocaleEnglishGB:    "English (Great Britain)",
	LocaleBulgarian:    "LocaleBulgarian",
	LocaleChineseCN:    "Chinese (China)",
	LocaleChineseTW:    "Chinese (Taiwan)",
	LocaleCroatian:     "LocaleCroatian",
	LocaleCzech:        "LocaleCzech",
	LocaleDanish:       "LocaleDanish",
	LocaleDutch:        "LocaleDutch",
	LocaleFinnish:      "LocaleFinnish",
	LocaleFrench:       "LocaleFrench",
	LocaleGerman:       "LocaleGerman",
	LocaleGreek:        "LocaleGreek",
	LocaleHindi:        "LocaleHindi",
	LocaleHungarian:    "LocaleHungarian",
	LocaleItalian:      "LocaleItalian",
	LocaleJapanese:     "LocaleJapanese",
	LocaleKorean:       "LocaleKorean",
	LocaleLithuanian:   "LocaleLithuanian",
	LocaleNorwegian:    "LocaleNorwegian",
	LocalePolish:       "LocalePolish",
	LocalePortugueseBR: "Portuguese (Brazil)",
	LocaleRomanian:     "LocaleRomanian",
	LocaleRussian:      "LocaleRussian",
	LocaleSpanishES:    "Spanish (Spain)",
	LocaleSpanishLATAM: "Spanish (LATAM)",
	LocaleSwedish:      "LocaleSwedish",
	LocaleThai:         "LocaleThai",
	LocaleTurkish:      "LocaleTurkish",
	LocaleUkrainian:    "LocaleUkrainian",
	LocaleVietnamese:   "LocaleVietnamese",
	LocaleUnknown:      "unknown",
}
