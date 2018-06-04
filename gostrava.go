package main

import (
	"encoding/json"
	"flag"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type RecentActivity struct {
	Name                 string  `json:"name"`
	Moving_time          int     `json:"moving_time"`
	Start_date_local     string  `json:"start_date_local"`
	Distance             float64 `json:"distance"`
	Total_elevation_gain float64 `json:"total_elevation_gain"`
}

type TotalActivities struct {
	All_run_totals struct {
		Distance       float64 `json:"distance"`
		Elevation_gain float64 `json:"elevation_gain"`
	} `json:"all_run_totals"`
	Ytd_run_totals struct {
		Distance       float64 `json:"distance"`
		Elevation_gain float64 `json:"elevation_gain"`
	} `json:"ytd_run_totals"`
	Canada  float64 `json: canada`
	Everest float64 `json: everest`
}

var STRAVA_ATHLETE string
var STRAVA_BEARER string

func GetStravaActivity(w http.ResponseWriter, r *http.Request) {
	//func GetStravaActivity() {
	var record []RecentActivity

	params := mux.Vars(r)
	count, err := strconv.Atoi(params["count"])
	if err != nil {
		count = 3
	}

	if count < 0 || count > 30 {
		count = 3
	}

	request, _ := http.NewRequest("GET", "https://www.strava.com/api/v3/athletes/"+STRAVA_ATHLETE+"/activities?per_page="+strconv.Itoa(count), nil)
	request.Header.Set("Authorization", "Bearer "+STRAVA_BEARER)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal("The HTTP request failed with error %s\n", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%#v", record)
	for i := 0; i < len(record); i++ {
		record[i].Distance = math.Ceil(record[i].Distance/10) / 100
	}
	log.Println(record)

        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(record)
}

func GetTotalActivities(w http.ResponseWriter, r *http.Request) {
	//func GetTotalActivities() {
	var record TotalActivities

	request, _ := http.NewRequest("GET", "https://www.strava.com/api/v3/athletes/"+STRAVA_ATHLETE+"/stats", nil)
	request.Header.Set("Authorization", "Bearer "+STRAVA_BEARER)                                                 
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal("The HTTP request failed with error %s\n", err)
	} // else {
	//data, _ := ioutil.ReadAll(resp.Body)

	//const data = `[{"name":"test1", "moving_time":123, "start_date_local":"datetime","distance":123.456,"total_elevation_gain":99.88}]`
	//fmt.Println(string(data))
	//json.Unmarshal([]byte(data), &record)
	//}
	defer resp.Body.Close()

	//fmt.Println(json.NewDecoder(resp.Body))
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Fatal(err)
	}

	record.Canada = math.Ceil(record.All_run_totals.Distance/5514) / 100
	record.Everest = math.Ceil(record.All_run_totals.Elevation_gain / 8848)

	//json.NewEncoder(w).Encode(record)
	//fmt.Printf("%#v", record)
	log.Println(record)
	/*	fmt.Println(record.All_run_totals.Distance)
		fmt.Println(record.All_run_totals.Elevation_gain)
		fmt.Println(record.Ytd_run_totals.Distance)
		fmt.Println(record.Ytd_run_totals.Elevation_gain)
		fmt.Println(record.Canada)*/

        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(record)

	//fmt.Println(json.Marshal(record))
	//if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
	//	log.Println(err)
	//}
}

func main() {
	STRAVA_ATHLETE = os.Getenv("STRAVA_ATHLETE")
	STRAVA_BEARER = os.Getenv("STRAVA_BEARER")

	port := flag.String("port", "9966", "Please specify what port to run gostrava...")
	flag.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/activity", GetStravaActivity).Methods("GET")
	router.HandleFunc("/activity/{count}", GetStravaActivity).Methods("GET")
	router.HandleFunc("/total", GetTotalActivities).Methods("GET")
	log.Println("gostrava is started successfully on port:[" + *port + "]...")
	log.Fatal(http.ListenAndServe(":"+*port, router))

	//	GetStravaActivity()
	//	GetTotalActivities()
}
