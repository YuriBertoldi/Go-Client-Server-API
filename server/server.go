package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Cotacao struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type CotacaoSimplificada struct {
	Bid string `json:"bid"`
}

func main() {
	http.HandleFunc("/cotacao", BuscaCotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

func BuscaCotacaoHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ctx := r.Context()

	select {
	case <-time.After(5 * time.Second):
		println("teste")
	case <-ctx.Done():
		http.Error(w, "Request cancelada pelo cliente", http.StatusRequestTimeout)
	}

	xCotacao, error := BuscaCotacao()
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	GravarCotacao(xCotacao)

	xCotacaoSimp := new(CotacaoSimplificada)
	xCotacaoSimp.Bid = xCotacao.USDBRL.Bid

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(xCotacaoSimp)
}

func BuscaCotacao() (*Cotacao, error) {
	resp, error := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if error != nil {
		return nil, error
	}
	defer resp.Body.Close()
	body, error := ioutil.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}
	var c Cotacao
	error = json.Unmarshal(body, &c)
	if error != nil {
		return nil, error
	}
	return &c, nil
}

func GravarCotacao(c *Cotacao) {
	db, err := sql.Open("sqlite3", "./db/cotacoes.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = insertCotacao(db, c)
	if err != nil {
		panic(err)
	}

}

func insertCotacao(db *sql.DB, c *Cotacao) error {
	stmt, err := db.Prepare("insert into USDBRL(Code, Bid, Name) values(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.USDBRL.Code, c.USDBRL.Bid, c.USDBRL.Name)
	if err != nil {
		return err
	}
	return nil
}
