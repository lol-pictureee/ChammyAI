package converter

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

func DecodeBase64Image(base64String string, outputPath string) error {

	imgBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return fmt.Errorf("ошибка декодирования: %w", err)
	}

	err = ioutil.WriteFile(outputPath, imgBytes, 0644)
	if err != nil {
		return fmt.Errorf("ошибка записи файла: %w", err)
	}

	return nil
}
