package curtaces

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

//CheckIntegrity
func Process(cvs map[string][]CurtVehicleApplication, avs map[string][]AcesVehicleApplication) error {
	//set up file
	f, err := os.Create("integrityCurtAces.csv")
	if err != nil {
		return err
	}
	writer := csv.NewWriter(f)

	//write header
	header := []string{"Part", "Year", "Make", "Model", "Style", "Integrated", "ACES Part Assoc.", "Year", "Make", "Model", "Submodel", "Configs", "Notes"}
	err = writer.Write(header)
	if err != nil {
		return err
	}

	//check integrity & write lines
	for key, cv := range cvs {
		if matchingAcesVehicle, ok := avs[key]; !ok {
			//write no match on Base Vehicle at all NO integrity
			line := []string{cv[0].Part, strconv.FormatFloat(cv[0].Year, 'f', 1, 64), cv[0].Make, cv[0].Model, cv[0].Style, "N", "", "", "", "", "", "", "Base Vehicle (year|make|model) does not exist"}
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
				lines = append(lines, []string{cv.Part, strconv.FormatFloat(cv.Year, 'f', 1, 64), cv.Make, cv.Model, cv.Style, integrity, av.Part, strconv.FormatFloat(av.Year, 'f', 1, 64), av.Make, av.Model, av.Submodel, cons, notes})
			}
		} else {
			//CurtVehicle style != "all" -- AUTOMATICALLY NO INTEGRITY
			integrity := "N" //Aces vehicle DOES NOT have a submodel or config(s)
			notes := "Curt Vehicle Style != all AND an Aces vehicle with either configs or a submodel DOES NOT exists in our ACES data"
			for _, av := range avs {
				if av.Submodel != "" || len(av.Configs) > 0 {
					//Aces vehicle DOES have a submodel or config(s)
					integrity = "N"
					notes = "Curt Vehicle Style != all AND an Aces vehicle with either configs or a submodel DOES exists (may/may not be a Curt style/Aces match) - REVIEW"
				}
			}

			for _, av := range avs {
				// log.Print(av)
				var cons string //stringify configs into comma separated sets
				for i, c := range av.Configs {
					if i > 0 {
						cons += ","
					}
					cons += c.Type + ":" + c.Value
				}
				lines = append(lines, []string{cv.Part, strconv.FormatFloat(cv.Year, 'f', 1, 64), cv.Make, cv.Model, cv.Style, integrity, av.Part, strconv.FormatFloat(av.Year, 'f', 1, 64), av.Make, av.Model, av.Submodel, cons, notes})
			}
		}
	}

	return lines
}
