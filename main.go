package main

import (
	"github.com/curt-labs/acesintegration/curtaces"
	"github.com/curt-labs/acesintegration/curtdmi"
	"log"
)

func main() {

	// curtAces()
	curtDmi()
	// acesDmi()

}

func curtAces() {
	cvs, err := curtaces.GetCurtVehicles()
	if err != nil {
		log.Print(err)
	}

	avs, err := curtaces.GetAcesVehicles()
	if err != nil {
		log.Print(err)
	}

	err = curtaces.Process(cvs, avs)
	if err != nil {
		log.Print(err)
	}
	log.Print("DONE")
}

func curtDmi() {
	dvs, err := curtdmi.GetDCIVehicles()
	if err != nil {
		log.Print(err)
	}
	cvs, err := curtaces.GetCurtVehicles()
	if err != nil {
		log.Print(err)
	}
	err = curtdmi.ProcessCurtToDci(cvs, dvs)
	if err != nil {
		log.Print(err)
	}
	log.Print("DONE")
}

//TODO ?
// func acesDmi() {
// 	dvs, err := curtdmi.GetDCIVehicles()
// 	if err != nil {
// 		log.Print(err)
// 	}
// 	avs, err := curtaces.GetAcesVehicles()
// 	if err != nil {
// 		log.Print(err)
// 	}
// 	err = curtdmi.ProcessAcesToDci(avs, dvs)
// 	if err != nil {
// 		log.Print(err)
// 	}
// 	log.Print("DONE")
// }
