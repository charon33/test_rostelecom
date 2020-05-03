package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"test_rostelecom/models"
)
const gorutineCount int = 5
const requiredLine string = "Go"

func makeRequest(site *models.Site) string{
	resp, err := http.Get(site.Url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	response := string(body)
	return response
}
	//Валидация урлов
func testUrls(siteSlises []string) []string{
	var patternArr = [4]string{
		"http://\\w+\\.\\w+$",
		"^https://\\w+\\.\\w+$",
		"www.\\w+\\.\\w+$",
		"^\\w+\\.\\w+$",
	}
	resultArr := make([]string,0)
	if len(siteSlises) != 0 {

		for _, url := range siteSlises {
			matched := false
			for _, pattern := range patternArr{
				compilPattrn, err := regexp.Compile(pattern)
				if err != nil {
					log.Fatalln(err)
				}
				if compilPattrn.MatchString(url) {
					matched = true
				}

			}
			if matched {
				resultArr = append(resultArr, url)
			}
		}
	}
	return resultArr
}
	//Подсчет количества вхождений
func CountEntries(response string, requiredLine string) int{
	var count int
	if len(response) != 0 {
		count = strings.Count(response, requiredLine)
	} else {
		log.Println("response line is null")
	}
	return count
}
	//Из массива провереных урлов создаем массив экземпляров сайтов
func createSiteExampl(sitesArr []string) []models.Site{
	resultArr := make([]models.Site,0)
	for _, site := range sitesArr {
		siteExampl := models.Site{Url: site}
		resultArr = append(resultArr, siteExampl)
	}
	return resultArr
}
	//Функция для сборки функций и методов для запуска через горутины
func AgregateFunctions(example models.Site, summ *int, waitG *sync.WaitGroup) {
	defer waitG.Done()
	response := makeRequest(&example)
	countEnt := CountEntries(response, requiredLine)
	example.ChangeEntries(countEnt)
	*summ = *summ + countEnt
	example.PrintCounts()
}

func main() {
	summ := new(int)
	waitG := new(sync.WaitGroup)
		//считываем аргументы стдин
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString(' ')
		//делим строку на отдельные урлы
	sitesSlices := strings.Split(text, "\n")
		//Проводим валидацию
	sitesSlices = testUrls(sitesSlices)
		//Получаем массив экземпляров структуры site
	siteExampls := createSiteExampl(sitesSlices)

	//Прошу прощения за костыль. Я понял, что нужно использовать channels,
	//но,к сожалению, голэнг начал изучать и выполнять тестовое только сегодня
	//с многими вещами разобраться не успел
	for i:=0; i < len(siteExampls); i =i + gorutineCount {
		if i + gorutineCount > len(siteExampls){
			for _, exampl := range siteExampls {
				waitG.Add(1)
				go AgregateFunctions(exampl, summ, waitG)
			}
		} else {
			for _, exampl := range siteExampls[i : i+gorutineCount] {
				waitG.Add(1)
				go AgregateFunctions(exampl, summ, waitG)
			}

		}
		waitG.Wait()
	}

	fmt.Print("Total: ", *summ, "\n")

}
