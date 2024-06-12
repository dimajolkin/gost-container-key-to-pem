package key_reader

import (
	"errors"
	"sort"
	"strings"
)

type KeyContainerValidator struct {
	fileMap map[string]bool
}

func NewKeyContainerValidator() *KeyContainerValidator {
	return &KeyContainerValidator{
		fileMap: map[string]bool{
			"header.key":   false,
			"masks.key":    false,
			"masks2.key":   false,
			"name.key":     false,
			"primary.key":  false,
			"primary2.key": false,
		},
	}
}

func (v *KeyContainerValidator) Add(fileName string) {
	if _, ok := v.fileMap[fileName]; ok {
		v.fileMap[fileName] = true
	}
}

func (v *KeyContainerValidator) GetError() error {
	errorMessage := make([]string, 0, 5)
	for name, value := range v.fileMap {
		if value == false {
			errorMessage = append(errorMessage, name)
		}
	}

	sort.Strings(errorMessage)

	if len(errorMessage) > 0 {
		return errors.New(strings.Join(errorMessage, ", ") + " файлы не найдены")
	}

	return nil
}
