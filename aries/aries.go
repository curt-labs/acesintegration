package aries

import (
	"strconv"
	"strings"

	"github.com/curt-labs/acesintegration/database"
	"gopkg.in/mgo.v2/bson"
)

type AriesVehicleApplication struct {
	Key   string  //Year|Make|Model
	Year  float64 `bson:"year"`
	Make  string  `bson:"make"`
	Model string  `bson:"model"`
	Style string  `bson:"style"`
	Part  string  `bson:"part"`
}

type MgoAriesVehicleApplication struct {
	Year    string `bson:"year"`
	Make    string `bson:"make"`
	Model   string `bson:"model"`
	Style   string `bson:"style"`
	PartIds []int  `bson:"parts"`
}

func GetAriesVehicleApplications() (map[string][]AriesVehicleApplication, error) {
	avas := make(map[string][]AriesVehicleApplication)
	var err error

	err = database.InitMongo()
	if err != nil {
		return avas, err
	}

	names, err := getCollectionNames()
	if err != nil {
		return avas, err
	}

	for _, name := range names {
		applications, err := getVehiclesFromCollection(name)
		if err != nil {
			return avas, err
		}
		ariesApps, err := convertMgoToAries(applications)
		if err != nil {
			return avas, err
		}

		// avas = append(avas, ariesApps...)
		for _, ariesApp := range ariesApps {
			avas[ariesApp.Key] = append(avas[ariesApp.Key], ariesApp)
		}
	}

	return avas, err
}

func getCollectionNames() ([]string, error) {
	err := database.InitMongo()
	if err != nil {
		return []string{}, err
	}
	return database.MongoSession.DB(database.MongoDB).CollectionNames()
}

func getVehiclesFromCollection(name string) ([]MgoAriesVehicleApplication, error) {
	var avas []MgoAriesVehicleApplication
	var err error
	err = database.InitMongo()
	if err != nil {
		return avas, err
	}
	err = database.MongoSession.DB(database.MongoDB).C(name).Find(bson.M{}).All(&avas)
	return avas, err
}

func convertMgoToAries(mgoApps []MgoAriesVehicleApplication) ([]AriesVehicleApplication, error) {
	var ariesApps []AriesVehicleApplication
	var err error
	for _, mgoApp := range mgoApps {
		y, err := strconv.ParseFloat(mgoApp.Year, 64)
		if err != nil {
			continue //no year?
		}

		for _, partId := range mgoApp.PartIds {
			a := AriesVehicleApplication{
				Year:  y,
				Key:   strings.TrimSpace(mgoApp.Year) + "|" + strings.TrimSpace(mgoApp.Make) + "|" + strings.TrimSpace(mgoApp.Model),
				Make:  strings.TrimSpace(mgoApp.Make),
				Model: strings.TrimSpace(mgoApp.Model),
				Style: strings.TrimSpace(mgoApp.Style),
				Part:  strconv.Itoa(partId),
			}
			ariesApps = append(ariesApps, a)
		}
	}
	return ariesApps, err
}
