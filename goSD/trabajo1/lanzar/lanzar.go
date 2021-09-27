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
		nombreusuario,
		ipjhost,
		22,
		"/root/.ssh/id_rsa",//"/Users/some-user/.ssh/id_rsa",
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
	var nombreusuario string
	var ipjhost string
	var comand string
	fmt.Print("Introduce el nombre del usuario del host = ")
    	fmt.Scanln(&nombreusuario)
    	fmt.Print("Introduce la ip del host = ")
    	fmt.Scanln(&ipjhost)
    	var eleccion int
    	fmt.Print("Elije que ejecutar el cliente(1) o el servidor(2) = ")
    	fmt.Scanln(&eleccion)
    	
 	switch eleccion{
 		case 1:
 			executeComand(nombreusuario,ipjhost,"//////")
 		case 2:
 			executeComand(nombreusuario,ipjhost,"//////")
 		default:
 			fmt.Println("La eleccion no existe")
 	}
 	//executeComand(nombreusuario,ipjhost,comand)
    	/*fmt.Print("Elije el comando a ejecutar = ")
    	fmt.Scanln(&comand)*/
    	
	/*ssh, err := sshcode.NewSshClient(
		"raulrcuni",
		"192.168.1.2",
		22,
		"/root/.ssh/id_rsa",//"/Users/some-user/.ssh/id_rsa",
		"")

	if err != nil {
		log.Printf("SSH init error %v", err)
	} else {
		output, err := ssh.RunCommand("ls")
		fmt.Println(output)
		if err != nil {
			log.Printf("SSH run command error %v", err)
		}
	}*/
}
