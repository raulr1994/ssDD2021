/*
* AUTOR: Rafael Tolosana Calasanz
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2021
* FICHERO: server.go
* DESCRIPCIÓN: contiene la funcionalidad esencial para realizar los servidores
*				correspondientes a la práctica 1
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
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

// PRE: verdad
// POST: IsPrime devuelve verdad si n es primo y falso en caso contrario
func IsPrime(n int) (foundDivisor bool) {
	foundDivisor = false
	for i := 2; (i < n) && !foundDivisor; i++ {
		foundDivisor = (n%i == 0)
	}
	return !foundDivisor
}

// PRE: interval.A < interval.B
// POST: FindPrimes devuelve todos los números primos comprendidos en el
// 		intervalo [interval.A, interval.B]
func FindPrimes(interval com.TPInterval) (primes []int) {
	for i := interval.A; i <= interval.B; i++ {
		if IsPrime(i) {
			primes = append(primes, i)
		}
	}
	return primes
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

func read(conn net.Conn) (interval com.TPInterval){
	decoder := gob.NewDecoder(conn)
	var request com.Request
	//for{
		err := decoder.Decode(&request)
		checkError(err)		
		interval = (request.Interval)
		
		fmt.Println(interval)
	//}
	return interval
}

func resp(conn net.Conn, listPrimes []int){
	MessageR := com.Reply{Id: 2,Primes: listPrimes}
	fmt.Println("Mensaje de vuelta: ", MessageR)
	encoder := gob.NewEncoder(conn)
	//for{
		err := encoder.Encode(MessageR)
		checkError(err)
	//}
	//conn.Close()
}

func handle(conn net.Conn){
	timeoutDuration := 2*time.Second
	fmt.Println("Iniciando servidor")
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))
	
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Cliente conectado desde " + remoteAddr)
	interval := read(conn)
	newPrimes := FindPrimes(interval)
	fmt.Println("Los primos encontrados son")
	fmt.Println(newPrimes)
	resp(conn,newPrimes)
}

func obtenerIPPuerto(vectDirPort [] string, pos int) (ip string, puerto string){
	s := strings.Split(vectDirPort[pos],":")
	ip = s[0] //La ip
	puerto = s[1] //El puerto
	return ip, puerto
}

func lecturaFichero(nameFile string) (vectDirPort [] string){
	file, err := os.Open(nameFile)
	
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}
	
	fileScanner := bufio.NewScanner(file)
	
	//vectDirPort = [] string{}
	
	for fileScanner.Scan(){
		//fmt.Println(fileScanner.Text())
		vectDirPort = append(vectDirPort,fileScanner.Text())
	}
	
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}
	
	file.Close()
	return vectDirPort
}

func onWorkers(id int, reqChan chan com.Request, repChan chan com.Reply){
	for{
		request := <- reqChan
		fmt.Println("Worker con el id ", id, " procesando peticion" , request)
		primes := []int{1,2,3,4,5}
		replyTest := com.Reply{Id: 2,Primes: primes}
		repChan <- replyTest
	}
}



type Job struct {
	Conn		net.Conn
	Request	com.Request
}

type Result struct {
	Conn		net.Conn
	Reply		com.Reply
}

func readRequest(conn net.Conn) (request com.Request){
	decoder := gob.NewDecoder(conn)
	var request com.Request
	//for{
		err := decoder.Decode(&request)
		checkError(err)		
		request = (request.Interval)
		
		fmt.Println(interval)
	//}
	return request
}

func respRequest(MessageR com.Reply){
	fmt.Println("Mensaje de vuelta: ", MessageR)
	encoder := gob.NewEncoder(conn)
	//for{
		err := encoder.Encode(MessageR)
		checkError(err)
	//}
	//conn.Close()
}

func handleRequests(conn net.Conn, reqChan chan Job, repChan chan Result){
	for{
		//Obtener datos peticion
		NewRequest := readRequest(conn)
		RequestToWorker := Job{Conn: conn, Request: NewRequest}
		//Mandar a worker datos
		reqChan <- RequestToWorker
		//Recibir resultados del worker(i)
		reply := <- repChan
		//Responder
	}
}

func main() {
	vectDirPort := lecturaFichero("./ipServer.txt")
	fmt.Println(vectDirPort)
	ip, puerto := obtenerIPPuerto(vectDirPort,0)
	fmt.Println("La IP es ", ip)
	fmt.Println("El puerto es ", puerto)
	fmt.Println("En espera por el puerto ", puerto)
	
	nWorkers := 3
	
	reqChan := make(chan Job)
	repChan := make(chan Result)
	
	for i:= 0; i < nWorkers; i++ {
		fmt.Println("Creando el worker ", i)
		go onWorkers(i, reqChan, repChan)
	}
	
	/*interval := com.TPInterval{10,70}
	requestTest := com.Request{Id: 2, Interval: interval}
	reqChan <- requestTest
	reply := <- repChan
	fmt.Println("Resultado obtenido ", reply)
	reqChan <- requestTest
	reply = <- repChan
	fmt.Println("Resultado obtenido ", reply)
	reqChan <- requestTest
	reply = <- repChan
	fmt.Println("Resultado obtenido ", reply)*/
	
	listener, err := net.Listen("tcp", ":"+puerto)
	checkError(err)
	fin := false
	for !fin {
		conn, err := listener.Accept()
		checkError(err)
		defer conn.Close()
		go handleRequests(conn)
	}
	fmt.Println("Servidor finalizado ")
}
