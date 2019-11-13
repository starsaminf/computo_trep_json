	package main

	import (
		"bufio"
		"fmt"
		"io/ioutil"
		"log"
		"net/http"
		"os"
		"strings"
	)

	func fileExists(system, filename string) bool {
		info, err := os.Stat(system+"/"+filename)
		if os.IsNotExist(err) {
			return false
		}
		return !info.IsDir()
	}

	func fileSave(system, mesa, json_ string){
		f, err := os.Create(system+"/"+mesa+".json")
		if err != nil {
			fmt.Println(err)
			return
		}
		l, err := f.WriteString(json_)
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
		fmt.Println(l, " -> "+mesa +"Ok")
		err = f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}

	}

	func main(){
		//Descargar todos los datos del trep o computo para obtener todas las mesas del pais
		mesa := ""
		var systemOep [2] string
		systemOep[0] = "trep"
		systemOep[1] = "computo"

		for index := 0; index < 2 ; index++  {
			file, err := os.Open("mesas.txt")
			if err != nil {
				log.Fatal(err)
			}
			file.Close()
			scanner := bufio.NewScanner(file)

			for scanner.Scan() {

				mesa = scanner.Text()
				response, err := http.Get("https://"+systemOep[index]+".oep.org.bo/resul/resulActa/"+mesa+"/1")
				if err != nil {
					log.Fatal(err)
				}

				response.Body.Close()

				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					log.Fatal(err)
				}

				if(strings.Contains(string(body), "ERROR") == false){
					if fileExists(systemOep[index],mesa) ==  true{
						log.Printf("Duplicado : " +mesa)
						fileSave(systemOep[index],mesa+"_2", string(body))
					}else{
						fileSave(systemOep[index],mesa, string(body))
					}
				}else if(strings.Contains(string(body), "Mesa no computada") == true){
					log.Printf(" Mesa no computada : " +mesa)
				}
			}
		}


	}