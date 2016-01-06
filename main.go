package main

import (
	"encoding/csv"
	"fmt"
	"github.com/curt-labs/acesintegration/database"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	cvs, err := getCurtVehicles()
	if err != nil {
		fmt.Errorf("%v", err)
	}

	avs, err := getAcesVehicles()
	if err != nil {
		fmt.Errorf("%v", err)
	}

	err = process(cvs, avs)
	if err != nil {
		fmt.Errorf("%v", err)
	}
	log.Print("DONE")
}

//THE DEAL
//If curtvehicle has style "all" => Integrity Y
//If curtvehicle style != "all" => if acesvehicle has submodels or configs => Integrity Y
//								=> if acesvehicle has no submodels and no configs => Integrity N

type CurtVehicleApplication struct {
	Key   string //Year|Make|Model
	Year  float64
	Make  string
	Model string
	Style string
	Part  string
}

type AcesVehicleApplication struct {
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

//CheckIntegrity
func process(cvs map[string][]CurtVehicleApplication, avs map[string][]AcesVehicleApplication) error {
	//set up file
	f, err := os.Create("integrity.csv")
	if err != nil {
		return err
	}
	writer := csv.NewWriter(f)

	//write header
	header := []string{"Part", "Year", "Make", "Model", "Style", "Integrated", "Year", "Make", "Model", "Submodel", "Configs", "Notes"}
	err = writer.Write(header)
	if err != nil {
		return err
	}

	//check integrity & write lines
	for key, cv := range cvs {
		if matchingAcesVehicle, ok := avs[key]; !ok {
			//write no match on Base Vehicle at all NO integrity
			line := []string{cv[0].Part, strconv.FormatFloat(cv[0].Year, 'f', 1, 64), cv[0].Make, cv[0].Model, cv[0].Style, "N", "", "", "", "", "", "Base Vehicle (year|make|model) does not exist"}
			err = writer.Write(line)
			if err != nil {
				return err
			}
			writer.Flush()
		} else {
			lines := determineIntegrity(cv, matchingAcesVehicle)
			err = writer.WriteAll(lines)
			if err != nil {
				return err
			}
			writer.Flush()
		}
	}
	return nil
}

func determineIntegrity(cvs []CurtVehicleApplication, avs []AcesVehicleApplication) [][]string {
	var lines [][]string

	for _, cv := range cvs {
		//CurtVehicle style == "all"
		if strings.ToLower(cv.Style) == "all" {
			integrity := "N"
			notes := "Curt Vehicle Style = all BUT an Aces non-config/non-submodel DOES NOT exists (need base)"
			for _, av := range avs {
				if av.Submodel == "" && len(av.Configs) == 0 {
					integrity = "Y"
					notes = "Curt Vehicle Style = all AND an Aces non-config/non-submodel exists"
				}
			}
			for _, av := range avs {
				var cons string //stringify configs into comma separated sets
				for i, c := range av.Configs {
					if i > 0 {
						cons += ","
					}
					cons += c.Type + ":" + c.Value
				}
				lines = append(lines, []string{cv.Part, strconv.FormatFloat(cv.Year, 'f', 1, 64), cv.Make, cv.Model, cv.Style, integrity, strconv.FormatFloat(av.Year, 'f', 1, 64), av.Make, av.Model, av.Submodel, cons, notes})
			}
		} else {
			//CurtVehicle style != "all"
			integrity := "N" //Aces vehicle DOES NOT have a submodel or config(s)
			notes := "Curt Vehicle Style != all AND an Aces vehicle with either configs or a submodel DOES NOT exists"
			for _, av := range avs {
				if av.Submodel != "" || len(av.Configs) > 0 {
					//Aces vehicle DOES have a submodel or config(s)
					integrity = "Y"
					notes = "Curt Vehicle Style != all AND an Aces vehicle with either configs or a submodel DOES exists"
				}
			}

			for _, av := range avs {
				var cons string //stringify configs into comma separated sets
				for i, c := range av.Configs {
					if i > 0 {
						cons += ","
					}
					cons += c.Type + ":" + c.Value
				}
				lines = append(lines, []string{cv.Part, strconv.FormatFloat(cv.Year, 'f', 1, 64), cv.Make, cv.Model, cv.Style, integrity, strconv.FormatFloat(av.Year, 'f', 1, 64), av.Make, av.Model, av.Submodel, cons, notes})
			}
		}
	}

	return lines
}

//GET ALL CURT && ACES Vehicle Apps
var (
	getCurtVehiclesStmt = `select y.Year, ma.Make, mo.Model, s.Style, p.oldPartNumber
		from Vehicle v 
		join VehiclePart vp on vp.vehicleID = v.vehicleID
		join Part p on p.partID = vp.partID
		join Year y on y.yearID = v.yearID
		join Make ma on ma.makeID = v.makeID
		join Model mo on mo.modelID = v.modelID
		join Style s on s.styleID = v.styleID
		where p.brandID = 1`
	getAcesVehiclesStmt = `select bv.yearID, vma.MakeName, vmo.ModelName, s.SubmodelName, cat.name, ca.value, p.oldPartNumber
		from vcdb_Vehicle vv
		join vcdb_VehiclePart vp on vp.vehicleID = vv.ID
		join Part p on p.partID = vp.partNumber
		join BaseVehicle bv on bv.ID = vv.BaseVehicleID
		join vcdb_Make vma on vma.ID = bv.MakeID
		join vcdb_Model vmo on vmo.ID = bv.ModelID
		left join Submodel s on s.ID = vv.SubmodelID
		left join VehicleConfigAttribute vca on vca.ID = vv.ConfigID
		left join ConfigAttribute ca on ca.ID = vca.AttributeID 
		left join ConfigAttributeType cat on ca.ConfigAttributeTypeID = cat.ID
		`
)

func getCurtVehicles() (map[string][]CurtVehicleApplication, error) {
	vs := make(map[string][]CurtVehicleApplication)
	if err := database.Init(); err != nil {
		return vs, err
	}

	rows, err := database.DB.Query(getCurtVehiclesStmt)
	if err != nil {
		return vs, err
	}
	var v CurtVehicleApplication
	var style, part *string

	for rows.Next() {
		err = rows.Scan(
			&v.Year,
			&v.Make,
			&v.Model,
			&style,
			&part,
		)
		if err != nil {
			return vs, err
		}

		if style != nil {
			v.Style = *style
		}
		if part != nil {
			v.Part = *part
		}
		v.Key = strconv.FormatFloat(v.Year, 'f', 1, 64) + "|" + v.Make + "|" + v.Model
		vs[v.Key] = append(vs[v.Key], v)
	}

	return vs, nil

}

func getAcesVehicles() (map[string][]AcesVehicleApplication, error) {
	vs := make(map[string][]AcesVehicleApplication)
	temp := make(map[string][]Config)
	if err := database.Init(); err != nil {
		return vs, err
	}

	rows, err := database.DB.Query(getAcesVehiclesStmt)
	if err != nil {
		return vs, err
	}
	var v AcesVehicleApplication
	var part, submodel, cat, ca *string

	for rows.Next() {
		var con Config
		err = rows.Scan(
			&v.Year,
			&v.Make,
			&v.Model,
			&submodel,
			&cat,
			&ca,
			&part,
		)
		if err != nil {
			return vs, err
		}
		if submodel != nil {
			v.Submodel = *submodel
		}
		if cat != nil {
			con.Type = strings.TrimSpace(*cat)
		}
		if ca != nil {
			con.Value = strings.TrimSpace(*ca)
		}
		if part != nil {
			v.Part = *part
		}
		mapKey := strconv.FormatFloat(v.Year, 'f', 1, 64) + "|" + v.Make + "|" + v.Model + "|" + v.Submodel
		temp[mapKey] = append(temp[mapKey], con)
	}

	for i, cons := range temp {
		varray := strings.Split(i, "|")
		yearFl, err := strconv.ParseFloat(varray[0], 64)
		if err != nil {
			return vs, err
		}

		v := AcesVehicleApplication{
			Key:      varray[0] + "|" + varray[1] + "|" + varray[2],
			Year:     yearFl,
			Make:     varray[1],
			Model:    varray[2],
			Submodel: strings.TrimSpace(varray[3]),
		}
		for _, c := range cons {
			if c.Type == "" && c.Value == "" {
				continue
			}
			var skip bool
			for _, vc := range v.Configs {
				if vc.Type == c.Type && vc.Value == c.Value {
					skip = true
				}
			}
			if skip {
				continue
			}
			v.Configs = append(v.Configs, Config{Type: c.Type, Value: c.Value})

		}
		vs[v.Key] = append(vs[v.Key], v)
	}

	return vs, nil

}