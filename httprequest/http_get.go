package httprequest

import "net/http"

func HttpGet(url string, headers map[string]string) (*http.Response, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Сет заголовоков
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
