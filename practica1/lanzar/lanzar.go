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
	"log"
	"fmt"
    	"example.ssh/ssh"
)


func executeComand (nombreusuario string, ipjhost string, comand string) {
	ssh, err := sshcode.NewSshClient(
		"a715657",
		"155.210.154.205",
		22,
		"/root/.ssh/id_rsa",//"/Users/some-user/.ssh/id_rsa",
		"bunkma29")

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

func main(){
	//nombreusuario := 
	var nombreusuario string
	var ipjhost string
	fmt.Print("Introduce el nombre del usuario del host = ")
    	fmt.Scanln(&nombreusuario)
    	fmt.Print("Introduce la ip del host = ")
    	fmt.Scanln(&ipjhost)
    	var eleccion int
    	fmt.Print("Elije que ejecutar \n -Descargar el cliente(1) \n -Descargar el servidor secuencial(2) \n  -Descargar el servidorConcurrente V1(3) \n -Descargar el servidorCOncurrenteV2(4) \n -Descargar el master(5) = \n -Descargar el worker(6) = \n -Iniciar el cliente(7) = \n -Iniciar el servidor secuencial/concurrV1/concurrV2(8) = \n -Iniciar el master(9) = \n -Iniciar el worker(10) = \n")
    	fmt.Scanln(&eleccion)
    	
 	switch eleccion{
 		case 1:
 			executeComand(nombreusuario,ipjhost,"rm -f ./client")
 			executeComand(nombreusuario,ipjhost,"rm -f ./ipClient.txt")
 			executeComand(nombreusuario,ipjhost,"wget https://raw.githubusercontent.com/raulr1994/ssDD2021/main/practica1/Client/client")
 			executeComand(nombreusuario,ipjhost,"wget https://raw.githubusercontent.com/raulr1994/ssDD2021/main/practica1/Client/ipClient.txt")
 			executeComand(nombreusuario,ipjhost,"chmod 777 ./client")
 		case 2:
 			executeComand(nombreusuario,ipjhost,"rm -f ./server-draft")
 			executeComand(nombreusuario,ipjhost,"rm -f ./ipServer.txt")
 			executeComand(nombreusuario,ipjhost,"wget https://raw.githubusercontent.com/raulr1994/ssDD2021/main/practica1/1cliente-servidor/Server/server-draft")
 			executeComand(nombreusuario,ipjhost,"wget https://raw.githubusercontent.com/raulr1994/ssDD2021/main/practica1/1cliente-servidor/Server/ipServer.txt")
 			executeComand(nombreusuario,ipjhost,"chmod 777 ./server-draft")
 		case 3:
 			executeComand(nombreusuario,ipjhost,"rm -f ./server-draft")
 			executeComand(nombreusuario,ipjhost,"rm -f ./ipServer.txt")
 			executeComand(nombreusuario,ipjhost,"wget https://raw.githubusercontent.com/raulr1994/ssDD2021/main/practica1/2cliente-servidorConcurrente/Server/server-draft")
 			executeComand(nombreusuario,ipjhost,"wget https://raw.githubusercontent.com/raulr1994/ssDD2021/main/practica1/2cliente-servidorConcurrente/Server/ipServer.txt")
 			executeComand(nombreusuario,ipjhost,"chmod 777 ./server-draft")
 		case 4:
 			executeComand(nombreusuario,ipjhost,"rm -f ./server-draft")
 			executeComand(nombreusuario,ipjhost,"rm -f ./ipServer.txt")
 			executeComand(nombreusuario,ipjhost,"wget https://raw.githubusercontent.com/raulr1994/ssDD2021/main/practica1/3cliente-servidorConcurrente2/Server/server-draft")
 			executeComand(nombreusuario,ipjhost,"wget https://raw.githubusercontent.com/raulr1994/ssDD2021/main/practica1/3cliente-servidorConcurrente2/Server/ipServer.txt")
 			executeComand(nombreusuario,ipjhost,"chmod 777 ./server-draft")
 		case 5:
 			executeComand(nombreusuario,ipjhost,"rm -f ./master")
 			executeComand(nombreusuario,ipjhost,"rm -f ./ipServer.txt")
 			executeComand(nombreusuario,ipjhost,"rm -f ./ipServer.txt")
 			executeComand(nombreusuario,ipjhost,"wget https://raw.githubusercontent.com/raulr1994/ssDD2021/main/practica1/4master-worker/Master/master")
 			executeComand(nombreusuario,ipjhost,"wget https://raw.githubusercontent.com/raulr1994/ssDD2021/main/practica1/4master-worker/Master/ipServer.txt")
 			executeComand(nombreusuario,ipjhost,"wget https://raw.githubusercontent.com/raulr1994/ssDD2021/main/practica1/4master-worker/Master/ipWorkers.txt")
 			executeComand(nombreusuario,ipjhost,"chmod 777 ./master")
 		case 6:
 			executeComand(nombreusuario,ipjhost,"rm -f ./worker")
 			executeComand(nombreusuario,ipjhost,"rm -f ./ipServer.txt")
 			executeComand(nombreusuario,ipjhost,"rm -f ./ipServer.txt")
 			executeComand(nombreusuario,ipjhost,"wget https://raw.githubusercontent.com/raulr1994/ssDD2021/main/practica1/4master-worker/Worker/worker")
 			executeComand(nombreusuario,ipjhost,"wget https://raw.githubusercontent.com/raulr1994/ssDD2021/main/practica1/4master-worker/Worker/ipWorker.txt")
 			executeComand(nombreusuario,ipjhost,"chmod 777 ./worker")
 		case 7:
 			executeComand(nombreusuario,ipjhost,"./client")
 		case 8:
 			executeComand(nombreusuario,ipjhost,"./server-draft")
 		case 9:
 			executeComand(nombreusuario,ipjhost,"./master")
 		case 10:
 			executeComand(nombreusuario,ipjhost,"./worker")
 		case 11:
 			executeComand(nombreusuario,ipjhost,"ls")
 		default:
 			fmt.Println("La eleccion no existe")
 	}
}
