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
	"fmt"
	"reflect"
	"github.com/DistributedClocks/GoVector/govec"
	"testing"
	"example.gestorfichero/gestorfichero"
	"strconv"
	"example.ms/ms"
)

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

func (ra *RASharedDB) release_mutex(){
	ra.PostProtocol()
}

func (ra *RASharedDB) listening(){
		for {
			ra.Logger.LogLocalEvent("Esperando peticion ", govec.GetDefaultLogOptions())
			msg := ra.ms.Receive()//Se queda a la espera de recivir una petición
			ra.Logger.LogLocalEvent("Recibido un " + reflect.TypeOf(msg).String() , govec.GetDefaultLogOptions())
			//fmt.Println("Recibido un ", reflect.TypeOf(msg))
			if(reflect.TypeOf(msg).String() == "ms.Request"){ //Recibido REQUEST
				fmt.Println(msg.(ms.Request).Type)
				fmt.Println("Id del mensaje Request :", msg.(ms.Request).Id)
				fmt.Println("Clock del mensaje Request :", msg.(ms.Request).Clock)
				fmt.Println("Estado del mensaje Request :", msg.(ms.Request).Escritor)
				fmt.Println("Mi id :", ra.me)
				fmt.Println("Mi reloj :", ra.OurSeqNum)
				fmt.Println("Mi estado :", ra.Escritor)
				ra.Mutex.Lock()
					ra.HigSeqNum = maxClock(msg.(ms.Request).Clock,ra.OurSeqNum)
					Defer_it := ra.ReqCS&&((msg.(ms.Request).Clock>ra.OurSeqNum)||((msg.(ms.Request).Clock==ra.OurSeqNum)&&(msg.(ms.Request).Id>ra.me)))&&exlusion(msg.(ms.Request).Escritor,ra.Escritor)
					fmt.Println("Defer_it ", Defer_it)
					fmt.Println("1 ", ra.ReqCS)
					fmt.Println("2 ", msg.(ms.Request).Clock>ra.OurSeqNum)
					fmt.Println("3 ", msg.(ms.Request).Clock==ra.OurSeqNum)
					fmt.Println("4 ", msg.(ms.Request).Id>ra.me)
					fmt.Println("5 ", exlusion(msg.(ms.Request).Escritor,ra.Escritor))
					if Defer_it { //Tengo más prioridad asi que pongo al nodo de la petición a esperar
						fmt.Println("Soy superior de ", msg.(ms.Request).Id)
						ra.RepDefd[msg.(ms.Request).Id] = true
					} else {//Tengo menos prioridad o no quiero acceder a la seccion critica le doy el permiso al nodo
						fmt.Println("Aceptando superioridad de ", msg.(ms.Request).Id)
						ra.ms.Send(msg.(ms.Request).Id,ms.Reply{TYPEREPLY,MSGFREE})
					}
					
				ra.Mutex.Unlock()
				
			}else{ //RECIBIDO REPLY
				fmt.Println("Id del mensaje Reply :", msg.(ms.Reply).Type)
				fmt.Println("Nombre del mensaje Reply :", msg.(ms.Reply).Response)
				ra.Mutex.Lock()
					fmt.Println("Contador de espera 1: ", ra.OutRepCnt)
					ra.OutRepCnt = ra.OutRepCnt - 1
					fmt.Println("Contador de espera 2: ", ra.OutRepCnt)
				ra.Mutex.Unlock()
			}
		}
}

func (ra *RASharedDB) SendWriteAll(linea string){  	
    	for j:= 1; j < (ra.NNodes+1); j++ {
    		if(j != ra.me){ //Enviar una peticion para que todos los nodos escriban en sus ficheros
    			fmt.Println("Enviando a: ", j)
    			ra.ms.Send(j, ms.Reply{TYPEWRITE,linea})//Operacion/Trozo codigo que escribir
    		}
    	}
}

func TestReader(t *testing.T){
    	//1º Recoger parametros
    	//pId := flag.Int("ID",0,"El id del nodo RA")
    	//2º inicializar variables
    	nExp := 1
    	opP := false //Leer(False),Escibir(True)
    	nNodes := 2
    	rai := New(1,"./users.txt",nNodes,opP);//Pid, nombre fichero, numero de nodos y si es escritor o lector
    	//bufferDeLectura := "" //Donde guardar lo que se lee
    	//bufferDeEscritura := "" //Que escribir
    	gestorfichero.LimpiarFichero("memory.txt")
	
    	go rai.listening() //Lanzar el proceso de escucha de peticiones
    	//3º Ejecutar Richard agravala
    	for i := 0; (i < nExp); i++ {
		rai.acquire_mutex(opP)//Pedir adquirir el mutex
		rai.Logger.LogLocalEvent("Entering to SC", govec.GetDefaultLogOptions())
		//Acceso a la CS para escibir o leer
		if opP {
			fmt.Println("Escribiendo")
			rai.Logger.LogLocalEvent("Writting by " + strconv.Itoa(rai.me), govec.GetDefaultLogOptions())
			gestorfichero.EscribirFichero("memory.txt", strconv.Itoa(i)) + "Escrito por el escritor " + strconv.Itoa(rai.me))
				//Mandar a todos los nodos que he escrito para que lo escriban ellos también
			ra.SendWriteAll(strconv.Itoa(i) + "Escrito por el escritor " + strconv.Itoa(rai.me))
		} else {
			fmt.Println("Leyendo")
			rai.Logger.LogLocalEvent("Reading by " + strconv.Itoa(rai.me), govec.GetDefaultLogOptions())
			lectura := gestorfichero.LeerFichero("memory.txt")
			gestorfichero.MostrarLectura(lectura)
		}
		fmt.Println("Exiting from SC")
		//Enviar al resto de nodos que ha escrito un escritor
		rai.Logger.LogLocalEvent("Exiting from SC", govec.GetDefaultLogOptions())
		rai.release_mutex()//Liberar el mutex
		//Dormir un rato antes de hacer la siguiente peticion
	}
	fmt.Println("Finishing " + strconv.Itoa(rai.me))
	for{
		//Para no dejar ningun nodo desatendido que le falte alguna peticion
	}
}
