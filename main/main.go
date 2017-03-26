package main

import (
	"curioz/api"
	"curioz/cameras"
	"curioz/worker"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
)

type camera struct {
	ID   string
	Name string
}

func main() {

	argSol := flag.Int("sol", -1, "sol - martian day")
	argCam := flag.String("cam", "MAST", "Curiosity rover camera")
	argDest := flag.String("dest", "", "Destination directory")

	flag.Parse()

	sol := *argSol
	cam := *argCam
	dest := *argDest

	if sol == -1 {
		printHelp()
		return
	}

	if cam != cameras.FHAZ &&
		cam != cameras.RHAZ &&
		cam != cameras.MAST &&
		cam != cameras.CHEMCAM &&
		cam != cameras.MAHLI &&
		cam != cameras.MARDI &&
		cam != cameras.NAVCAM &&
		cam != cameras.PANCAM {
		printHelp()
		return
	}

	err := ensureDestinationIsCreated(&dest)
	if err != nil {
		fmt.Println(err)
		return
	}

	for sol <= sol+50 {

		response, err := api.Get(sol, cam)
		if err != nil {
			fmt.Println(err)
			break
		}

		err = worker.Fetch(sol, cam, response, dest)
		if err != nil {
			fmt.Println(err)
		}

		sol++
	}

}

func ensureDestinationIsCreated(dest *string) (err error) {
	if len(*dest) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		*dest = path.Join(cwd, "img_result")

		//create dest dir if it does not exists
		if _, err := os.Stat(*dest); os.IsNotExist(err) {
			err = os.MkdirAll(*dest, 0777)
			if err != nil {
				return err
			}
		}
	} else {
		if _, err := os.Stat(*dest); os.IsNotExist(err) {
			return errors.New("Destination does not exist")
		}
	}

	return nil
}

func printHelp() {
	fmt.Println("Arguments:")
	fmt.Printf("\t--sol : the martian day for which you want to download pictures\n")
	fmt.Printf("\t--cam : mars rover rover camera")
	fmt.Printf("\tCameras:\n")
	fmt.Printf("\t\tFHAZ - Front Hazard Avoidance Camera\n")
	fmt.Printf("\t\tRHAZ - Rear Hazard Avoidance Camera\n")
	fmt.Printf("\t\tMAST - Mast Camera\n")
	fmt.Printf("\t\tCHEMCAM - Chemistry and Camera Complex\n")
	fmt.Printf("\t\tMAHLI - Mars Hand Lens Imager\n")
	fmt.Printf("\t\tMARDI - Mars Descent Imager\n")
	fmt.Printf("\t\tNAVCAM - Navigation Camera\n")
	fmt.Printf("\t\tPANCAM - Panoramic Camera\n")
	fmt.Printf("Usage example:\n")
	fmt.Printf("\tcurioz --sol=1000 --cam=FHAZ")
}
