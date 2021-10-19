/*
* AUTOR: Raúl Rustarazo Carmona
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* NIA: 715657
* FECHA: octubre de 2021
* FICHERO: worker.go
* DESCRIPCIÓN: contiene la funcionalidad esencial para implementar el worker
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

func read(conn net.Conn) (idProcess int, interval com.TPInterval){
	decoder := gob.NewDecoder(conn)
	var request com.Request
	
	err := decoder.Decode(&request)
	checkError(err)
	idProcess = (request.Id)		
	interval = (request.Interval)
		
	fmt.Println(interval)
		
	return idProcess, interval
}

func resp(conn net.Conn, listPrimes []int, id int){
	MessageR := com.Reply{Id: id,Primes: listPrimes}
	//fmt.Println("Mensaje de vuelta: ", MessageR)
	fmt.Println("Enviando al cliente el resultado")
	encoder := gob.NewEncoder(conn)

	err := encoder.Encode(MessageR)
	checkError(err)
}

func handle(conn net.Conn){
	timeoutDuration := 2*time.Second
	fmt.Println("Iniciando servidor")
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))
	
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Cliente conectado desde " + remoteAddr)
	id,interval := read(conn)
	newPrimes := FindPrimes(interval)
	fmt.Println(len(newPrimes), " primos encontrados")
	//fmt.Println(newPrimes)
	resp(conn,newPrimes,id)
	conn.Close()
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

func main() {
	vectDirPort := lecturaFichero("./ipWorker.txt")
	fmt.Println(vectDirPort)
	ip, puerto := obtenerIPPuerto(vectDirPort,0)
	fmt.Println("La IP es ", ip)
	fmt.Println("El puerto es ", puerto)
	fmt.Println("En espera por el puerto ", puerto)
	listener, err := net.Listen("tcp", ":"+puerto)
	checkError(err)
	
	fin := false
	for !fin {
		conn, err := listener.Accept()
		checkError(err)
		defer conn.Close()
		handle(conn)
		fmt.Println("Cliente respondido")
	}
	fmt.Println("Servidor finalizado ")
}

