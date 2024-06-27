package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	//pacote nao padrao que oferece um parser HTML
	"golang.org/x/net/html"
)

var w = flag.String("w", "", "Set Wordlist (Nao disponivel)")
var o = flag.String("o", "", "Cria um arquivo de saída") //cria uma variavel do tipo string que armazena o nome de um arquivo que vai salvar todos os resultados finais que irão ser gerados no fim da execucao do programa
var p = flag.Bool("p", false, "Exibe a personalidade do Specter")
var wurl string

func main() {
	fmt.Println("")
	fmt.Println("      ••CrawlSpectrum••")
	fmt.Println("       •••Lilithh3x•••")
	fmt.Println("")
	fmt.Println("            /\\_/\\")
	fmt.Println("           ( o.o )	<-- Specter")
	fmt.Println("            > ^ <")
	fmt.Println("           /     \\")
	fmt.Println("          |       |")
	fmt.Println("         /| |   | |\\")
	fmt.Println("")
	fmt.Println("")

	u := flag.String("u", "", "URL de entrada")        //cria uma variavel do tipo string, que define a url de input do programa, e inclue uma descriçao ao comando "u"
	d := flag.Int("b", 1, "Nível de profundidade BFS") //cria uma variavel do tipo inteiro, que define o nivel de profundidade do BFS, e inclue uma descrição ao comando "d"
	flag.Parse()                                       //processa o input do usuario

	//if *u == "" {
	//	log.Fatal("A URL obrigatoria.")
	//}

	if *p {
		showPersonality()
		return
	}

	url1, err := url.Parse(*u)
	if err != nil {
		log.Fatalf("Erro ao analisar a URL: %v", err)

	}

	wurl = url1.Host

	var urlsLidas, urlsParaLer []string   //cria uma variavel do tipo string para armazenar e separar as urls que ja foram lidas e as que faltam ainda ler
	urlsParaLer = append(urlsParaLer, *u) //le a url de input armazenada no valor u
	for i := 0; i < *d; i++ {             //cria um loop de repetiçao que procura e armazena as urls encontradas na pagina
		urlsLidas = append(urlsLidas, urlsParaLer...)  //subistitui as urls que ja foram lidas para as novas que ainda faltam ler
		urlsParaLer = breadthFirst(crawl, urlsParaLer) //armazena as urls que ja foram lidas das novas que tem para ler
	}

	urlsLidas = append(urlsLidas, urlsParaLer...) //copia os elementos da slice "urlsParaLer" e add no final da slice "urlsLidas"

	if *o != "" {
		writeOutput(urlsLidas, *o)
	} else {
		printOutput(urlsLidas)
	}
}

func showPersonality() {
	fmt.Println("      •• CrawlSpectrum ••")
	fmt.Println("       ••• Lilithh3x •••")
	fmt.Println("")
	fmt.Println("        /\\_/\\")
	fmt.Println("       ( o.o )      <-- Specter")
	fmt.Println("        > ^ <")
	fmt.Println("       /     \\")
	fmt.Println("      |       |")
	fmt.Println("     /| |   | |\\")
	fmt.Println("")
	fmt.Println("Personalidade do Specter:")
	fmt.Println("")
	fmt.Println("1. Curioso:")
	fmt.Println("   Specter é extremamente curioso, sempre explorando novos cantos da web para descobrir informações escondidas. Ele adora investigar links profundos e encontrar dados ocultos.")
	fmt.Println("")
	fmt.Println("2. Silencioso:")
	fmt.Println("   Ele se move silenciosamente pela web, como um verdadeiro espectro. Um mestre na arte da discrição.")
	fmt.Println("")
	fmt.Println("3. Inteligente:")
	fmt.Println("   Com olhos atentos e mente afiada, Specter é capaz de entender rapidamente a estrutura de qualquer site. Ele sabe como evitar armadilhas e navegar com eficiência pelos labirintos digitais.")
	fmt.Println("")
	fmt.Println("4. Leal:")
	fmt.Println("   Specter é leal a seua dona, Lilithh3x. Ele sempre está pronto para ajudar a encontrar informações e dados precisos, sem nunca decepcionar.")
	fmt.Println("")
	fmt.Println("5. Misterioso:")
	fmt.Println("   Pouco se sabe sobre a origem de Specter. Sua presença é notada apenas pelos resultados de seu trabalho impecável. Ele é um enigma, envolto em mistério.")
	fmt.Println("")
	fmt.Println("6. Brincalhão:")
	fmt.Println("   Apesar de sua natureza séria e focada, Specter tem um lado brincalhão. Ele gosta de surpreender com descobertas inesperadas e curiosas durante suas explorações.")
}

func fuzzScan(worklist string) { // definicao da funcao fuzzScan que recebe um argumento 'worklist' d tip string
	readFile, err := os.Open(*w) // abre o arquivo relaciondo ao '*w' (um ponteiro para o nome de arquivo)
	if err != nil {              //exibe uma mensagem de erro caso de erro ao abrir o arquivo
		log.Fatalf("Falha ao ler o arquivo : %s", err) //mensagem de erro
	}

	defer readFile.Close()

	fileScan := bufio.NewScanner(readFile) //cria um scanner para ler o conteudo do arquivo armazenado no '*w'
	fileScan.Split(bufio.ScanLines)        // configura o scanner para dividir em linhas
	var fileTextLines []string             // declara uma slice vazia para armazenar as linhas do arquivo

	for fileScan.Scan() { //intera cada linha usando o scanner
		fileTextLines = append(fileTextLines, fileScan.Text()) // a cada interacao a linha é add a slice 'fleTextLines' usando a funcao 'append'

	}

	for _, eachline := range fileTextLines { // cria um loop para armazenar cada linha a slice 'fileTextLines'
		http.Get(worklist + "/" + eachline) // cria uma url concatenando o valor da url principal com o valor atual da linha 'eachline' do arquivo e em seguida faz uma requisicao get e captura a resosta e armazena em 'resp'
		if err != nil {                     // verifica se ocorreu um erro e armazena em 'err'
			fmt.Println(" A URL Deve ser com 'HTTP://' ou 'HTTPS://'", err) // printa uma mensagem de erro caso tenha algum erro armazenado em 'err'
			log.Fatalln(err)

		}

	}

}

//funçao que cria o arquivo de saida com o resultado da variavel urlsLidas
func writeOutput(urlsLidas []string, o string) { // especifica dois argumentos uma slice de strings 'urlsLidas' e representa 'o' o nome do arquivo de saida
	f, err := os.Create(o) // cria um novo arquivo de saida com o nome que esta armazenado em 'o'.
	if err != nil {
		log.Fatalf("Erro ao criar o arquivo de saída: %s", err)
	}
	defer f.Close() //agenda o fechamento do arquivo

	write := bufio.NewWriter(f) // cria um buffrt associado ao arquivo criado
	for _, url := range urlsLidas {
		write.WriteString(url + "\n") // escreve cada url dentro do arquivo
	}
	write.Flush() // garante que todos os dados foram gravados antes do fechamento do arquivo

}

func printOutput(urlsLidas []string) {
	for _, url := range urlsLidas {
		fmt.Println(url)
	}
}

func breadthFirst(f func(string) []string, worklist []string) []string { //cria uma funçao f do tipo funçao com entrada e saida do tipo string e grava o valor do "urlsLidas" na funçao worklist do tipo string
	seen := make(map[string]bool) //cria um map que armazena um conjunto de chave/valor do tipo string boleano
	for len(worklist) > 0 {       // inicia um loop que exibe as urls ja acessadas
		items := worklist            //copia as urls do worklist para a variavel 'itens'
		worklist = nil               // apaga o valor gravado do worklist para receber os proximos
		for _, item := range items { // cria um par de valores
			if !seen[item] { //verifica se o  'item' ja foi visitado, se nao foi entra no bloco condicional
				seen[item] = true                       // marca o 'item' atual como visitado correspondente no 'seen' como 'true'
				worklist = append(worklist, f(item)...) // regrava o valor atual da funcao f na worklist com as novas urls
			}
		}

	}

	worklist = nil        // apaga "worklist" outra vez
	for k := range seen { // pega o resultado do tipo string e ignora o resultado tipo bool do mapa "seen"
		worklist = append(worklist, k) // reescreve o "worklist" com o resultado do tipo string
	}
	return worklist //retorna o valor armazenado na variavel worklist

}

func crawl(u string) []string { // cria uma entrada que aceita string(url) como argumento e retorna uma slice de strings como resultado
	url2, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}

	if url2.Host != wurl {
		return nil
	}

	list, err := Extract(u) // extrai infos especificas da pagina web da url e armazena o resultado na variavel 'list'
	if err != nil {
		log.Print(err)
	}
	return list //retorna uma lista

}

// Responsavel pelo pre-processamento e pos-processamento em cada no visitado
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
	http.DefaultClient.Timeout = time.Second * 2 //define um timeout de 2 seg para cada request HTTP feito
	//faz um request HTTP e armazena a resposta em 'resp' se ocorrer um erro
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro acessando %s: %s\n", url, err)
	}
	defer resp.Body.Close()
	//printa o status code da url
	fmt.Println(url, "> Status code:", resp.StatusCode)

	//retorna o status code da requisição se der algum erro inesperado (403,404,500 e etc...)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("recebendo %s: %s", url, resp.Status)
	}

	// a função "html.Parse" le uma sequencia de bytes e faz o seu parse e devolve a raiz do documento HTML
	doc, err := html.Parse(resp.Body)
	//resp.Body.Close()
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
