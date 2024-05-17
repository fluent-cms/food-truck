package models

type Facility struct {
	LocationID              string  `json:"locationID"`
	Applicant               string  `json:"applicant"`
	FacilityType            string  `json:"facilityType"`
	CNN                     string  `json:"cnn"`
	LocationDescription     string  `json:"locationDescription"`
	Address                 string  `json:"address"`
	BlockLot                string  `json:"blockLot"`
	Block                   string  `json:"block"`
	Lot                     string  `json:"lot"`
	Permit                  string  `json:"permit"`
	Status                  string  `json:"status"`
	FoodItems               string  `json:"foodItems"`
	X                       float64 `json:"x"`
	Y                       float64 `json:"y"`
	Latitude                float64 `json:"latitude"`
	Longitude               float64 `json:"longitude"`
	Schedule                string  `json:"schedule"`
	DaysHours               string  `json:"daysHours"`
	NOISent                 string  `json:"NOISent"`
	Approved                string  `json:"approved"`
	Received                string  `json:"received"`
	PriorPermit             string  `json:"priorPermit"`
	ExpirationDate          string  `json:"expirationDate"`
	Location                string  `json:"location"`
	FirePreventionDistricts string  `json:"firePreventionDistricts"`
	PoliceDistricts         string  `json:"policeDistricts"`
	SupervisorDistricts     string  `json:"supervisorDistricts"`
	ZipCodes                string  `json:"zipCodes"`
	NeighborhoodsOld        string  `json:"neighborhoodsOld"`
}

func GetFacilityLocation(f Facility) (float64, float64) {
	return f.Latitude, f.Longitude
}

func GetFacilityKey(f Facility) string {
	return f.LocationID
}

func GetFacilityScore(f Facility) float64 {
	return 0
}
