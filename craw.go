package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	//pacote nao padrao que oferece um parser HTML
	"golang.org/x/net/html"
)

var w = flag.String("w", "", "Set Wordlist")

func main() {
	println("")
	println("")
	println("")
	println("")
	println(" 		    •• Web-CrawlF •• \n")
	println(" 	   	••• Srt.Lilith_64 ••• \n")
	println(" 	•••• https://github.com/SrtLilith64 ••••\n")
	println("")
	println("	 		( ⌐■  __ ■ ) ")
	println("								  ")
	println("")
	println("")
	println("")

	u := flag.String("u", "", "URL de entrada")           //cria uma variavel do tipo string, que define a url de input do programa, e inclue uma descriçao ao comando "u"
	d := flag.Int("d", 1, "Nível de profundidade BFS")    //cria uma variavel do tipo inteiro, que define o nivel de profundidade do BFS, e inclue uma descrição ao comando "d"
	o := flag.String("o", "", "Cria um arquivo de saída") //cria uma variavel do tipo string que armazena o nome de um arquivo que vai salvar todos os resultados finais que irão ser gerados no fim da execucao do programa

	flag.Parse()                          //processa o input do usuario
	var urlsLidas, urlsParaLer []string   //cria uma variavel do tipo string para armazenar e separar as urls que ja foram lidas e as que faltam ainda ler
	urlsParaLer = append(urlsParaLer, *u) //le a url de input armazenada no valor u
	for i := 0; i < *d; i++ {             //cria um loop de repetiçao que procura e armazena as urls encontradas na pagina
		urlsLidas = append(urlsLidas, urlsParaLer...)  //subistitui as urls que ja foram lidas para as novas que ainda faltam ler
		urlsParaLer = breadthFirst(crawl, urlsParaLer) //armazena as urls que ja foram lidas das novas que tem para ler
	}
	urlsLidas = append(urlsLidas, urlsParaLer...)
	for _, url := range urlsLidas { //cria um loop que salva todas as urls encontradas
		fmt.Println(url) //printa as urls encontradas
	}

	writeOutput(urlsLidas, *o) //chama a funcao writeOutput para criar um arquivo de saida
}

func fuzzScan(worklist string) {
	readFile, err := os.Open(*w)
	if err != nil {
		log.Fatalf("Falha ao ler o arquivo : %s", err)
	}

	fileScan := bufio.NewScanner(readFile)
	fileScan.Split(bufio.ScanLines)
	var fileTextLines []string

	for fileScan.Scan() {
		fileTextLines = append(fileTextLines, fileScan.Text())
	}

	for _, eachline := range fileTextLines {
		http.Get(worklist + "/" + eachline)
		if err != nil {
			fmt.Println(" A URL Deve ser com 'HTTP://' ou 'HTTPS://'")
			log.Fatalln(err)
		}

	}

	readFile.Close()

}

//funçao que cria o arquivo de saida com o resultado da variavel urlsLidas
func writeOutput(urlsLidas []string, o string) {
	f, err := os.Create(o)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	write := bufio.NewWriter(f)
	for _, url := range urlsLidas {
		write.WriteString(url + "\n")
	}
	write.Flush()
}

func breadthFirst(f func(string) []string, worklist []string) []string { //cria uma funçao f do tipo funçao com entrada e saida do tipo string e grava o valor do "urlsLidas" na funçao worklist do tipo string
	seen := make(map[string]bool) //cria um map que armazena um conjunto de chave/valor do tipo string boleano
	for len(worklist) > 0 {       //cria um loop que exibe as urls ja acessadas
		items := worklist            //copia as urls do worklist para o itens
		worklist = nil               // apaga o valor gravado do worklist
		for _, item := range items { // cria um par de valores
			// verifica se a url ja foi acessada
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...) // regrava o valor atual da funcao f na worklist com as novas urls
			}
		}

	}
	return worklist //retorna o valor armazenado na variavel worklist
}

func crawl(url string) []string { // cria uma entrada(url) e uma saida do tipo string
	fmt.Println(url) //printa o valor que foi gravado na url
	list, err := Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list //retorna uma lista
}

//trata o percurso
//recursao
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}

}

/*
Faz uma requisição HTTP GET ao o URL especificado fazendo a analise do HTML
e devolve os links encontrados no HTML do site especificado
*/
func Extract(url string) ([]string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	//retorna o status code da requisição se der algum erro inesperado (403,404,500 e etc...)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("recebendo %s: %s", url, resp.Status)
	}
	// a função "html.Parse" le uma sequencia de bytes e faz o seu parse e devolve a raiz do documento HTML
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Analisando %s as HTML: %v", url, err)
	}
	var links []string                // declara uma variavel links com saida do tipo string
	visitNode := func(n *html.Node) { //cria uma variavel de funcao
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue //ignora as urls que nao sao compativeis
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}
