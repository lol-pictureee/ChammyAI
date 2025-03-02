package converter

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

//func ImageFileToBase64(imagePath string) string {
//	fmt.Println("Converting")
//	return "base64_encoded_image_data"
//}

func ImageUrlToBase64(imageUrl string) string {
	fmt.Println("Converting...")
	response, err := http.Get(imageUrl)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		fmt.Println("Status code is: ", response.StatusCode)
		return ""
	}
	imageData, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}
	encoded := base64.StdEncoding.EncodeToString(imageData)
	return encoded
}
