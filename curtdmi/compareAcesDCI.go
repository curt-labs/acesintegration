package curtdmi

// import (
// 	"github.com/curt-labs/acesintegration/curtaces"

// 	"encoding/csv"
// 	"os"
// 	"strconv"
// 	"strings"
// )

// //CheckIntegrity
// func ProcessAcesToDci(avs map[string][]curtaces.AcesVehicleApplication, dvs map[string][]DmiVehicleApplication) error {
// 	//set up file
// 	f, err := os.Create("integrityCurtDCI.csv")
// 	if err != nil {
// 		return err
// 	}
// 	writer := csv.NewWriter(f)

// 	//write header
// 	header := []string{"Part", "Year", "Make", "Model", "Style", "Integrated", "Year", "Make", "Model", "Submodel", "Configs", "Notes"}
// 	err = writer.Write(header)
// 	if err != nil {
// 		return err
// 	}

// 	//check integrity & write lines
// 	for key, cv := range avs {
// 		if matchingDciVehicle, ok := dvs[key]; !ok {
// 			//write no match on Base Vehicle at all NO integrity
// 			line := []string{cv[0].Part, strconv.FormatFloat(cv[0].Year, 'f', 1, 64), cv[0].Make, cv[0].Model, cv[0].Style, "N", "", "", "", "", "", "Base Vehicle (year|make|model) does not exist"}
// 			err = writer.Write(line)
// 			if err != nil {
// 				return err
// 			}
// 			writer.Flush()
// 		} else {
// 			lines := determineIntegrity(cv, matchingDciVehicle)
// 			err = writer.WriteAll(lines)
// 			if err != nil {
// 				return err
// 			}
// 			writer.Flush()
// 		}
// 	}
// 	return nil
// }

// func determineIntegrity(avs []curtaces.AcesVehicleApplication, dvs []DmiVehicleApplication) [][]string {
// 	var lines [][]string

// 	for _, cv := range avs {
// 		//CurtVehicle style == "all"
// 		if strings.ToLower(cv.Style) == "all" {
// 			integrity := "N"
// 			notes := "Curt Vehicle Style = all BUT an DCI non-config/non-submodel DOES NOT exists (need base)"
// 			for _, av := range dvs {
// 				if av.Submodel == "" && len(av.Configs) == 0 {
// 					integrity = "Y"
// 					notes = "Curt Vehicle Style = all AND an DCI non-config/non-submodel exists"
// 				}
// 			}
// 			for _, av := range dvs {
// 				var cons string //stringify configs into comma separated sets
// 				for i, c := range av.Configs {
// 					if i > 0 {
// 						cons += ","
// 					}
// 					cons += c.Type + ":" + c.Value
// 				}
// 				lines = append(lines, []string{cv.Part, strconv.FormatFloat(cv.Year, 'f', 1, 64), cv.Make, cv.Model, cv.Style, integrity, strconv.FormatFloat(av.Year, 'f', 1, 64), av.Make, av.Model, av.Submodel, cons, notes})
// 			}
// 		} else {
// 			//CurtVehicle style != "all"
// 			integrity := "N" //DCI vehicle DOES NOT have a submodel or config(s)
// 			notes := "Curt Vehicle Style != all AND an DCI vehicle with either configs or a submodel DOES NOT exists"
// 			for _, av := range dvs {
// 				if av.Submodel != "" || len(av.Configs) > 0 {
// 					//DCI vehicle DOES have a submodel or config(s)
// 					integrity = "Y"
// 					notes = "Curt Vehicle Style != all AND an DCI vehicle with either configs or a submodel DOES exists"
// 				}
// 			}

// 			for _, av := range dvs {
// 				var cons string //stringify configs into comma separated sets
// 				for i, c := range av.Configs {
// 					if i > 0 {
// 						cons += ","
// 					}
// 					cons += c.Type + ":" + c.Value
// 				}
// 				lines = append(lines, []string{cv.Part, strconv.FormatFloat(cv.Year, 'f', 1, 64), cv.Make, cv.Model, cv.Style, integrity, strconv.FormatFloat(av.Year, 'f', 1, 64), av.Make, av.Model, av.Submodel, cons, notes})
// 			}
// 		}
// 	}

// 	return lines
// }
