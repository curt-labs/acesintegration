package main

import (
	"github.com/curt-labs/acesintegration/curtaces"
	"github.com/curt-labs/acesintegration/curtdmi"
	"log"
)

func main() {

	// curtAces()
	curtDmi()

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
	dvs, err := curtdmi.GetDMIVehicleApplications()
	if err != nil {
		log.Print(err)
	}
	log.Print(dvs)
}
