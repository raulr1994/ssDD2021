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
	/*ssh, err := sshcode.NewSshClient(
		"a715657",
		"155.210.154.205",
		22,
		"/root/.ssh/id_rsa",//"/Users/some-user/.ssh/id_rsa",
		"")*/
	ssh, err := sshcode.NewSshClient(
		nombreusuario,
		ipjhost,
		22,
		"./rsa",//"/root/.ssh/id_rsa",//"/Users/some-user/.ssh/id_rsa",
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

func main(){
	//nombreusuario := 
	var nombreusuario string
	var ipjhost string
	fmt.Print("Introduce el nombre del usuario del host = ")
    	fmt.Scanln(&nombreusuario)
    	fmt.Print("Introduce la ip del host = ")
    	fmt.Scanln(&ipjhost)
    	var eleccion int
    	fmt.Print("Elije que ejecutar \n -Descargar el cliente(1) \n -Descargar el servidor(2) \n -Ejecutar el cliente(3) \n -Ejecutar el servidor(4) = ")
    	fmt.Scanln(&eleccion)
    	
 	switch eleccion{
 		case 1:
 			executeComand(nombreusuario,ipjhost,"rm -f ./client")
 			executeComand(nombreusuario,ipjhost,"rm -f ./ipClient.txt")
 			executeComand(nombreusuario,ipjhost,"wget https://raw.githubusercontent.com/raulr1994/ssDD2021/main/practica1/1cliente-servidor/Client/client")
 			executeComand(nombreusuario,ipjhost,"wget https://raw.githubusercontent.com/raulr1994/ssDD2021/main/practica1/1cliente-servidor/Client/ipClient.txt")
 			executeComand(nombreusuario,ipjhost,"chmod 777 ./client")
 		case 2:
 			executeComand(nombreusuario,ipjhost,"rm -f ./server-draft")
 			executeComand(nombreusuario,ipjhost,"rm -f ./ipServer.txt")
 			executeComand(nombreusuario,ipjhost,"wget https://raw.githubusercontent.com/raulr1994/ssDD2021/main/practica1/1cliente-servidor/Server/server-draft")
 			executeComand(nombreusuario,ipjhost,"wget https://raw.githubusercontent.com/raulr1994/ssDD2021/main/practica1/1cliente-servidor/Server/ipServer.txt")
 			executeComand(nombreusuario,ipjhost,"chmod 777 ./server-draft")
 		case 3:
 			executeComand(nombreusuario,ipjhost,"./client")
 		case 4:

 			executeComand(nombreusuario,ipjhost,"./server-draft")
 		case 5:
 			//executeComand(nombreusuario,ipjhost,"./hello -port 23")
 			executeComand(nombreusuario,ipjhost,"ls")
 			//executeComand(nombreusuario,ipjhost,"ssh raulrcuni@192.168.1.2 ./server -port 30000")
 		default:
 			fmt.Println("La eleccion no existe")
 	}
}
