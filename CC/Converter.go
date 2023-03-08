package cur_conv

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type ExchangeRates struct {
	Success bool              `json:"success"`
	Rates   map[string]string `json:"symbols"`
}

type Convert struct {
	Result float64 `json:"result"`
}

func Sym() string {
	url := "https://api.apilayer.com/exchangerates_data/symbols"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", "S5zUr9PkYd4GbHaoQwGh5xtbf70YTMGB")
	Err(err)
	res, err := client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	}
	Err(err)
	body, _ := ioutil.ReadAll(res.Body)
	var s ExchangeRates
	err = json.Unmarshal(body, &s)
	Err(err)
	var out string
	for k, v := range s.Rates {
		out += k + " : " + v + "\n"
	}
	return out
}

func Conv(from, to, amount string) string {
	url := "https://api.apilayer.com/exchangerates_data/convert?to=" + to + "&from=" + from + "&amount=" + amount
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", "S5zUr9PkYd4GbHaoQwGh5xtbf70YTMGB")

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	}
	Err(err)
	body, err := ioutil.ReadAll(res.Body)
	Err(err)
	var s Convert
	err = json.Unmarshal(body, &s)
	Err(err)
	out := strconv.FormatFloat(s.Result, 'f', -1, 64)
	return string(out)
}

func Err(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
