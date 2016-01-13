package curtdmi

import (
	"encoding/xml"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/curt-labs/acesintegration/database"
)

type Result struct {
	XMLName xml.Name
	Header  Header `xml:"Header"`
	Apps    []App  `xml:"App"`
}

//Xml parsed
type Header struct {
	XMLName xml.Name
	Company string `xml:"Company"`
}

type App struct {
	XMLName xml.Name `xml:"App"`
	// MfrLabel string   `xml:"MfrLabel"`
	// Action   string   `xml:"id,attr"`
	BaseVehicle  BaseVehicle `xml:"BaseVehicle"`
	SubModel     SubModel    `xml:"SubModel"`
	Notes        []string    `xml:"Note"`
	PartID       int         `xml:"Part"`
	BodyType     BodyType
	BodyNumDoors BodyNumDoors
	WheelBase    WheelBase
	DriveType    DriveType
	BedLength    BedLength
	BedType      BedType
}

type BaseVehicle struct {
	ID    int `xml:"id,attr"`
	Year  int
	Make  string
	Model string
}
type SubModel struct {
	ID   int `xml:"id,attr"`
	Name string
}

type BodyType struct {
	ID int `xml:"id,attr"`
}
type BodyNumDoors struct {
	ID int `xml:"id,attr"`
}
type WheelBase struct {
	ID int `xml:"id,attr"`
}
type DriveType struct {
	ID int `xml:"id,attr"`
}
type BedLength struct {
	ID int `xml:"id,attr"`
}
type Aspiration struct {
	ID int `xml:"id,attr"`
}
type BedType struct {
	ID int `xml:"id,attr"`
}

//dmi Vehicle
type DmiVehicleApplication struct {
	Key      string //Year|Make|Model
	Year     float64
	Make     string
	Model    string
	Submodel string
	Configs  []Config
	Part     string
}

type Config struct {
	Type  string
	Value string
}

func GetDCIVehicles() (map[string][]DmiVehicleApplication, error) {
	vehicleMap := make(map[string][]DmiVehicleApplication)
	ds, err := GetDMIVehicleApplications()
	if err != nil {
		return vehicleMap, err
	}

	for _, d := range ds {
		vehicleMap[d.Key] = append(vehicleMap[d.Key], d)
	}
	return vehicleMap, err
}

func GetDMIVehicleApplications() ([]DmiVehicleApplication, error) {
	var ds []DmiVehicleApplication

	res, err := ParseXML()
	if err != nil {
		return ds, err
	}

	basemap, err := getVcdbBaseMap()
	if err != nil {
		return ds, err
	}

	submap, err := getVcdbSubmodelMap()
	if err != nil {
		return ds, err
	}

	err = getVcdbConfigMaps()
	if err != nil {
		return ds, err
	}
	for _, app := range res.Apps {
		var d DmiVehicleApplication

		//assign base vehicle
		if vcdbBase, ok := basemap[app.BaseVehicle.ID]; !ok {
			//TODO non-vcdb base
			continue
		} else {
			d.Year = float64(vcdbBase.Year)
			d.Make = strings.TrimSpace(vcdbBase.Make)   //TODO lower to match curt?
			d.Model = strings.TrimSpace(vcdbBase.Model) //TODO lower to match curt?
		}
		//assign submodel
		if app.SubModel.ID > 0 {
			if vcdbSub, ok := submap[app.SubModel.ID]; !ok {
				//TODO non-vcdb submodel
				continue
			} else {
				d.Submodel = strings.TrimSpace(vcdbSub.Name)
			}
		}
		//assign configs
		d.processConfigs(app)
		//assign part
		d.Part = strconv.Itoa(app.PartID)

		//make vehicle key for lookup comparison
		d.Key = strconv.FormatFloat(d.Year, 'f', 1, 64) + "|" + d.Make + "|" + d.Model

		//add to array
		ds = append(ds, d)
	}
	log.Print(len(ds), " individual DCI Applications.")
	return ds, err
}

func ParseXML() (Result, error) {
	var res Result

	//Get File
	f, err := os.Open("CUR20151219_ACESV3.xml")
	if err != nil {
		return res, err
	}
	defer f.Close()
	err = xml.NewDecoder(f).Decode(&res)
	return res, err
}

func (d *DmiVehicleApplication) processConfigs(app App) {

	if app.BedLength.ID > 0 {
		c := Config{
			Type:  "Bed Length",
			Value: BedLengthMap[app.BedLength.ID],
		}
		d.Configs = append(d.Configs, c)
	}
	if app.BedType.ID > 0 {
		c := Config{
			Type:  "Bed Type",
			Value: BedTypeMap[app.BedType.ID],
		}
		d.Configs = append(d.Configs, c)
	}
	if app.BodyType.ID > 0 {
		c := Config{
			Type:  "Body Type",
			Value: BodyTypeMap[app.BodyType.ID],
		}
		d.Configs = append(d.Configs, c)
	}
	if app.BodyNumDoors.ID > 0 {
		c := Config{
			Type:  "Number of Doors",
			Value: BodyNumDoorsMap[app.BodyNumDoors.ID],
		}
		d.Configs = append(d.Configs, c)
	}
	if app.WheelBase.ID > 0 {
		c := Config{
			Type:  "Wheel Base",
			Value: WheelBaseMap[app.WheelBase.ID],
		}
		d.Configs = append(d.Configs, c)
	}
	if app.DriveType.ID > 0 {
		c := Config{
			Type:  "Drive Type",
			Value: DriveTypeMap[app.DriveType.ID],
		}
		d.Configs = append(d.Configs, c)
	}
	return
}

//vcdb Maps
var (
	vcdbBaseVehicleStmt = `select b.BaseVehicleID, b.YearID, ma.Makename, mo.ModelName
		from BaseVehicle b
		join Make ma on ma.MakeID = b.MakeID
		join Model mo on mo.ModelID = b.ModelID`
	vcdbSubmodelStmt = `select s.SubmodelID, s.SubmodelName from Submodel s`

	BodyTypeMap     map[int]string
	BodyNumDoorsMap map[int]string
	WheelBaseMap    map[int]string
	DriveTypeMap    map[int]string
	BedLengthMap    map[int]string
	BedTypeMap      map[int]string
)

func getVcdbBaseMap() (map[int]BaseVehicle, error) {
	themap := make(map[int]BaseVehicle)

	if err := database.Init(); err != nil {
		return themap, err
	}

	rows, err := database.VCDB.Query(vcdbBaseVehicleStmt)
	if err != nil {
		return themap, err
	}

	var b BaseVehicle

	for rows.Next() {
		err = rows.Scan(
			&b.ID,
			&b.Year,
			&b.Make,
			&b.Model,
		)
		if err != nil {
			return themap, err
		}
		themap[b.ID] = b
	}

	return themap, nil
}

func getVcdbSubmodelMap() (map[int]SubModel, error) {
	themap := make(map[int]SubModel)
	if err := database.Init(); err != nil {
		return themap, err
	}

	rows, err := database.VCDB.Query(vcdbSubmodelStmt)
	if err != nil {
		return themap, err
	}
	var s SubModel
	for rows.Next() {
		err = rows.Scan(
			&s.ID,
			&s.Name,
		)
		if err != nil {
			return themap, err
		}
		themap[s.ID] = s
	}
	return themap, nil
}

func getVcdbConfigMaps() error {
	var err error
	configTypeArray := map[string]string{
		"BodyType":     "BodyTypeName",
		"BodyNumDoors": "BodyNumDoors",
		"WheelBase":    "WheelBase",
		"DriveType":    "DriveTypeName",
		"BedLength":    "BedLength",
		"BedType":      "BedTypeName",
	}
	for configType, name := range configTypeArray {
		id := configType + "ID"
		stmt := "select " + id + ", " + name + " from " + configType
		switch configType {
		case "BodyType":
			BodyTypeMap, err = getVcdbConfigMap(stmt)
			if err != nil {
				return err
			}
		case "BodyNumDoors":
			BodyNumDoorsMap, err = getVcdbConfigMap(stmt)
			if err != nil {
				return err
			}
		case "WheelBase":
			WheelBaseMap, err = getVcdbConfigMap(stmt)
			if err != nil {
				return err
			}
		case "DriveType":
			DriveTypeMap, err = getVcdbConfigMap(stmt)
			if err != nil {
				return err
			}
		case "BedLength":
			BedLengthMap, err = getVcdbConfigMap(stmt)
			if err != nil {
				return err
			}
		case "BedType":
			BedTypeMap, err = getVcdbConfigMap(stmt)
			if err != nil {
				return err
			}
		default:
			log.Fatal("MISS", configType)
		}
	}
	return nil
}

func getVcdbConfigMap(stmt string) (map[int]string, error) {
	themap := make(map[int]string)
	if err := database.Init(); err != nil {
		return themap, err
	}

	rows, err := database.VCDB.Query(stmt)
	if err != nil {
		return themap, err
	}
	var s string
	var i int
	for rows.Next() {
		err = rows.Scan(&i, &s)
		if err != nil {
			return themap, err
		}
		themap[i] = strings.TrimSpace(s)
	}
	return themap, err
}

//NOT USED ACES CONFG TYPES that are in our configAttributeType table:
// FuelType
// Engine
// BrakeABS
// BrakeSystem
// CylinderHeadType
// EngineDesignation
// EngineManafacturer
// EngineVersion
// EngineVIN
// FrontBrakeType
// FrontSpringType
// FuelDeliverySubType
// FuelDeliveryType
// FuelSystemControlType
// FuelSystemDesign
// IgnitionSystemType
// ManufacturerBodyCode
// PowerOutput
// RearBrakeType
// RearSpringType
// SteeringSystem
// SteeringType
// TransmissionElectronicControlled
// Transmission
// TransmissionBase
// TransmissionControlType
// TransmissionManufacturer
// TransmissionnumberOfSpeeds
// TransmissionType
// ValvesPerEngine
