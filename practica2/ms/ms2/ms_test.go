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
)

type Request struct {
	Id int
}

type Reply struct{
	Response string
}

func TestSendReceiveMessage(t *testing.T) {
	p1 := New(1, "./users2.txt", []Message{Request{}, Reply{}})
	
	//p1.Send(2, Request{1})
	p1.Send(2, Message1{Operation: "ataque", Sender: 1, Clock: 12, Write: true})
}


