package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const numberOfChecks = 3
const delayChecks = 5;

func main() {
	showIntroduction()

	for {
		showMenu()
		readSitesFromArchive()

		comand := readCommand()

		switch comand {
		case 1:
			startMonitoring()
		case 2:
			fmt.Println("Exibindo Logs...")
			fmt.Println("")
			printLogs()
		case 0:
			fmt.Println("Saindo do programa...")
			os.Exit(0)
		default:
			fmt.Println("Comando desconhecido")
			os.Exit(-1)
		}
	}
}

func showIntroduction() {
	version := 1.0
	fmt.Println("Olá!")
	fmt.Println("Este programa está na versão", version)
	fmt.Println("---")
}

func showMenu() {
	fmt.Println("1- Iniciar Monitoramento")
	fmt.Println("2- Exibir Logs")
	fmt.Println("0- Sair do programa")
	fmt.Println("---")
	fmt.Println("")
}

func readCommand() int {
	var userComand int
	fmt.Scan(&userComand)
	fmt.Println("O comando escolhido foi", userComand)

	return userComand
}

func startMonitoring() {
	fmt.Println("Monitorando...")

	sites := readSitesFromArchive()

	for index := 0; index < numberOfChecks; index += 1 {
		for _, site := range sites {
			testSites(site)
		}

		fmt.Println("")
		time.Sleep(delayChecks * time.Second)
	}
	
	fmt.Println("")
}

func readSitesFromArchive() []string {
	var sites []string
	
	file, err := os.Open("sites.txt")
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}
	
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		
		line = strings.TrimSpace(line)

		sites = append(sites, line)
	}

	file.Close()

	return sites
}

func testSites(site string) {
	resp, err := http.Get(site)

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("O site", site, "está OK. Status Code:", resp.StatusCode)
		registerLog(site, true)
	} else {
		fmt.Println("O site", site, "está com problemas. Status Code:", resp.StatusCode)
		registerLog(site, false)
	}
}

func registerLog(site string, isOnline bool) {
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	file.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - online: " + strconv.FormatBool(isOnline) + "\n" )

	file.Close()
}

func printLogs() {
	file, err := ioutil.ReadFile("log.txt")
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	fmt.Println(string(file))
}
