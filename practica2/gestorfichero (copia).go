/*
* AUTOR: Rafael Tolosana Calasanz
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2021
* FICHERO: ms.go
* DESCRIPCIÓN: Implementación de un sistema de mensajería asíncrono, insipirado en el Modelo Actor
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"log"
	"io/ioutil"
)

func LeerFichero(nameFile string) (lectura [] string){
	file, err := os.Open(nameFile)
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan(){
		//fmt.Println(fileScanner.Text())
		lectura = append(lectura,fileScanner.Text())
	}
	
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}
	
	file.Close()
	return lectura
}

func lecturaFichero(nameFile string) (vectDirPort [] string, nworkers int){
	file, err := os.Open(nameFile)
	
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}
	
	fileScanner := bufio.NewScanner(file)
	
	//vectDirPort = [] string{}
	nworkers = 0
	for fileScanner.Scan(){
		//fmt.Println(fileScanner.Text())
		vectDirPort = append(vectDirPort,fileScanner.Text())
		nworkers++
	}
	
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}
	
	file.Close()
	return vectDirPort,nworkers
}

func EscribirFichero(nameFile string, linea string){	
	file, err := os.OpenFile(nameFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}
	
	_,err = file.Write([]byte("\n"+linea))
	if err !=nil {
		log.Fatal(err)
	}
}

func limpiarFichero(nameFile string){
	b := []byte("")
	err := ioutil.WriteFile(nameFile, b ,0644)
	if err !=nil {
		log.Fatal(err)
	}
}

func mostrarLectura(lectura [] string) {
	for linea := range lectura {
		fmt.Println(lectura[linea])
	}
}

func main(){
	limpiarFichero("memory.txt")
	lectura := LeerFichero("memory.txt")
	mostrarLectura(lectura)
	EscribirFichero("memory.txt","Hola1")
	lectura = LeerFichero("memory.txt")
	mostrarLectura(lectura)
	EscribirFichero("memory.txt","Hola2")
	lectura = LeerFichero("memory.txt")
	mostrarLectura(lectura)
}
