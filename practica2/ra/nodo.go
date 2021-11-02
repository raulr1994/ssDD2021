/*
* AUTOR: Rafael Tolosana Calasanz
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2021
* FICHERO: ricart-agrawala.go
* DESCRIPCIÓN: Implementación del algoritmo de Ricart-Agrawala Generalizado en Go
*/
package ra

import (
	"testing"
)

func acquire_mutex(tipoOp bool){ //Escibir true, leer false

}

func release_mutex(){

}

func main(){
    	//1º Recoger parametros
    	//2º inicializar variables
    	//3º Ejecutar Richard agravala
    	nExp := 1
    	rai := New();//Pid y nombre fichero
    	opP := true //Leer(False),Escibir(True)
    	bufferDeLectura := "" //Donde guardar lo que se lee
    	bufferDe Escritura := "" //Que escribir
    	for i := 0; (i < nExp); i++ {
		acquire_mutex(opP)//Pedir adquirir el mutex
		//Acceso a la CS para escibir o leer
		release_mutex()//Liberar el mutex
	}
}
