/*
* AUTOR: Rafael Tolosana Calasanz
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2021
* FICHERO: ms.go
* DESCRIPCIÓN: Implementación de un sistema de mensajería asíncrono, insipirado en el Modelo Actor
*/
package ms

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	
	"reflect"
	"github.com/DistributedClocks/GoVector/govec"
)

type Message interface {

}

type MessageSystem struct {
	mbox  chan Message
	peers []string
	done  chan bool
	me    int
}

const (
	MAXMESSAGES = 10000
	TYPEREQUEST = "REQUEST"
	TYPEREPLY = "REPLY"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func parsePeers(path string) (lines []string) {
	file, err := os.Open(path)
	checkError(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

// Pre: pid en {1..n}, el conjunto de procesos del SD
// Post: envía el mensaje msg a pid
func (ms *MessageSystem) Send(pid int, msg Message) {
	Logger := govec.InitGoVector("client", "clientlogfile", govec.GetDefaultConfig())

	conn, err := net.Dial("tcp", ms.peers[pid - 1])
	checkError(err)
	
	fmt.Println("Enviando un 1: ", reflect.TypeOf(msg))
	//fmt.Println("Enviando un 2 con ID: ", msg.(Request).Id)
	outBuf := Logger.PrepareSend("Sending packet", msg, govec.GetDefaultLogOptions())
	fmt.Println("Enviando un 2: ", reflect.TypeOf(outBuf))
	_, err = conn.Write(outBuf)
	if err != nil {
		fmt.Println("GOt a conn write failure, retrying...")
		//conn.Close()
	}
	
	conn.Close()
}

// Pre: True
// Post: el mensaje msg de algún Proceso P_j se retira del mailbox y se devuelve
//		Si mailbox vacío, Receive bloquea hasta que llegue algún mensaje
func (ms *MessageSystem) Receive() (msg Message) {
	msg = <-ms.mbox
	return msg
}

func register(messageTypes []Message){
	for _, msgTp := range messageTypes {
		gob.Register(msgTp)
	}
}

// Pre: whoIam es el pid del proceso que inicializa este ms
//		usersFile es la ruta a un fichero de texto que en cada línea contiene IP:puerto de cada participante
//		messageTypes es un slice con todos los tipos de mensajes que los procesos se pueden intercambiar a través de este ms
func New(whoIam int, usersFile string, messageTypes []Message) (ms MessageSystem) {
	ms.me = whoIam
	ms.peers = parsePeers(usersFile)
	ms.mbox = make(chan Message, MAXMESSAGES)
	ms.done = make(chan bool)
	register(messageTypes)
	
	Logger := govec.InitGoVector("client", "clientlogfileR", govec.GetDefaultConfig())
	fmt.Println("Type of Logger: ", reflect.TypeOf(Logger))
	go func() {
		listener, err := net.Listen("tcp", ms.peers[ms.me-1])
		checkError(err)
		fmt.Println("Process listening at " + ms.peers[ms.me-1])
		defer close(ms.mbox)
		for {
			select {
			case <-ms.done:
				return
			default:
				conn, err := listener.Accept()
				checkError(err)
				
				/*decoder := gob.NewDecoder(conn)
				var msg Message
				err = decoder.Decode(&msg)*/
				
				inBuf := make([]byte, 2048)
				n, errRead := conn.Read(inBuf)
				if errRead != nil {
					fmt.Println("Got a conn read failure, retrying...")
					//conn.Close()
				}
				var msg Message
				Logger.UnpackReceive("Received Message from server", inBuf[0:n], &msg, govec.GetDefaultLogOptions())
								
				conn.Close()
				ms.mbox <- msg
			}
		}
	}()
	return ms
}

//Pre: True
//Post: termina la ejecución de este ms
func (ms *MessageSystem) Stop() {
	ms.done <- true
}

func printErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
