package curtaces

import (
	"log"
	"strconv"

	"github.com/curt-labs/acesintegration/database"
)

type CurtVehicleApplication struct {
	Key   string //Year|Make|Model
	Year  float64
	Make  string
	Model string
	Style string
	Part  string
}

//GET ALL CURTVehicle Apps
var (
	getCurtVehiclesStmt = `select y.Year, lower(ma.Make), lower(mo.Model), s.Style, p.oldPartNumber
		from Vehicle v 
		join VehiclePart vp on vp.vehicleID = v.vehicleID
		join Part p on p.partID = vp.partID
		join Year y on y.yearID = v.yearID
		join Make ma on ma.makeID = v.makeID
		join Model mo on mo.modelID = v.modelID
		join Style s on s.styleID = v.styleID
		where p.brandID = 1`
)

func GetCurtVehicles() (map[string][]CurtVehicleApplication, error) {
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
	count := 0
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
		count++
	}
	log.Print(count, " individual CURT applications.")
	return vs, nil
}
