/*
* AUTOR: Rafael Tolosana Calasanz
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2021
* FICHERO: ms_test.go
* DESCRIPCIÓN: Implementación de un sistema de mensajería asíncrono, insipirado en el Modelo Actor
*/
package ms

import (
	"testing"
	"reflect"
	"fmt"
)

type Request struct {
	Type string
	Id int
}

type Reply struct{
	Type string
	Response string
}


func TestSendReceiveMessage(t *testing.T) {
		p1 := New(1, "./users2.txt", []Message{Request{}, Reply{}})
	
		//p1.Send(2, Request{1})
		for {
			msg := p1.Receive()
			fmt.Println("Recibido un ", reflect.TypeOf(msg))
			if(reflect.TypeOf(msg).String() == "ms.Request"){
				fmt.Println("Id del mensaje Request :", msg.(Request).Type)
				fmt.Println("Id del mensaje Request :", msg.(Request).Id)
			}else{
				fmt.Println("Id del mensaje Request :", msg.(Reply).Type)
				fmt.Println("Nombre del mensaje Reply :", msg.(Reply).Response)
			}
		}
}	


