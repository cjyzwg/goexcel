package model

type Attendancelist struct {
	Userid       int64 `json:"userid"`
	Username string `json:"username"`
	Signedtime 	map[string]*Signedtimelist
	All float64 `all`
}
type Signedtimelist struct {
	//Signedday string `json:signedday`
	//Signedlist struct{
	//	Firstsigned string `firstsigned`
	//	Lastsigned string `lastsigned`
	//	Lackperiod string `lackperiod`
	//	Period string `period`
	//}
	Firstsigned string `firstsigned`
	Lastsigned string `lastsigned`
	Lackperiod string `lackperiod`
	Period string `period`
}

type Tempuser struct {
	Userid   string `json:"userid"`
	Username string `json:"username"`
}
