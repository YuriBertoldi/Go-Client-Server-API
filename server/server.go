package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
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

type ErroJson struct {
	Mensagem string `json:"mensagem"`
}

func main() {
	log.Println("Servidor iniciado")
	http.HandleFunc("/cotacao", BuscaCotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

func BuscaCotacaoHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	xCotacao, error := BuscaCotacao()
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		xErr := new(ErroJson)
		xErr.Mensagem = error.Error()
		json.NewEncoder(w).Encode(xErr)

		log.Println(error.Error())
		return
	}

	err := GravarCotacao(xCotacao)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		xErr := new(ErroJson)
		xErr.Mensagem = error.Error()
		json.NewEncoder(w).Encode(xErr)

		log.Println(err.Error())
		return
	}

	xCotacaoSimp := new(CotacaoSimplificada)
	xCotacaoSimp.Bid = xCotacao.USDBRL.Bid

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(xCotacaoSimp)
}

func BuscaCotacao() (*Cotacao, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, error := ioutil.ReadAll(res.Body)
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

func GravarCotacao(c *Cotacao) error {
	db, err := sql.Open("sqlite3", "db/cotacoes.db")
	if err != nil {
		return err
	}
	defer db.Close()
	err = insertCotacao(db, c)
	if err != nil {
		return err
	}
	return nil
}

func insertCotacao(db *sql.DB, c *Cotacao) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	stmt, err := db.PrepareContext(ctx, "insert into USDBRL(Code, Bid, Name) values(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, c.USDBRL.Code, c.USDBRL.Bid, c.USDBRL.Name)
	if err != nil {
		return err
	}
	return nil
}
