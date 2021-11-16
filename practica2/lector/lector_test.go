/*
* AUTOR: Raúl Rustarazo Carmona
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* NIA: 715657
* FECHA: noviembre de 2021
* FICHERO: lector_test.go
* DESCRIPCIÓN: Implementación de un lector para el algoritmo de ricart-agrawala
*/
package lector

import (
	//"flag"
	"github.com/DistributedClocks/GoVector/govec"
	"testing"
	"example.gestorfichero/gestorfichero"
	"strconv"
	"example.com/RA"
)


func TestReader(t *testing.T){
    	nExp := 3 //Nº total de pruebas de lectura o escritura
    	opP := false //Leer(False),Escibir(True)
    	nNodes := 2 //Nº total de nodos en la red
    	rai := ra.New(2,"./users.txt",nNodes,opP);//Pid, nombre fichero, numero de nodos y si es escritor o lector
	if nNodes > 1{
    		go rai.Listening() //Lanzar el proceso de escucha de peticiones
    	}
    	for i := 0; (i < nExp); i++ {
		rai.Acquire_mutex()//Pedir adquirir el mutex
		//Acceso a la CS para escibir o leer
		rai.Logger.LogLocalEvent("Entering to SC by Lector " + strconv.Itoa(rai.MiId()) + " and reading", govec.GetDefaultLogOptions())
		
		rai.Mutex.Lock()
		lectura := gestorfichero.LeerFichero("memory.txt")
		rai.Mutex.Unlock()
		
		rai.Logger.LogLocalEvent(gestorfichero.LeerLectura(lectura), govec.GetDefaultLogOptions())
		//Enviar al resto de nodos que ha escrito un escritor
		rai.Logger.LogLocalEvent("Exiting from SC by Lector" + strconv.Itoa(rai.MiId()) + "and free the SC " + strconv.Itoa(rai.MiId()), govec.GetDefaultLogOptions())
		rai.Release_mutex("",false)//Liberar el mutex
		//Dormir un rato antes de hacer la siguiente peticion
	}
	rai.WaitListening()
	rai.Logger.LogLocalEvent("Finishing the Lector " + strconv.Itoa(rai.MiId()), govec.GetDefaultLogOptions())
	rai.Stop()
}
