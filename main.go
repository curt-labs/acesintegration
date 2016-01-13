package main

import (
	"github.com/curt-labs/acesintegration/curtaces"
	"github.com/curt-labs/acesintegration/curtdmi"
	"log"
)

func main() {

	curtAces()
	curtDmi()

}

func curtAces() {
	cvs, err := curtaces.GetCurtVehicles()
	if err != nil {
		log.Print(err)
	}
	log.Print(len(cvs), " Curt Base vehicles in DB")

	avs, err := curtaces.GetAcesVehicles()
	if err != nil {
		log.Print(err)
	}

	err = curtaces.Process(cvs, avs)
	if err != nil {
		log.Print(err)
	}
	log.Print(len(avs), " ACES Base vehicles in DB")

	log.Print("DONE")
}

func curtDmi() {
	dvs, err := curtdmi.GetDCIVehicles()
	if err != nil {
		log.Print(err)
	}
	log.Print(len(dvs), " DCI Base vehicles")

	cvs, err := curtaces.GetCurtVehicles()
	if err != nil {
		log.Print(err)
	}
	log.Print(len(cvs), " Curt Base vehicles in DB")

	err = curtdmi.ProcessCurtToDci(cvs, dvs)
	if err != nil {
		log.Print(err)
	}
	log.Print("DONE")
}
