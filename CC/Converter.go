package cur_conv

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type ExchangeRates struct {
	Success bool              `json:"success"`
	Rates   map[string]string `json:"symbols"`
}

type Convert struct {
	Result float64 `json:"result"`
}

type User struct {
	User_id int64
	MyMoney float64
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

func Wallet(in, do string, id int64) string {
	money, err := strconv.ParseFloat(in, 64)
	Err(err)
	db, err := sql.Open("mysql", "root:@/Wallet_tg")
	if err != nil {
		log.Print("Ошибка с ДБ", err)
	}
	defer db.Close()
	sel := fmt.Sprintf("SELECT * FROM `User` WHERE id =%d", id)
	var u User
	err = db.QueryRow(sel).Scan(&u.User_id, &u.MyMoney)
	if err == sql.ErrNoRows {
		insert := fmt.Sprintf("INSERT INTO `User`(`id`, `Money`) VALUES ('%d','%f')", id, money)
		ins_db, err := db.Query(insert)
		if err != nil {
			log.Fatal(err)
		}
		u = User{
			User_id: id,
			MyMoney: money,
		}
		defer ins_db.Close()
	} else if err != nil {
		log.Fatal(err)
	} else {

		sel_db, err := db.Query(sel)
		if err != nil {
			log.Fatal(err)
		}
		for sel_db.Next() {
			err = sel_db.Scan(&u.User_id, &u.MyMoney)
			if err != nil {
				panic(err)
			}
		}
		if do == "sum" {
			u.MyMoney += money
		} else {
			u.MyMoney -= money
		}
		update_db := fmt.Sprintf("UPDATE `User` SET `Money`='%f' WHERE id =%d", u.MyMoney, id)
		upd, err := db.Query(update_db)
		if err != nil {
			log.Fatal(err)
		}
		defer upd.Close()

	}
	return strconv.FormatFloat(u.MyMoney, 'f', 5, 64)
}

func Err(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
