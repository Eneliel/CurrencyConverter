package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
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
	http.HandleFunc("/Sym", HandlerSym)
	http.HandleFunc("/", HandlerConv)
	http.HandleFunc("/calc", HandlerCalc)
	err := http.ListenAndServe(":8000", nil)
	log.Fatal(err)
}

func HandlerSym(w http.ResponseWriter, r *http.Request) {
	html, err := template.ParseFiles("Sym.html")
	Err(err)
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
	err = html.Execute(w, s)
	Err(err)
}

func HandlerConv(w http.ResponseWriter, r *http.Request) {
	html, err := template.ParseFiles("Conv.html")
	Err(err)
	err = html.Execute(w, nil)
	Err(err)
}

func HandlerCalc(w http.ResponseWriter, r *http.Request) {
	to := r.FormValue("to")
	from := r.FormValue("from")
	amount := r.FormValue("amount")
	Result := Convert(from, to, amount)
	hmtl, err := template.ParseFiles("Calc.html")
	Err(err)
	err = hmtl.Execute(w, Result)
	Err(err)
}

func Convert(from, to, amount string) float64 {
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
	return s.Result
}

func Err(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
