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
    "example.ms/ms"
    "sync"
    "github.com/DistributedClocks/GoVector/govec"
    //"time"
)

/*type Message interface {}

type Request struct {
	Type string
	Id int
	Clock int
	Escritor bool //El mensaje es de un escritor(True) o lector (False)
}

type Reply struct{
	Type string
	Response string
}*/

const (
	TYPEREQUEST = "REQUEST"
	TYPEREPLY = "REPLY"
	MSGFREE = "LIBERAR"
)

type RASharedDB struct {
    OurSeqNum   int //Nuestro reloj
    HigSeqNum   int //Reloj más alto
    OutRepCnt   int //Nº de reply necesarios para salir de la SC
    ReqCS       bool //True para querer acceder a la seccion critica
    RepDefd     []bool //Todos nodos que aun estan esperando permiso para entra en la SC
    ms          *ms.MessageSystem
    done        chan bool //Son canales
    chrep       chan bool 
    Mutex       sync.Mutex // mutex para proteger concurrencia sobre las variables
    // TODO: completar
    me 	int //Id del nodo
    NNodes 	int //Numero total de nodos conectados
    Escritor	bool //El nodo es escritor(True) o lector (False)
    Logger	*govec.GoLog
}


func New(me int, usersFile string, nnodes int, escritor bool) (*RASharedDB) {
    messageTypes := []ms.Message{ms.Request{},ms.Reply{}}
    Logger := govec.InitGoVector("client", "LogFileEventInt", govec.GetDefaultConfig())
    msgs := ms.New(me, usersFile, messageTypes)
    ra := RASharedDB{0, 0, 0, false, []bool{}, &msgs,  make(chan bool),  make(chan bool), sync.Mutex{}, me, nnodes, escritor, Logger}
    // TODO completar
    return &ra
}

//Pre: Verdad
//Post: Realiza  el  PreProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PreProtocol(){
    // TODO completar
    	ra.Mutex.Lock()
    	ra.ReqCS = true //Pedir acceso a la sección critica
    	ra.OurSeqNum = ra.HigSeqNum + 1 //Actualizar el reloj interno
    	ra.OutRepCnt = ra.NNodes-1 //Nodos a los que se esta esperando respuesta
    	
    	for j:= 0; j < ra.NNodes; j++ {
    		if(j != ra.me){ //Enviar una peticion para acceder a la SC a todos los nodos salvo a si mismo
    			ra.ms.Send(j+1, ms.Request{TYPEREQUEST,ra.me,ra.OurSeqNum,ra.Escritor})//Operacion/IdNodo/ClockNodo/Escritor/Lector
    		}
    	}
    	
    	for {
    		ra.Mutex.Lock()//Asegurar la exlusión mutua cuabndo se consulta la variable
    		if ra.OutRepCnt == 0 {//Esperar a que haya recibido todas las respuestas
    			break
    		}
    		ra.Mutex.Unlock()    		
    	}
}
//Pre: Verdad
//Post: Realiza  el  PostProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PostProtocol(){
	// TODO completa
    	ra.Mutex.Lock()
    	ra.ReqCS = false
    	
    	for d := 0; d<ra.NNodes ;d++ { //Liberar a todos los que estaban a la espera de permiso(True)
    		if(ra.RepDefd[d]){ //Si estaba a la espera, hay que liberarlo
    			ra.RepDefd[d] = false;
    			ra.ms.Send(d,ms.Reply{TYPEREPLY,MSGFREE})
    		}
    	}
    	ra.Mutex.Lock()
}

func (ra *RASharedDB) Stop(){
    ra.ms.Stop()
    ra.done <- true
}
