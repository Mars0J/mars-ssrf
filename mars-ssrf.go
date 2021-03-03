package main


import (
"fmt"
"time"
"sync"
"bufio"
"os"
"strings"
"io/ioutil"
"net/http"
"net/url"
"flag"
"log"
)

func main() {

	banner:=`


 	  __  __    _    ____  ____    ____ ____  ____  _____
       	 |  \/  |  / \  |  _ \/ ___|  / ___/ ___||  _ \|  ___|
	 | |\/| | / _ \ | |_) \___ \  \___ \___ \| |_) | |_
	 | |  | |/ ___ \|  _ < ___) |  ___) |__) |  _ <|  _|
  	 |_|  |_/_/   \_\_| \_\____/  |____/____/|_| \_\_|
				1.0 - @Mars

		`

	fmt.Printf("%s\n\n", banner)




	var concurrency int
	var payloads string
	var match string
	var paths bool
	var appendMode bool
	var silent bool
	flag.IntVar(&concurrency, "c", 30, " Daha Yüksek Bir Hız için Concurrency'i Ayarlayın ")
	flag.StringVar(&payloads, "pL", "", " Payload Listesi ")
	flag.StringVar(&match, "m", "", "Yanıtı Bir Kalıpla eşleştirin (e.g) 'Success:'")
	flag.BoolVar(&appendMode, "a", false, "Payload'ı Parametreye Ekleyin")
	flag.BoolVar(&paths, "p", false, "Yalnızca SSRF'yi test edin.")
	flag.BoolVar(&silent, "s", false, "Yalnızca Savumasız Ana Bilgisayarları Yazdır")
	flag.Parse()

	if payloads != "" {

		var wg sync.WaitGroup
		for i:=0; i<=concurrency; i++ {
			wg.Add(1)
			go func () {
				test_ssrf(payloads, match, appendMode, silent, paths)
				wg.Done()
			}()
			wg.Wait()
		}

	}

}

func test_ssrf(payloads string, match string, appendMode bool, silent bool, paths bool) {

	file,err := os.Open(payloads)

	if err != nil {
		log.Fatal("Dosya Okunamadı.")
	}

	defer file.Close()

	time.Sleep(time.Millisecond * 10)
	scanner:=bufio.NewScanner(os.Stdin)

	pScanner:=bufio.NewScanner(file)

	for scanner.Scan() {
		for pScanner.Scan() {
			link:=scanner.Text()
			payload:=pScanner.Text()

			u,err :=  url.Parse(link)

			if err != nil {
				return
			}
			if paths == false {

				qs:= url.Values{}
				for param, vv := range u.Query() {
					if appendMode {
						qs.Set(param, vv[0]+payload)
					}else {
						qs.Set(param, payload)
					}
				}



				u.RawQuery = qs.Encode()
				if silent == false {
					fmt.Printf("[+] Testing: \t  %s\n", u)
				}
				make_request(u.String(), match)
			}else {
				newLink:=link+payload
				fmt.Println(newLink)
				if silent == false {
                                        fmt.Printf("[+] Testing: \t  %s\n", newLink)
                                }
                                make_request(newLink, match)

			}

		}
	}
}


func make_request(url string, match string) {
	client:=&http.Client{}
	req,err:= http.NewRequest("GET", url, nil)
	if err != nil {
		return
        }

	resp,err:=client.Do(req)
	if err != nil {
		return
	}

	bodyBytes,err  := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
	if strings.Contains(bodyString, match) {
		fmt.Println(url + " Savumasız")
	} 


}

