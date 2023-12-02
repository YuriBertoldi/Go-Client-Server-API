package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.println(err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.println(err.Error())
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.println(err)
	}
	var c Cotacao
	err = json.Unmarshal(body, &c)
	if err != nil {
		log.println(err)
	}

	gravarArquivo(c.Bid)
}

func gravarArquivo(dolar string) {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		log.println(err.Error())
	}
	defer f.Close()

	tamanho, err := f.WriteString("Dólar: " + dolar)
	if err != nil {
		log.println(err.Error())
	}
	log.println("Arquivo criado com sucesso! Tamanho: %d bytes\n", tamanho)

}
