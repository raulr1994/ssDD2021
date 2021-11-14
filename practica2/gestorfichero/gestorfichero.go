/*
* AUTOR: Raúl Rustarazo Carmona
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* NIA: 715657
* FECHA: noviembre de 2021
* FICHERO: ms.go
* DESCRIPCIÓN: Implementación de un sistema de gestion de lectura y escritura de ficheros
*/
package gestorfichero

import (
	"bufio"
	"fmt"
	"os"
	"log"
	"io/ioutil"
)

/*
	Lee el contenido de un fichero y devuelve un array de string con las lineas leidas del fichero
*/
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

/*
	Añade al fichero una linea nueva con un nuevo string
*/
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


/*
	Borra el contenido del fichero
*/
func LimpiarFichero(nameFile string){
	b := []byte("")
	err := ioutil.WriteFile(nameFile, b ,0644)
	if err !=nil {
		log.Fatal(err)
	}
}

