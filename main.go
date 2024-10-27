package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ncruces/zenity"
)

// Estrutura para mapear os dados da API de geolocalização
type RespostaAPI struct {
	Cidade string `json:"city"`
	Estado string `json:"region"`
	Pais   string `json:"country"`
	Coord  string `json:"loc"`
}

// Estrutura para salvar os dados de localização no JSON
type Localizacao struct {
	Cidade      string `json:"cidade"`
	Coordenadas string `json:"coordenadas"`
	NomePC      string `json:"nome_pc"`
}

// Função para obter o IP público automaticamente
func obterIPPublico() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}

// Função para obter a geolocalização com base no IP
func obterGeolocalizacao(ip, token string) (*Localizacao, error) {
	urlAPI := fmt.Sprintf("https://ipinfo.io/%s?token=%s", ip, token)
	resp, err := http.Get(urlAPI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var resposta RespostaAPI
	if err := json.NewDecoder(resp.Body).Decode(&resposta); err != nil {
		return nil, err
	}

	// Preenche a estrutura Localizacao com os dados da resposta
	localizacao := &Localizacao{
		Cidade:      resposta.Cidade,
		Coordenadas: resposta.Coord,
		NomePC:      "Computador Teste 1",
	}

	return localizacao, nil
}

// Função para salvar a localização em um arquivo JSON
func salvarLocalizacaoJSON(localizacao *Localizacao) error {
	arquivo, err := os.Create("localizacao.json")
	if err != nil {
		return err
	}
	defer arquivo.Close()

	// Codificar os dados como JSON e gravar no arquivo
	encoder := json.NewEncoder(arquivo)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(localizacao); err != nil {
		return err
	}

	return nil
}

func main() {
	// Exibe uma caixa de diálogo para o usuário inserir o token
	token, err := zenity.Entry("Digite seu token IPINFO_TOKEN:", zenity.Title("Token de Autenticação"))
	if err != nil {
		fmt.Println("Erro ao inserir o token:", err)
		return
	}

	// Obter IP público automaticamente
	ip, err := obterIPPublico()
	if err != nil {
		fmt.Println("Erro ao obter IP público:", err)
		return
	}

	// Obter informações de geolocalização
	localizacao, err := obterGeolocalizacao(ip, token)
	if err != nil {
		fmt.Println("Erro ao obter geolocalização:", err)
		return
	}

	// Salvar a localização em um arquivo JSON
	if err := salvarLocalizacaoJSON(localizacao); err != nil {
		fmt.Println("Erro ao salvar dados no JSON:", err)
		return
	}

	fmt.Println("Informações de localização salvas no arquivo 'localizacao.json'.")

	// Servir arquivos da pasta "web" como conteúdo estático
	http.Handle("/", http.FileServer(http.Dir("./web")))

	// Endpoint para servir o arquivo JSON de localização
	http.HandleFunc("/dados", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "localizacao.json")
	})

	// Iniciar o servidor em localhost:8080
	fmt.Println("Servidor rodando em http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}
