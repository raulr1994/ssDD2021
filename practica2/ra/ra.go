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
    "ms"
    "sync"
)

type Request struct{
    Clock   int
    Pid     int
}

type Reply struct{}

type RASharedDB struct {
    OurSeqNum   int //Nuestro reloj
    HigSeqNum   int //Reloj más alto
    OutRepCnt   int //Nº de reply necesarios para salir de la SC
    ReqCS       boolean //True para querer acceder a la seccion critica
    RepDefd     int[] 
    ms          *MessageSystem
    done        chan bool
    chrep       chan bool //TOdos nodos que aun estan esperando permiso para entra en la SC
    Mutex       sync.Mutex // mutex para proteger concurrencia sobre las variables
    // TODO: completar
    me 	int //Id del nodo
    NNodes 	int //Numero total de nodos conectados
    Escritor	bool //El nodo es escritor(True) o lector (False)
}

//Variables globales
var mensPet = "REQUEST"
var mensPer = "PERMISION"

func New(me int, usersFile string, nnodes int, escritor bool) (*RASharedDB) {
    messageTypes := []Message{Request, Reply}
    msgs = ms.New(me, usersFile string, messageTypes)
    ra := RASharedDB{0, 0, 0, false, []int{}, &msgs,  make(chan bool),  make(chan bool), &sync.Mutex{}, me, nnodes, escritor}
    // TODO completar
    return &ra
}

//Pre: Verdad
//Post: Realiza  el  PreProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PreProtocol(){
    // TODO completar
    	ra.Mutex.Lock()
    	ra.ReqCS = true
    	ra.OurSeqNum = HigSeqNum + 1
    	OutRepCnt := ra.NNodes-1
    	
    	for j:= 0; j < NNodes; j++ {
    		if(j != ra.me){ //Enviar una peticion para acceder a la SC a todos los nodos salvo a mi mismo
    			ra.ms.Send(j,Message{Operation: mensPet, Sender: ra.me,Clock: ra.OurSeqNum, Write: ra.Escritor})
    		}
    	}
    	
    	ra.Mutex.Unlock()
    	
    	for OutRepCnt > 0 { //Esperar a que haya recibido todas las respuestas
    	}
    	
    	ra.Mutex.Lock()
    	
    	ra.Mutex.Lock()
}

//Pre: Verdad
//Post: Realiza  el  PostProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PostProtocol(){
	// TODO completa
    	ra.Mutex.Lock()
    	ReqCS = false
    	
    	for d := 0, d < ra.NNodes, d++ { //Liberar a todos los que estaban a la espera de permiso
    		if(chrep[d]){
    			chrep[d] = false;
    			ra.ms.Send(d,Message{Operation: mensPer, Sender: ra.me})
    		}
    	}
    	ra.Mutex.Lock()
}

func (ra *RASharedDB) Stop(){
    ra.ms.Stop()
    ra.done <- true
}
