package nicefish

import (
	"fmt"
	"net/http"
	"os/exec"
)

func Cookies(twkUuid string) []*http.Cookie {
	return []*http.Cookie{
		{Name: "user_id", Value: "7d40ed97-f964-4b49-a822-4656a57c3a6d"},
		{Name: "NUXT_LOCALE", Value: "ru"},
		{Name: "TawkConnectionTime", Value: "0"},
		{Name: "twk_uuid_66f0cb3de5982d6c7bb2f3cb", Value: twkUuid}, // setCookieSelenium
	}
}

func setCookieSelenium() {
	cmd := exec.Command("python3", "selenium_cookie.py")
	out, errGetCookieSelenium := cmd.Output()
	if errGetCookieSelenium != nil {
		fmt.Println(errGetCookieSelenium)
		return
	}

	fmt.Println(string(out))

}
