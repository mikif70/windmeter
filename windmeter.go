// windmeter
package main

import (
	"gopkg.in/mgo.v2"	
	"fmt"
	"strings"
	"net/http"
//	"io/ioutil"
	"os"
	"bufio"
)

var (
	mongodb  = "mongodb://10.39.81.85:27018"
)

const dateForm = "08/06/2015 17:09"

func main() {
	
	session, err := mgo.Dial(mongodb)
	if err != nil {
		panic(err)
	}
	defer session.Close()
//	db := session.DB("wind").C("kitezone")
	
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
			data := strings.Split(retval, " ")
			fmt.Println(data)
		}
	}	
}
