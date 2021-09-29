/*
* AUTOR: Rafael Tolosana Calasanz
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2021
* FICHERO: client.go
* DESCRIPCIÓN: cliente completo para los cuatro escenarios de la práctica 1
*/
package main

import (
    "fmt"
    //"time"
    "os"
    "net"
    //"./com"
    //"/home/raulrcuni/goSD/trabajo1/com"
    "example.com/com"
    //"flag"
    "encoding/binary"
    //"bytes"
    "bufio"
	"log"
	"strings"
	"strconv"
)

type Message struct {
	ID string
	VectP int
}


func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func enviarPeticion(conn net.Conn, desde int, hasta int) {

	fmt.Println("Enviar peticion de primos")
	
	interval := []int{desde,hasta}
	
	bs := make([]byte,4)
	intervalSend := make([]byte, 2*4)
	j:= 0
	for i:= range interval{
		binary.LittleEndian.PutUint32(bs, uint32(interval[i]))
		for i:= range bs {
			intervalSend[j + i] = bs[i]
		}
		j = j +4
	}
	conn.Write([]byte(intervalSend))
}

func recibirRespuesta(conn net.Conn) {
	bs := make([]byte, 4)
	//n := 1
	buf := make([]byte, 4)
	_, err := conn.Read(buf)
	checkError(err)
	nPrimes := binary.LittleEndian.Uint32(buf)
	fmt.Println("El numero de primos recibidos es= ", nPrimes)
	
	fmt.Println("Enviar OK al servidor")
	binary.LittleEndian.PutUint32(bs, 1)
	conn.Write([]byte(bs))	
	
	/*n = int(nPrimes)
	fmt.Println("N ahora vale ", n)*/
	
	newVectPrimes := [] uint32{}
	t2 := make([]byte, nPrimes*4)
	_, err2 := conn.Read(t2)
	checkError(err2)

	fmt.Println("traduciendo los datos")
	j := 0
	for i:= 0; i < int(nPrimes); i++ {
		for z:= 0; z < 4; z++{
			bs[z] = t2[z+j]
		}
		newVectPrimes = append(newVectPrimes, binary.LittleEndian.Uint32(bs))
		j = j + 4
	}
	fmt.Println(newVectPrimes)
}


func obtenerIPPuerto(vectDirPort [] string, pos int) (ip string, puerto string, desde int, hasta int){
	s := strings.Split(vectDirPort[pos],":")
	ip = s[0] //La ip
	puerto = s[1] //El puerto
	desde, err := strconv.Atoi(s[2])
	fmt.Println(err)
	hasta , err = strconv.Atoi(s[3])
	fmt.Println(err)
	return ip, puerto, desde, hasta
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

func main(){
	vectDirPort := lecturaFichero("./ipClient.txt")
	fmt.Println(vectDirPort)
	ip, puerto, desde, hasta := obtenerIPPuerto(vectDirPort,0)
	fmt.Println("La IP es ", ip)
	fmt.Println("El puerto es ", puerto)
	
	fmt.Println("Desde = ", desde)
	fmt.Println("Hasta = ", hasta)
	
	interval := com.TPInterval{desde, hasta}
    	fmt.Println("Intervalo= ", interval)
    	
    	fmt.Print("Intentado establecer conexión\n")
    	tcpAddr, err := net.ResolveTCPAddr("tcp", ip + ":" + puerto)
	checkError(err)
	fmt.Print("Conexion extablecida\n")

	fmt.Print("Enviando peticion conexion\n")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	
	defer conn.Close()
	
	fmt.Print("Enviando datos\n")
	enviarPeticion(conn,desde,hasta)
	fmt.Print("Datos enviados\n")
		
	fmt.Print("Recibiendo respuesta\n")	
	recibirRespuesta(conn)
	fmt.Print("Respuesta recibida\n")
}
