package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var comand string = ""

//clientes conectados
var clients []*Client

const INFO string = `Lista de comandos - cliente:

send_file -chn ruta  ---> Envia el archivo especificado en el parametro ruta por el canal correspondiente
ej: send_file -ch3 ./mi-archivo.txt

receive -chn ---> Permite al cliente recibir la informacion que se quiera enviar por el canal especificado
ej: receive -ch1

suscribe -chn ---> Permite al cliente recibir toda la informacion que se envie por el canal especificado
ej: suscribe -ch5

info ---> muestra informacion de los comandos 

exit ---> cierra la conexion del cliente

NOTA: en los metodos donde se requiere un canal la n representa el numero del canal el cual debe estar en un rango de 1 a 5 incluyendo extremos`

type Client struct {
	Conn        net.Conn
	Channel     string
	Suscription string
}
type Response struct {
	Message string
	Ok      bool
	Package []byte
}

//ejecuto el cli para procesar los comandos del usuario, en el momento
//en el que el usuario digite "start" se lanza una go routine para ejecutar el servidor concurrentemente
func main() {

	var reader *bufio.Reader = bufio.NewReader(os.Stdin)
	fmt.Println("File server")

	fmt.Println("-------------------------")
	for {
		fmt.Print("-->")
		comand, _ := reader.ReadString('\n')
		comand = strings.Replace(comand, "\n", "", 1)
		comand = strings.Replace(comand, "\r", "", 1)
		comand = strings.ReplaceAll(comand, " ", "")
		comand = strings.ToLower(comand)

		switch comand {
		case "start":
			go server()
		case "stadistics":
			showStadistcs()
		case "help":
			fmt.Println("LISTA DE COMANDOS")
			fmt.Println("start ----> Inicia el servidor en el puerto 8080")
			fmt.Println("stadistics ----> Muestra estadisticas del servidor como numero de conecciones, entre otras.")
			fmt.Println("help ----> Muestra informacion de los comandos")
			fmt.Println("exit ----> Detiene la aplicacion")
		case "exit":
			return
		default:
			fmt.Println("el comando no se reconoce, por favor digita 'help' para mas informacion acerca de los comandos")
		}
	}
}

//servidor corriendo y aceptando cada cliente
func server() {
	server, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	for {
		client, err := server.Accept()
		if err != nil {
			fmt.Println("hubo problemas en la conexion del cliente")
			continue
		}
		//gouroutine para manejar cada conexion
		go handleConnection(client)

	}

}

//se da manejo a la conexion
func handleConnection(client net.Conn) {
	c := Client{
		Conn:        client,
		Channel:     "",
		Suscription: "",
	}
	clients = append(clients, &c)
	for {
		var cmd string
		err := gob.NewDecoder(client).Decode(&cmd)
		if err != nil {
			fmt.Println("Hubo un error en la decodificacion del paquete")
			fmt.Println(err)
			return
		}

		process_request(cmd, &c)
	}

}

//se procesa la solicitud del cliente y se le da respuesta a la misma
func process_request(cmd string, c *Client) {
	params := strings.Split(cmd, " ")
	file := make([]byte, 5242880) //buffer de 5mb
	switch params[0] {
	case "send_file":
		if len(params) != 3 {
			resp := Response{
				Message: "Parametros insuficientes, por favor digita info para mas informacion acerca de los comandos ",
				Ok:      false,
				Package: nil,
			}
			err := gob.NewEncoder(c.Conn).Encode(resp)
			if err != nil {
				println("Error en la codificacion de la informacion")
				return
			}
			return
		}
		fmt.Println("entre al caso")

		c.Conn.Read(file)
		sendFile(params[1], file, c)

	case "receive":
		c.Channel = params[1]
		resp := Response{
			Message: "Se encuentra disponible para recibir datos",
			Ok:      false,
			Package: nil,
		}
		err := gob.NewEncoder(c.Conn).Encode(resp)
		if err != nil {
			println("Error en la codificacion de la informacion")
			return
		}
	case "suscribe":
		c.Suscription = params[1]
		resp := Response{
			Message: "Se ha suscrito al canal",
			Ok:      true,
			Package: nil,
		}
		err := gob.NewEncoder(c.Conn).Encode(resp)
		if err != nil {
			println("Error en la codificacion de la informacion")
			return
		}
	case "info":
		r := Response{
			Message: INFO,
			Ok:      false,
			Package: nil,
		}
		err := gob.NewEncoder(c.Conn).Encode(r)
		if err != nil {
			println("Error en la codificacion de la informacion")
			return
		}
	case "exit":
		c.Conn.Close()
	default:
		r := Response{
			Message: "El comando no se reconoce, por favor digite info para mas informacion acerca de los comandos",
			Ok:      false,
			Package: nil,
		}
		err := gob.NewEncoder(c.Conn).Encode(r)
		if err != nil {
			println("Error en la codificacion de la informacion")
			return
		}

	}
}

//Se envia el archivo a los clientes suscritos a un canal o a los que estan disponibles para recibir en el mismo
//En el caso de que no hayan suscripciones o clientes en modo "receive" no se envia el archivo

func sendFile(channel string, file []byte, sender *Client) {
	cont := 0
	for _, c := range clients {
		if (channel == c.Channel) || (channel == c.Suscription) {
			r := Response{
				Message: "Archivo recibido",
				Ok:      true,
				Package: file,
			}
			err := gob.NewEncoder(c.Conn).Encode(r)
			if err != nil {
				println("Error en la codificacion de la informacion")
				return
			}
			cont++
		}
	}
	if cont == 0 {
		r := Response{
			Message: "El archivo no ha sido enviado porque no hay clientes suscritos a este canal o no hay ninguno en modo 'receive' en este canal, intentalo de nuevo mas tarde.",
			Ok:      false,
			Package: nil,
		}
		err := gob.NewEncoder(sender.Conn).Encode(r)
		if err != nil {
			println("Error en la codificacion de la informacion")
			return
		}
	}

}

func showStadistcs() {
	panic("unimplemented")
}
