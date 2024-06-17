package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

// Cores para o terminal
const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorReset  = "\033[0m"
)

func main() {
	var domainListFile string
	var singleDomain string
	var threads int
	var outputFile string
	var cpanelPort int

	flag.StringVar(&domainListFile, "l", "", "Arquivo contendo a lista de URLs/domínios")
	flag.StringVar(&singleDomain, "d", "", "Domínio único para verificar (sem protocolo, exemplo: example.com)")
	flag.IntVar(&threads, "t", 10, "Número de threads/concorrência")
	flag.StringVar(&outputFile, "o", "cpnalvalid.txt", "Nome do arquivo de saída para domínios válidos")
	flag.IntVar(&cpanelPort, "p", 2083, "Porta para verificar o cPanel (padrão: 2083)")
	flag.Parse()

	if singleDomain != "" {
		verifySingleDomain(singleDomain, cpanelPort, outputFile)
	} else if domainListFile != "" {
		verifyMultipleDomains(domainListFile, threads, outputFile, cpanelPort)
	} else {
		fmt.Println("Erro: Especifique um domínio único com -d ou um arquivo de lista de URLs/domínios com -l.")
		os.Exit(1)
	}
}

func verifySingleDomain(domain string, port int, outputFile string) {
	if strings.HasPrefix(domain, "https://") {
		domain = strings.TrimPrefix(domain, "https://")
	}

	if strings.HasPrefix(domain, "http://") {
		fmt.Println("Erro: O protocolo http não é suportado. Use apenas domínios sem protocolo (exemplo: example.com).")
		os.Exit(1)
	}

	if isValidCPanel(domain, port) {
		fmt.Printf("%s[+] %s:%d: cPanel Login encontrado%s\n", ColorGreen, domain, port, ColorReset)
		writeValidURL(domain+":"+fmt.Sprint(port), outputFile)
	} else {
		fmt.Printf("%s[~] %s:%d: Não é um cPanel válido%s\n", ColorRed, domain, port, ColorReset)
	}
}

func verifyMultipleDomains(filename string, threads int, outputFile string, port int) {
	urls, err := readURLList(filename)
	if err != nil {
		fmt.Printf("Erro ao ler o arquivo de URLs/domínios: %v\n", err)
		os.Exit(1)
	}

	urlChan := make(chan string)

	wg := sync.WaitGroup{}

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range urlChan {
				if isValidCPanel(url, port) {
					fmt.Printf("%s[+] %s:%d | cPanel GOOD%s\n", ColorGreen, url, port, ColorReset)
					writeValidURL(url+":"+fmt.Sprint(port), outputFile)
				} else {
					fmt.Printf("%s[~] %s:%d | BAD%s\n", ColorRed, url, port, ColorReset)
				}
			}
		}()
	}

	for _, url := range urls {
		urlChan <- url
	}

	close(urlChan)

	wg.Wait()
}

func readURLList(filename string) ([]string, error) {
	var urls []string
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url != "" {
			urls = append(urls, url)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func isValidCPanel(domain string, port int) bool {
	resp, err := http.Get(fmt.Sprintf("http://%s:%d", domain, port))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	title := getTitleFromHTML(body)
	return title == "cPanel"
}

func getTitleFromHTML(html []byte) string {
	titleStart := strings.Index(string(html), "<title>")
	if titleStart == -1 {
		return ""
	}
	titleEnd := strings.Index(string(html), "</title>")
	if titleEnd == -1 || titleEnd <= titleStart {
		return ""
	}
	return strings.TrimSpace(string(html[titleStart+len("<title>") : titleEnd]))
}

func writeValidURL(url, outputFile string) {
	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Erro ao abrir o arquivo de saída: %v\n", err)
		return
	}
	defer f.Close()
	if _, err := f.WriteString(url + "\n"); err != nil {
		fmt.Printf("Erro ao escrever no arquivo de saída: %v\n", err)
	}
}
