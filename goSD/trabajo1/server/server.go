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
	//"encoding/gob"
	"fmt"
	"net"
	"os"
	"io"
	//"./com"
	"example.com/com"
	//"flag"
	//"bytes"
	//"time"
	"encoding/binary"
	"bufio"
	"log"
	"strings"
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
	fmt.Println("Leyendo intervalo")
	//for{
		/*buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if logerr(err){
			break
		}*/
		/*fmt.Println("Leido el buffer")
		data := binary.LittleEndian.Uint32(buf)
		fmt.Println(data)*/
		var intervalV[2] uint32
		t2 := make([]byte, 2*4)
		
		bufferPrimesB := make([]byte, 2*4)
		_, err := conn.Read(bufferPrimesB)
		checkError(err)
		
		t2 = bufferPrimesB
		bs := make([]byte, 4)
		fmt.Println("traduciendo los datos")
		j := 0
		for i := 0; i < 2; i++ {
			for z := 0; z < 4; z++ {
				bs[z] = t2[z+j]
			}
			intervalV[i] = binary.LittleEndian.Uint32(bs)
			j = j + 4
		}
		
		/*n := int64(data)
		fmt.Println("Hasta= ",int(n))*/
		fmt.Println("El intervalo pedido es ", intervalV)
		newInterval := com.TPInterval{A:int(intervalV[0]),B:int(intervalV[1])}
		interval = newInterval
	//}
	return interval
}

func resp(conn net.Conn, listPrimes []int){
	bs := make([]byte,4)
	binary.LittleEndian.PutUint32(bs,uint32(len(listPrimes)))
	fmt.Println("Enviando el numero de primos encontrados ", len(listPrimes))
	conn.Write([]byte(bs))
	
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	checkError(err)
	fmt.Println("Recibido OK del cliente")
	data := binary.LittleEndian.Uint32(buf)
	fmt.Println(data)
	
	vectorPrimes := make([]byte, len(listPrimes)*4)
	j := 0
	for i := range listPrimes {
		binary.LittleEndian.PutUint32(bs, uint32(listPrimes[i]))
		for i:= range bs {
			vectorPrimes[j + i] = bs[i]
		}
		j = j + 4
	}
	conn.Write([]byte(vectorPrimes))
	
	conn.Close()
}

func handle(conn net.Conn){
	//defer conn.Close()
	/*timeoutDuration := 2*time.Second
	fmt.Println("Iniciando servidor")
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))*/
	
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Cliente conectado desde " + remoteAddr)
	//conn.Close()
	interval := read(conn)
	fmt.Println(interval)
	newPrimes := FindPrimes(interval)
	fmt.Println("El numero de primos encontrados es ", len(newPrimes))
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

func main() {
	//port := flag.String("port","","El puerto de escucha")
	vectDirPort := lecturaFichero("./ipServer.txt")
	fmt.Println(vectDirPort)
	ip, puerto := obtenerIPPuerto(vectDirPort,0)
	fmt.Println("La IP es ", ip)
	fmt.Println("El puerto es ", puerto)
	
	/*Usando el protocolo CONN_TYPE a la dirección/es CONN_HOST por medio del puerto CONN_PORT*/
	/*flag.Parse()*/
	fmt.Println("En espera por el puerto ", puerto)

	listener, err := net.Listen("tcp", ":"+ puerto)
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
}
