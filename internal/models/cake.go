package models

type Cake struct {
	ID      int
	Status  string 
	X, Y    float64 
	TargetX float64
	TargetY float64
	Alpha   float64
	InOven  bool
}
