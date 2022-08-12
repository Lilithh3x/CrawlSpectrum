package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	//pacote nao padrao que oferece um parser HTML
	"golang.org/x/net/html"
)

func main() {
	breadthFirst(crawl, os.Args[1:]) // chama a função crawl e executa junto com um argumento de linha de comando

}

type UsageSection struct {
	Nome          string
	Descrição     string
	Flags         []UsageFlag
	Hidden        bool
	ExpectedFlags []string
}

func (u *UsageSection) PrintSection(max_length int, extended bool) {
	if !extended && u.Hidden {
		return
	}
	fmt.Printf("%s:\n", u.Nome)
	for _, f := range u.Flags {
		f.PrintFlag(max_length)
	}
	fmt.Printf("\n")
}

type UsageFlag struct {
	Name      string
	Descrição string
	Default   string
}

func (f *UsageFlag) PrintFlag(max_length int) {
	format := fmt.Sprintf("  -%%-%ds", max_length)
	if f.Default != "" {
		format = format + "(default: %s)\n"
		fmt.Printf(format, f.Name, f.Descrição, f.Default)
	} else {
		format = format + "\n"
		fmt.Printf(format, f.Name, f.Descrição)
	}
}

func Usage() {
	u_http := UsageSection{
		Nome:          "HTTP OPTIONS",
		Descrição:     "Options contolling the HTTP request and its parts.",
		Flags:         make([]UsageFlag, 0),
		Hidden:        false,
		ExpectedFlags: []string{"u", "X", "timeout"},
	}
	u_general := UsageSection{
		Nome:          "GENERAL OPTIONS",
		Descrição:     "",
		Flags:         make([]UsageFlag, 0),
		Hidden:        false,
		ExpectedFlags: []string{"t", "v"},
	}
	u_output := UsageSection{
		Nome:          "OUTPUT OPTIONS",
		Descrição:     "Options for output. Output file formats and file names .",
		Flags:         make([]UsageFlag, 0),
		Hidden:        false,
		ExpectedFlags: []string{"o", "of"},
	}
	sections := []UsageSection{u_http, u_general, u_output}

	max_length := 0
	flag.VisitAll(func(f *flag.Flag) {
		found := false
		for i, section := range sections {
			if strInSlince(f.Name, section.ExpectedFlags) {
				sections[i].Flags = append(sections[i].Flags, UsageFlag{
					Name:      f.Name,
					Descrição: f.Usage,
					Default:   f.DefValue,
				})
				found = true
			}
		}
		if !found {
			fmt.Printf("DEBUG: Flag %s was found but not defined in craw.go.\n", f.Name)
			os.Exit(1)
		}
		if len(f.Name) > max_length {
			max_length = len(f.Name)
		}
	})

	for _, section := range sections {
		section.PrintSection(max_length, false)
	}
}

func strInSlince(val string, slice []string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func breadthFirst(f func(string) []string, worklist []string) { //cria uma funçao f do tipo funçao com entrada e saida do tipo string
	// e grava o valor do "so.Args[1:]" na funçao worklist do tipo string
	seen := make(map[string]bool) //cria um map que armazena um conjunto de chave/valor do tipo string boleano
	for len(worklist) > 0 {
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
}

func crawl(url string) []string { // cria uma entrada(url) e uma saida do tipo string
	fmt.Println(url) //printa o valor que foi gravado na url
	list, err := Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
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
# Web-Craw
# Web-Craw
# Web-Craw
# Web-Craw
# Web-Craw
