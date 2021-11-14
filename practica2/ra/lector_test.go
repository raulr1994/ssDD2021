/*
* AUTOR: Raúl Rustarazo Carmona
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* NIA: 715657
* FECHA: noviembre de 2021
* FICHERO: lector_test.go
* DESCRIPCIÓN: Implementación de un lector para el algoritmo de ricart-agrawala
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
    	rai := New(2,"./users.txt",nNodes,opP);//Pid, nombre fichero, numero de nodos y si es escritor o lector
    	gestorfichero.LimpiarFichero("memory.txt")
	
    	go rai.listening() //Lanzar el proceso de escucha de peticiones
    	for i := 0; (i < nExp); i++ {
		rai.acquire_mutex(opP)//Pedir adquirir el mutex
		rai.Logger.LogLocalEvent("Entering to SC by Lector " + strconv.Itoa(rai.me), govec.GetDefaultLogOptions())
		//Acceso a la CS para escibir o leer

		rai.Logger.LogLocalEvent("Reading by Lector " + strconv.Itoa(rai.me), govec.GetDefaultLogOptions())
		
		rai.Mutex.Lock()
		lectura := gestorfichero.LeerFichero("memory.txt")
		rai.Mutex.Unlock()
		
		rai.Logger.LogLocalEvent(gestorfichero.LeerLectura(lectura), govec.GetDefaultLogOptions())
		//Enviar al resto de nodos que ha escrito un escritor
		rai.Logger.LogLocalEvent("Exiting from SC by Lector " + strconv.Itoa(rai.me), govec.GetDefaultLogOptions())
		rai.release_mutex("",false)//Liberar el mutex
		//Dormir un rato antes de hacer la siguiente peticion
	}
	rai.waitListening()
	rai.Logger.LogLocalEvent("Finishing the Lector " + strconv.Itoa(rai.me), govec.GetDefaultLogOptions())
	rai.ms.Stop()
	/*for{
		//Para no dejar ningun nodo desatendido que le falte alguna peticion
	}*/
}
