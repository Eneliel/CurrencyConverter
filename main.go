package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ExchangeRates struct {
	Success bool              `json:"success"`
	Rates   map[string]string `json:"symbols"`
}
type Conv struct {
	Result float64 `json:"result"`
}

func main() {

	var in string
	for {
		fmt.Println("Print any command:\n'Sym', 'Convert','Exit'")
		fmt.Scan(&in)
		if in == "Sym" {
			Sym()
		} else if in == "Convert" {
			var to, from, amount string
			fmt.Println("Convert to >")
			fmt.Scan(&to)
			fmt.Println("Convert from >")
			fmt.Scan(&from)
			fmt.Println("Convert amount >")
			fmt.Scan(&amount)
			Convert(to, from, amount)
		} else if in == "Exit" || in == "exit" {
			break
		}
	}
}
func Sym() {
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
	fmt.Println("Список доступных валют:")
	for i, k := range s.Rates {
		fmt.Println(i, k)
	}
}

func Convert(from, to, amount string) {
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
	var s Conv
	json.Unmarshal(body, &s)
	fmt.Println(s.Result)
}

func Err(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
