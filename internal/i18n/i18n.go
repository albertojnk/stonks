package i18n

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"sync"

	"github.com/albertojnk/stonks/internal/context"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	bundleOnce         sync.Once
	bundleInstance     *i18n.Bundle
	pageEntry          map[string][]string
	supportedLanguages []string
)

// NewBundle create a singleton instance of *i18n.Bundle, who know where to look and find the translation files
func NewBundle() *i18n.Bundle {
	bundleOnce.Do(func() {
		bundleInstance = i18n.NewBundle(language.AmericanEnglish)
		bundleInstance.RegisterUnmarshalFunc("json", json.Unmarshal)

		files, err := ioutil.ReadDir(path.Join(".", "translations"))
		if err != nil {
			log.Fatal(err)
		}

		// Using fixed file to get all entries
		fileBytes, _ := ioutil.ReadFile(path.Join(".", "translations", "en-US.json"))
		data := make(map[string]map[string]interface{})
		json.Unmarshal(fileBytes, &data)

		pageEntry = make(map[string][]string)
		for keyHighLevel, valueHighLevel := range data {
			for keyLowLevel := range valueHighLevel {
				pageEntry[keyHighLevel] = append(pageEntry[keyHighLevel], keyLowLevel)
			}
		}

		for _, file := range files {
			_, err = bundleInstance.LoadMessageFile(path.Join(".", "translations", file.Name()))
			if err != nil {
				fmt.Printf("error when loading translation file: %v", file.Name())
				continue
			}

			supportedLanguages = append(supportedLanguages, file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))])
		}
	})

	return bundleInstance
}

// NewLocalizer create a new i18n localizer
func NewLocalizer(lang, accept string) *i18n.Localizer {
	return i18n.NewLocalizer(bundleInstance, lang, accept)
}

// GetSupportedLang is a helper to check if language passed is supported by system
func GetSupportedLang(lang string) string {
	for _, supportedLang := range supportedLanguages {
		if lang == supportedLang {
			return lang
		}
	}
	return "en-US"
}

// GetSupportedLanguages get all system suported laguages codes
func GetSupportedLanguages() []string {
	return supportedLanguages
}

// SetLang change the localizer context language
func SetLang(c *context.Context, lang string) *context.Context {
	NewBundle()
	c.Localizer = NewLocalizer(lang, lang)
	return c
}

// Translate is a shortcut for `c.Localizer.Localize(&i18n.LocalizeConfig{MessageID: id})` and also treat errors with default message
func Translate(c *context.Context, id string) string {
	defaultMessage := c.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "TranslationNotFound"})
	message, err := c.Localizer.Localize(&i18n.LocalizeConfig{MessageID: id})
	if err != nil {
		c.Logger.Errorf("Error when searching a translation: ", err.Error())
		return defaultMessage
	}
	return message
}

// TranslateNestedItem is a shortcut for `Translate(c, fmt.Sprintf("%s.%s", nestedKey, id))`
func TranslateNestedItem(c *context.Context, nestedKey, id string) string {
	return Translate(c, fmt.Sprintf("%s.%s", nestedKey, id))
}

// TranslatePageFull retrive all nested traslation based on nestedKey
func TranslatePageFull(c *context.Context, nestedKey string) (map[string]string, error) {
	pageTranslation := make(map[string]string)
	defaultMessage := c.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "TranslationNotFound"})

	if _, ok := pageEntry[nestedKey]; !ok {
		return pageTranslation, fmt.Errorf(defaultMessage)
	}

	for _, entryKey := range pageEntry[nestedKey] {
		translation := TranslateNestedItem(c, nestedKey, entryKey)
		pageTranslation[entryKey] = translation
	}

	return pageTranslation, nil
}

func SetBody(c *context.Context, body gin.H, translations map[string]string) gin.H {
	for k, v := range translations {
		if _, ok := body[k]; !ok {
			body[k] = v
		} else {
			c.Logger.Warn("trying to overwrite body key: %v", k)
		}
	}
	return body
}
