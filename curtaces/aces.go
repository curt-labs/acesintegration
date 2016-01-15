package curtaces

import (
	"log"
	"strconv"
	"strings"

	"github.com/curt-labs/acesintegration/database"
)

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

//GET ALL ACES Vehicle Apps
var (
	getAcesVehiclesStmt = `select bv.yearID, lower(vma.MakeName), lower(vmo.ModelName), s.SubmodelName, cat.name, ca.value, p.oldPartNumber
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
		where p.brandID = 1
		`
)

func GetAcesVehicles() (map[string][]AcesVehicleApplication, error) {
	vs := make(map[string][]AcesVehicleApplication)
	temp := make(map[string][]Config)
	if err := database.Init(); err != nil {
		return vs, err
	}

	rows, err := database.DB.Query(getAcesVehiclesStmt)
	if err != nil {
		return vs, err
	}
	var part, submodel, cat, ca *string
	count := 0
	for rows.Next() {
		var v AcesVehicleApplication

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
		mapKey := strconv.FormatFloat(v.Year, 'f', 1, 64) + "|" + v.Make + "|" + v.Model + "|" + v.Submodel + "|" + v.Part
		temp[mapKey] = append(temp[mapKey], con)
	}

	//combine configs, by Base+Submodel+PartID
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
			Part:     varray[4],
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
		count++
	}
	log.Print(count, " individual ACES applications.")
	return vs, nil
}
