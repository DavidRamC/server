package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

//estructura que representa las respuestas del servidor
type Response struct {
	Ok      bool
	Message string
	Package []byte
}

func main() {
	//Se conecta al servidor
	serveCnn, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	go writeMessage(serveCnn)
	readMessages(serveCnn)
}

//lee las respuestas a las peticiones realizadas
//recibe el archivo y lo crea (en bytes)
func readMessages(cnn net.Conn) {
	var i int64 = 0
	for {
		var response Response

		err := gob.NewDecoder(cnn).Decode(&response)

		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(response.Message)
		if response.Package != nil {
			i++
			valueString := fmt.Sprint(i)
			name := "received_files/000" + valueString
			_, err := os.Create(name)
			if err != nil {
				fmt.Println("error en la creacion del archivo")
				fmt.Println(err)
				continue
			}
			os.WriteFile(name, response.Package, 0666)

		}
	}
}

//representa el CLI, el cual se ejecuta concurrentemente
func writeMessage(conn net.Conn) {
	var scanner = bufio.NewReader(os.Stdin)
	fmt.Println("Cliente.....")
	fmt.Println("")
	fmt.Println("-----------------")
	fmt.Println("")

	for {
		fmt.Print("-->")
		text, _ := scanner.ReadString('\n')
		text = normalize(text)
		params := strings.Split(text, " ")
		flag := validateCh(params)
		flag2 := validate(params)

		if flag && flag2 {
			err := gob.NewEncoder(conn).Encode(text)
			if err != nil {
				fmt.Println("Error al codificar la informacion")
				return
			}
			if params[0] == "exit" {
				conn.Close()
				os.Exit(0)
			}

			if params[0] == "send_file" {
				data, err := ioutil.ReadFile(params[2])
				if err != nil {
					fmt.Println(err)
					continue
				}
				conn.Write(data)
			}
		} else {
			continue
		}

	}

}

//valida que el canal ingresado si sea valido
func validateCh(params []string) bool {
	flag := true
	if len(params) == 2 {
		switch params[1] {
		case "-ch1":
		case "-ch2":
		case "-ch3":
		case "-ch4":
		case "-ch5":
		default:
			flag = false
		}
	}
	return flag

}

//elimina los caracteres no deseados para facilitar el manejo de los comandos
func normalize(c string) string {
	c = strings.Replace(c, "\n", "", 1)
	c = strings.Replace(c, "\r", "", 1)
	c = strings.Replace(c, "-->", "", 1)
	return c
}

//valida la cantidad de parametros ingresados
func validate(params []string) bool {
	var v bool = true
	if params[0] == "sendFile" && len(params) != 3 {
		v = false
	} else if (params[0] == "receive" || params[0] == "suscribe") && len(params) != 2 {
		v = false
	}
	return v
}
