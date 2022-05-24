package todoapp

type Type string

type Location struct {
	// - name: location_type
	//  in: location_type
	//  description: location type
	//  schema:
	//  type: object
	//  required: true
	LocationType Type `bson:"location_type",json:"location_type"`

	// - name: coordinates
	//  in: coordinates
	//  description: coordinate of location
	//  schema:
	//  type: object
	//  required: true
	Coordinates Coordinate `bson:"coordinates",json:"coordinates"`
}

type Coordinate struct {
	// latitude of location
	// in: float64
	Latitude float64 `bson:"latitude",json:"latitudxse",csv:"Latitude`

	// longtitude of location
	// in: float64
	Longtitude float64 `bson:"longtitude",json:"longtitude",csv:"Latitude`
}
