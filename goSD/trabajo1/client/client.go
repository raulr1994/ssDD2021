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
    "encoding/gob"
    "bytes"
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

func send(conn net.Conn, intervalNew com.TPInterval) {
	
	MessageS := com.Request{Id: 1, Interval: intervalNew}
	//MessageS := new(Request)
	bin_buf := new(bytes.Buffer)
	
	gobobj := gob.NewEncoder(bin_buf)
	
	gobobj.Encode(MessageS)
	
	conn.Write(bin_buf.Bytes())
}

func recv(conn net.Conn) {
	tmp := make([]byte, 500)
	conn.Read(tmp)
	
	tmpbuff := bytes.NewBuffer(tmp)
	tmpstruct := new(com.Reply)
	
	gobobj := gob.NewDecoder(tmpbuff)
	gobobj.Decode(tmpstruct)
	
	fmt.Println(tmpstruct)	
}

func main(){
	/*desde := flag.Int("desde",0,"primer valor del rango")
	hasta := flag.Int("hasta",0,"ultimo valor del rango")
	IPPort := flag.String("IPPort","","La IP destino y el puerto")
	
	flag.Parse()
	fmt.Println("Desde = ", *desde)
	fmt.Println("Hasta = ", *hasta)
	fmt.Println("Destino = ", *IPPort)
	
	interval := com.TPInterval{*desde, *hasta}
    	fmt.Println("Intervalo= ", interval)
    	
    	fmt.Print("Intentado establecer conexión\n")
    	tcpAddr, err := net.ResolveTCPAddr("tcp", *IPPort)
    	
    	fmt.Println("Desde = ", *desde)
	fmt.Println("Hasta = ", *hasta)
	fmt.Println("Destino = ", *IPPort)
	
	interval := com.TPInterval{*desde, *hasta}
    	fmt.Println("Intervalo= ", interval)
    	
    	fmt.Print("Intentado establecer conexión\n")
    	tcpAddr, err := net.ResolveTCPAddr("tcp", *IPPort)*/
    	
    	var desde int
    	fmt.Print("Elige el rango de primos\n")
    	fmt.Print("Desde = ")
    	fmt.Scanln(&desde)
    	var hasta int
    	fmt.Print("Hasta = ")
    	fmt.Scanln(&hasta)
    	var ip string
    	fmt.Print("Elige la ip del host = ")
    	fmt.Scanln(&ip)
    	var port string
    	fmt.Print("Elige el puerto del host = ")
    	fmt.Scanln(&port)
    	
    	fmt.Println("Desde = ", desde)
	fmt.Println("Hasta = ", hasta)
	fmt.Println("Destino = ", ip + ":" + port)
	
	interval := com.TPInterval{desde, hasta}
    	fmt.Println("Intervalo= ", interval)
    	
    	fmt.Print("Intentado establecer conexión\n")
    	tcpAddr, err := net.ResolveTCPAddr("tcp", ip + ":" + port)

    	
	checkError(err)
	fmt.Print("Conexion extablecida\n")

	fmt.Print("Enviando datos\n")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	
	defer conn.Close()
	
	fmt.Print("Enviando datos\n")
	send(conn,interval)
	fmt.Print("Datos enviados\n")
		
	fmt.Print("Recibiendo respuesta\n")	
	recv(conn)
	fmt.Print("Respuesta recibida\n")
	
	/*Dirección IP y puerto por el que nos conectamos al servidor*/
    //endpoint := "155.210.154.200:30000"
    
    // TODO: crear el intervalo solicitando dos números por teclado
	/*TPInterval es un objeto definido en la propia conexión para asignar el rango pedido de 1000 a 7000(es un vector{primer numero,ultimo numero})*/
    /*interval := com.TPInterval{1000, 70000}*/
	
	/*Para establecer la configruación de la conexión al determinado servidor con cierta IP y PUERTO*/
    //tcpAddr, err := net.ResolveTCPAddr("tcp", endpoint)
    //checkError(err)

    //conn, err := net.DialTCP("tcp", nil, tcpAddr)
    //checkError(err)

    // la variable conn es de tipo *net.TCPconn
}
