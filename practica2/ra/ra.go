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
    "fmt"
    "strconv"
    "reflect"
    "example.gestorfichero/gestorfichero"
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
	TYPEWRITE = "WRITE"
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
    var Logger *govec.GoLog
    if escritor {
    	Logger = govec.InitGoVector("Escritor_" + strconv.Itoa(me) , "LogFileEventInt from escritor " +  strconv.Itoa(me), govec.GetDefaultConfig())
    } else {
    	Logger = govec.InitGoVector("Lector_" + strconv.Itoa(me) , "LogFileEventInt from lector " + strconv.Itoa(me), govec.GetDefaultConfig())
    }
    msgs := ms.New(me, usersFile, messageTypes, Logger)
    ra := RASharedDB{0, 0, 0, false, []bool{}, &msgs,  make(chan bool),  make(chan bool), sync.Mutex{}, me, nnodes, escritor, Logger}
    
    for i:=0; i<nnodes; i++ {
    	ra.RepDefd = append(ra.RepDefd,false)
    }
    //fmt.Println("Vector de booleanos ", ra.RepDefd)
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
    	
    	for j:= 1; j < (ra.NNodes+1); j++ {
    		
    		if(j != ra.me){ //Enviar una peticion para acceder a la SC a todos los nodos salvo a si mismo
    			ra.Logger.LogLocalEvent("Enviando REQUEST a: " + strconv.Itoa(j), govec.GetDefaultLogOptions())
    			ra.ms.Send(j, ms.Request{TYPEREQUEST,ra.me,ra.OurSeqNum,ra.Escritor})//Operacion/IdNodo/ClockNodo/Escritor/Lector
    		}
    	}
    	ra.Mutex.Unlock()
    	ra.Logger.LogLocalEvent("Waitting for  " + strconv.Itoa(ra.OutRepCnt) + " replies", govec.GetDefaultLogOptions())
    	for {
    		ra.Mutex.Lock()//Asegurar la exlusión mutua cuabndo se consulta la variable
    		if ra.OutRepCnt == 0 {//Esperar a que haya recibido todas las respuestas
    			ra.Logger.LogLocalEvent("Exiting from Waitting for  " + strconv.Itoa(ra.OutRepCnt) + " replies", govec.GetDefaultLogOptions())
    			break
    		}
    		ra.Mutex.Unlock()    		
    	}
    	ra.Mutex.Unlock()
}
//Pre: Verdad
//Post: Realiza  el  PostProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PostProtocol(linea string, escribir bool){
	// TODO completa
    	ra.Mutex.Lock()
    	ra.ReqCS = false
    	ra.Logger.LogLocalEvent("Liberando a los que quedan", govec.GetDefaultLogOptions())
    	for d := 0; d<ra.NNodes ;d++ { //Liberar a todos los que estaban a la espera de permiso(True)
    		if(ra.RepDefd[d] && d != ra.me-1){ //Si estaba a la espera, hay que liberarlo
    			ra.Logger.LogLocalEvent("Liberando a " + strconv.Itoa(d+1), govec.GetDefaultLogOptions())
    			ra.RepDefd[d] = false;
    			ra.ms.Send(d+1,ms.Reply{TYPEREPLY,linea,escribir})
    		}
    	}
    	ra.Mutex.Unlock()
}

func (ra *RASharedDB) Stop(){
    ra.ms.Stop()
    ra.done <- true
}

func maxClock(reloj1 int, reloj2 int) int {
	if(reloj1 > reloj2){
		return reloj1
	}else {
		return reloj2
	}
}

func exlusion(evento1 bool, evento2 bool) bool{
	return evento1 || evento2
}

func (ra *RASharedDB) acquire_mutex(tipoOp bool){ //Escibir true, leer false
	ra.Logger.LogLocalEvent("Trying to SC", govec.GetDefaultLogOptions())
	ra.PreProtocol()
}

func (ra *RASharedDB) release_mutex(linea string, escritor bool){
	ra.PostProtocol(linea,escritor)
}

func (ra *RASharedDB) listening(){
		for {
			ra.Logger.LogLocalEvent("Esperando peticion ", govec.GetDefaultLogOptions())
			msg := ra.ms.Receive()//Se queda a la espera de recivir una petición
			ra.Logger.LogLocalEvent("Recibido un " + reflect.TypeOf(msg).String() , govec.GetDefaultLogOptions())
			if(reflect.TypeOf(msg).String() == "ms.Request"){ //Recibido REQUEST
				ra.Mutex.Lock()
					ra.HigSeqNum = maxClock(msg.(ms.Request).Clock,ra.OurSeqNum)
					Defer_it := ra.ReqCS&&((msg.(ms.Request).Clock>ra.OurSeqNum)||((msg.(ms.Request).Clock==ra.OurSeqNum)&&(msg.(ms.Request).Id>ra.me)))&&exlusion(msg.(ms.Request).Escritor,ra.Escritor)
					if Defer_it { //Tengo más prioridad asi que pongo al nodo de la petición a esperar
						ra.RepDefd[msg.(ms.Request).Id-1] = true
					} else {//Tengo menos prioridad o no quiero acceder a la seccion critica le doy el permiso al nodo
						ra.ms.Send(msg.(ms.Request).Id,ms.Reply{TYPEREPLY,"",false})
					}
					
				ra.Mutex.Unlock()
				
			} else if reflect.TypeOf(msg).String() == "ms.Reply" { //RECIBIDO REPLY
				if msg.(ms.Reply).Mode { //Recibida orden de escritura
					ra.Mutex.Lock()
						fmt.Println("Tipo del mensaje Write :", msg.(ms.Reply).Type)
						fmt.Println("Linea del mensaje Write :", msg.(ms.Reply).Response)
						fmt.Println("Modo del mensaje Write :", msg.(ms.Reply).Mode)
						//Escribir en el archivo (slice de escrituras)
						gestorfichero.EscribirFichero("memory.txt", msg.(ms.Reply).Response)
					ra.Mutex.Unlock()
				} else { 
					fmt.Println("Tipo del mensaje Write :", msg.(ms.Reply).Type)
					fmt.Println("Linea del mensaje Write :", msg.(ms.Reply).Response)
					fmt.Println("Modo del mensaje Write :", msg.(ms.Reply).Mode)
					ra.Mutex.Lock()
						fmt.Println("Contador de espera 1: ", ra.OutRepCnt)
						ra.OutRepCnt = ra.OutRepCnt - 1
						fmt.Println("Contador de espera 2: ", ra.OutRepCnt)
					ra.Mutex.Unlock()
				}
			}
		}
}
