/*
* AUTOR: Raúl Rustarazo Carmona
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* NIA: 715657
* FECHA: octubre de 2021
* FICHERO: server-draft.go
* DESCRIPCIÓN: contiene la funcionalidad esencial para implementar el master para la arquitectura master-worker
*				
*/
package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"io"
	//"./com"
	"example.com/com"
	"bufio"
	"log"
	"strings"
	//"strconv"
	//"bytes"
	"time"
	"example.ssh/ssh"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func logerr(err error) bool {
	if err != nil {
		netErr, ok := err.(net.Error)
		if ok && netErr.Timeout() {
			fmt.Println("read timeout:", err)
		} else if err == io.EOF {
			fmt.Println("read EOF:", err)
		} else{
			fmt.Println("read error:", err)
		}
		return true
	}
	return false
}

func obtenerIPPuerto(vectDirPort [] string, pos int) (ip string, puerto string){
	s := strings.Split(vectDirPort[pos],":")
	ip = s[0] //La ip
	puerto = s[1] //El puerto
	return ip, puerto
}

func lecturaFichero(nameFile string) (vectDirPort [] string, nworkers int){
	file, err := os.Open(nameFile)
	
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}
	
	fileScanner := bufio.NewScanner(file)
	
	//vectDirPort = [] string{}
	nworkers = 0
	for fileScanner.Scan(){
		//fmt.Println(fileScanner.Text())
		vectDirPort = append(vectDirPort,fileScanner.Text())
		nworkers++
	}
	
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}
	
	file.Close()
	return vectDirPort,nworkers
}

func read(conn net.Conn) (idProcess int, interval com.TPInterval){
	decoder := gob.NewDecoder(conn)
	var request com.Request
	//for{
		err := decoder.Decode(&request)
		checkError(err)
		idProcess = (request.Id)		
		interval = (request.Interval)
		
		fmt.Println(interval)
	//}
	return idProcess, interval
}

func resp(conn net.Conn, listPrimes []int, id int){
	MessageR := com.Reply{Id: id,Primes: listPrimes}
	//fmt.Println("Mensaje de vuelta al cliente: ", MessageR)
	fmt.Println("Enviando al cliente el resultado")
	encoder := gob.NewEncoder(conn)
	//for{
		err := encoder.Encode(MessageR)
		checkError(err)
	//}
	//conn.Close()
}

func sendToWorker(conn net.Conn, id int, interval com.TPInterval){
	request := com.Request{Id: id,Interval: interval}
	fmt.Println("Enviando peticion al worker: ", request)
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(request) //Enviando al worker
	checkError(err)
}

func receiveFromWorker(conn net.Conn) (primes [] int){
	decoder := gob.NewDecoder(conn)
	var reply com.Reply
	//for{
		err := decoder.Decode(&reply)
		checkError(err)		
		primes = (reply.Primes)
		
		//fmt.Println(primes)
	//}
	return primes
}

func handle(conn net.Conn, dirWorker string, connW net.Conn){
	timeoutDuration := 2*time.Second
	fmt.Println("Iniciando servidor")
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))
	
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Cliente conectado desde " + remoteAddr)
	id,interval := read(conn) //Obtengo los datos de la peticion del cliente
	
    	//Enviando datos al worker
    	sendToWorker(connW,id,interval)
    	
    	//Recibiendo datos del worker
	fmt.Println("Recibiendo datos del worker con dir ", dirWorker)
	newPrimes := receiveFromWorker(connW)
	fmt.Println(len(newPrimes), " primos encontrados")
	//fmt.Println(newPrimes)
	
	//Respondiendo al cliente
	resp(conn,newPrimes,id)
}

var jobs = make(chan net.Conn)

func onWorkers(id int, newDirWorker string){
	idWorker := id
	dirWorker := newDirWorker
	
	//Estableciendo conexión con el Worker
	
	tcpAddr, err := net.ResolveTCPAddr("tcp", dirWorker)
	checkError(err)

	for job := range jobs {
		

		connW, err := net.DialTCP("tcp", nil, tcpAddr)
		checkError(err)
		
		fmt.Println("Conectado con el worker ", idWorker, " con dir ", dirWorker)
	
		handle(job,dirWorker,connW)
	}
}

func executeComand (nombreusuario string, ipjhost string, comand string) {
	ssh, err := sshcode.NewSshClient(
		nombreusuario,
		ipjhost,
		22,
		"./rsa",
		"")

	if err != nil {
			log.Printf("SSH init error %v", err)
	} else {
		output, err := ssh.RunCommand(comand)
		fmt.Println(output)
		if err != nil {
			log.Printf("SSH run command error %v", err)
		}
	}
}

func main() {
	vectDirPort, _ := lecturaFichero("./ipServer.txt")
	fmt.Println(vectDirPort)
	ip, puerto := obtenerIPPuerto(vectDirPort,0)
	fmt.Println("La IP es ", ip)
	fmt.Println("En espera por el puerto ", puerto)
	
	vectDirWorkers, nWorkers := lecturaFichero("./ipWorkers.txt")
	fmt.Println(vectDirWorkers)
	fmt.Println(nWorkers, "nWorkers creados")
	//ip, puerto := obtenerIPPuerto(vectDirPort,0)
	
	for i:= 0; i < nWorkers; i++ { //Despertando los worker
		nombreusuario := "raulrcuni"
		fmt.Println("Despertando el worker ", i)
		ip,_ := obtenerIPPuerto(vectDirWorkers,i)
		executeComand(nombreusuario,ip,"chmod 777 ./worker")
		go executeComand(nombreusuario,ip,"./worker")
	}
	
	for i:= 0; i < nWorkers; i++ {
		fmt.Println("Creando el worker ", i)
		go onWorkers(i,vectDirWorkers[i])
	}
	
	listener, err := net.Listen("tcp", ":"+puerto)
	checkError(err)
	
	fin := false
	for !fin {
		fmt.Println("En espera por el puerto ", puerto)
		conn, err := listener.Accept()
		checkError(err)
		//defer conn.Close()
		jobs <- conn
	}
	fmt.Println("Servidor finalizado ")
}

