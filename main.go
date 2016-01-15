package main

import (
	"github.com/curt-labs/acesintegration/aries"
	"github.com/curt-labs/acesintegration/curtaces"
	"github.com/curt-labs/acesintegration/curtdmi"
	"github.com/curt-labs/acesintegration/database"
	"log"
)

func main() {
	err := database.TestMongoConnection()
	if err != nil {
		log.Fatal(err)
	}
	curtAces()
	curtDmi()
	ariesDci()

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
	dvs, err := curtdmi.GetDCIVehicles("CUR20151219_ACESV3.xml")
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

func ariesDci() {
	dvs, err := curtdmi.GetDCIVehicles("CUR20151219_ACESV3.xml")
	if err != nil {
		log.Print(err)
	}
	log.Print(len(dvs), " DCI Base vehicles")
	avs, err := aries.GetAriesVehicleApplications()
	if err != nil {
		log.Print(err)
	}
	log.Print(len(avs), " Aries Base vehicles in Mongo")

	err = aries.ProcessAriesToDci(avs, dvs)
	if err != nil {
		log.Print(err)
	}

	log.Print("DONE")
}
