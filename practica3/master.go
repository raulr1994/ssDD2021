/*
* AUTOR: Rafael Tolosana Calasanz
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: octubre de 2021
* FICHERO: worker.go
* DESCRIPCIÓN: contiene la funcionalidad esencial para realizar los servidores
*				correspondientes la practica 3
*/
package main

import (
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"
	"fmt"
	"practica3/com"
)

const (
	NORMAL   = iota // NORMAL == 0
	DELAY    = iota // DELAY == 1
	CRASH    = iota // CRASH == 2
	OMISSION = iota // IOTA == 3
)

type PrimesImpl struct {
	delayMaxMilisegundos int
	delayMinMiliSegundos int
	behaviourPeriod      int
	behaviour            int
	i                    int
	mutex                sync.Mutex
}

func sendRequestToWorker(endpoint string, interval com.TPInterval, reply *[]int){
	client, err := rpc.DialHTTP("tcp", endpoint)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	//var reply []int
	err = client.Call("PrimesImpl.FindPrimes", interval, &reply)
	if err != nil {
		log.Fatal("primes error:", err)
	}
}


func (p *PrimesImpl) FindPrimesToWorker(interval com.TPInterval, primeList *[]int) error {
		//Tener en cuenta aqui que luego el worker podría caerse
		//Elegir al primer worker que esté disponible todos funcionan de forma ininterrumpida
		var primes []int
		go sendRequestToWorker("192.168.1.3:30000", interval, &primes)
		return nil
}

func main() {
	if len(os.Args) == 2 {
		time.Sleep(10 * time.Second)
		rand.Seed(time.Now().UnixNano())
		primesImpl := new(PrimesImpl)
		primesImpl.delayMaxMilisegundos = 4000
		primesImpl.delayMinMiliSegundos = 2000
		primesImpl.behaviourPeriod = 4
		primesImpl.i = 1
		primesImpl.behaviour = NORMAL
		rand.Seed(time.Now().UnixNano())

		rpc.Register(primesImpl)
		rpc.HandleHTTP()
		l, e := net.Listen("tcp", os.Args[1])
		if e != nil {
			log.Fatal("listen error:", e)
		}
		http.Serve(l, nil)
	} else {
		fmt.Println("Usage: go run worker.go <ip:port>")
	}
}
