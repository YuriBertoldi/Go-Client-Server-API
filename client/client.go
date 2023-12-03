package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

type ErroJson struct {
	Mensagem string `json:"mensagem"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Println(err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err.Error())
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != http.StatusOK {
		var mensagem ErroJson
		err = json.Unmarshal(body, &mensagem)
		if err != nil {
			panic(err)
		}
		panic(mensagem.Mensagem)
	}

	var c Cotacao
	err = json.Unmarshal(body, &c)
	if err != nil {
		panic(err)
	}

	gravarArquivo(c.Bid)
}

func gravarArquivo(dolar string) {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		log.Println(err.Error())
	}
	defer f.Close()

	f.WriteString("Dólar: " + dolar)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("Arquivo criado com sucesso!")
	}

}
