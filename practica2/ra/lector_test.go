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
	//"flag"
	"github.com/DistributedClocks/GoVector/govec"
	"testing"
	"example.gestorfichero/gestorfichero"
	"strconv"
)


func TestReader(t *testing.T){
    	nExp := 1 //Nº total de pruebas de lectura o escritura
    	opP := false //Leer(False),Escibir(True)
    	nNodes := 2 //Nº total de nodos en la red
    	rai := New(1,"./users.txt",nNodes,opP);//Pid, nombre fichero, numero de nodos y si es escritor o lector
    	gestorfichero.LimpiarFichero("memory.txt")
	
    	go rai.listening() //Lanzar el proceso de escucha de peticiones
    	//3º Ejecutar Richard agravala
    	for i := 0; (i < nExp); i++ {
		rai.acquire_mutex(opP)//Pedir adquirir el mutex
		rai.Logger.LogLocalEvent("Entering to SC", govec.GetDefaultLogOptions())
		//Acceso a la CS para escibir o leer

		rai.Logger.LogLocalEvent("Reading by " + strconv.Itoa(rai.me), govec.GetDefaultLogOptions())
		lectura := gestorfichero.LeerFichero("memory.txt")
		gestorfichero.MostrarLectura(lectura)
		
		//Enviar al resto de nodos que ha escrito un escritor
		rai.Logger.LogLocalEvent("Exiting from SC", govec.GetDefaultLogOptions())
		rai.release_mutex("",false)//Liberar el mutex
		//Dormir un rato antes de hacer la siguiente peticion
	}
	rai.Logger.LogLocalEvent("Finishing " + strconv.Itoa(rai.me), govec.GetDefaultLogOptions())
	for{
		//Para no dejar ningun nodo desatendido que le falte alguna peticion
	}
}
