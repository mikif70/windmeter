// windmeter
package main

import (
	"fmt"
	"github.com/mikif70/Pid"
	"gopkg.in/mgo.v2"
	"net/http"
	"strings"
	//	"io/ioutil"
	"bufio"
	"os"
	"regexp"
	"time"
)

var (
	mongodb = "mongodb://10.39.81.85:27018"
)

type kitezone struct {
	Time time.Time `json:"time"`
	Data string    `json:"data"`
	Ora  string    `json:"ora"`
	Vel  string    `json:"vel"`
	Dir  string    `json:"dir"`
	Temp string    `json:"temp"`
}

func parseKitezone(session *mgo.Session) {
	dateForm := "02/01/2006 15:04"
	db := session.DB("wind").C("kitezone")
	resp, err := http.Get("http://www.kitezone.it/public/dati_short.txt")
	if err != nil {
		fmt.Println("HTTP Error: ", err.Error())
		os.Exit(-1)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		retval := scanner.Text()
		if len(retval) > 0 {
			data := strings.Split(retval, "\t")
			tm, _ := time.Parse(dateForm, data[0]+" "+data[1])
			kz := kitezone{
				Time: tm,
				Data: data[0],
				Ora:  data[1],
				Dir:  data[3],
				Vel:  data[4],
				Temp: data[7],
			}
			db.Insert(kz)
		}
	}
}

func parsePortodagumu(session *mgo.Session) {
	dateForm := "02/01/06 15.04"
	db := session.DB("wind").C("portodagumu")
	resp, err := http.Get("http://www.portodagumu.it/public/davis/dati_short.txt")
	if err != nil {
		fmt.Println("HTTP Error: ", err.Error())
		os.Exit(-1)
	}
	defer resp.Body.Close()

	regex, _ := regexp.Compile(`^\d{2}\/\d{2}\/\d{2}.*`)

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		retval := scanner.Text()
		if len(retval) > 0 && regex.MatchString(retval) {
			data := strings.Split(retval, "\t")
			do := strings.Split(data[0], " ")
			tm, _ := time.Parse(dateForm, data[0])
			kz := kitezone{
				Time: tm,
				Data: do[0],
				Ora:  do[1],
				Dir:  data[7],
				Vel:  data[8],
				Temp: data[3],
			}
			db.Insert(kz)
		}
	}
}

func main() {

	p := pid.New()

	p.Write()
	defer p.Remove()

	session, err := mgo.Dial(mongodb)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	parseKitezone(session)
	parsePortodagumu(session)
}
