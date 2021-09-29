/*
* AUTOR: Rafael Tolosana Calasanz
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2021
* FICHERO: server.go
* DESCRIPCIÓN: contiene la funcionalidad esencial para realizar los servidores
*				correspondientes al trabajo 1
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
	"flag"
	"bytes"
	"time"
)

/*
	Imprime por pantalla el error ocurrido
*/
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

type Message struct {
	ID string
	Data string
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
	tmp := make([]byte, 500)
	//interval := new(com.TPInterval)
	for{
		_, err := conn.Read(tmp)
		if logerr(err){
			break
		}
		
		tmpbuff := bytes.NewBuffer(tmp)
		tmpstruct := new(com.Request)
		
		gobobj := gob.NewDecoder(tmpbuff)
		gobobj.Decode(tmpstruct)
		
		interval = (tmpstruct.Interval)
		
		fmt.Println(tmpstruct.Interval)
	}
	return interval
}

func resp(conn net.Conn, listPrimes []int){
	MessageR := com.Reply{Id: 2,Primes: listPrimes}
	bin_buf := new(bytes.Buffer)
	
	gobobje := gob.NewEncoder(bin_buf)
	gobobje.Encode(MessageR)
	
	conn.Write(bin_buf.Bytes())
	conn.Close()
}

func handle(conn net.Conn){
	//defer conn.Close()
	timeoutDuration := 2*time.Second
	fmt.Println("Iniciando servidor")
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))
	
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Cliente conectado desde " + remoteAddr)
	//conn.Close()
	interval := read(conn)
	newPrimes := FindPrimes(interval)
	fmt.Println("Los primos encontrados son")
	fmt.Println(newPrimes)
	resp(conn,newPrimes)
}

func main() {
	//prot := flag.String("prot","","El protocolo para transmitir")
	//iphost := flag.String("iphost","","La ip del host")
	port := flag.String("port","","El puerto de escucha")
	
	/*Usando el protocolo CONN_TYPE a la dirección/es CONN_HOST por medio del puerto CONN_PORT*/
	flag.Parse()
	fmt.Println("En espera por el puerto ", *port)
	//listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	//listener, err := net.Listen(*prot, (*iphost) +":"+ (*port))
	//listener, err := net.Listen("tcp", "192.168.1.0:2000")
	//listener, err := net.Listen("tcp", (*iphost) +":"+ (*port))
	//listener, err := net.Listen("tcp", ":30000")
	/*var port string
	fmt.Print("Elige el puerto= ")
	fmt.Scanln(&port)
	fmt.Print("En espera\n")*/
	listener, err := net.Listen("tcp", ":"+ *port)
	/*Comprueba si se ha producido un error en la conexión*/
	checkError(err)
	
	/*Se pone en espera para decibir una llamada y devuelve una conexion genérica conn*/
	conn, err := listener.Accept()
	fmt.Println("Cliente escuchado\n")
	/*Nos aseguramos que cuando se acabe la tarea, se cierre la conexión entre el emisor y el receptor*/
	defer conn.Close()
	checkError(err)
	
	handle(conn) //go quitado
	fmt.Println("Servidor finalizado ")
    // TO DO
}
