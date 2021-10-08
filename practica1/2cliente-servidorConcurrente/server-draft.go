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
	"./com"
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
	vectDirPort := lecturaFichero("./ipServer.txt")
	fmt.Println(vectDirPort)
	ip, puerto := obtenerIPPuerto(vectDirPort,0)
	fmt.Println("La IP es ", ip)
	fmt.Println("El puerto es ", puerto)

	listener, err := net.Listen("tcp", ip+":"+puerto)
	checkError(err)
	fin := false
	for !fin {
		conn, err := listener.Accept()
		defer conn.Close()
		checkError(err)
		handle(conn)
	}
	fmt.Println("Servidor finalizado ")
}

