package models

type Hotel struct {
	RoomCount			int					`json:"room_count"`
	CheckIn				string				`json:"check_in"`
	CheckOut			string				`json:"check_out"`
	DistanceToCenter	int					`json:"distance_to_center"`
	Address				string				`json:"address"`
	Name				string				`json:"name"`
	Location_id			int					`json:"location_id"`
	PhotosIds			[]string			`json:"photos_ids"`
	Minprice			int					`json:"minprice"`
	Property_type		string				`json:"property_type"`
	Rating				float32				`json:"rating"`
	Stars				int					`json:"stars"`
	YearOpened			int					`json:"year_opened"`
	YearRenovated		int					`json:"year_renovated"`
	LocationsIds		[]int				`json:"locations_ids"`
	Popularity2			int					`json:"popularity2"`
	Reviews_count		int					`json:"reviews_count"`
	Chain				string				`json:"chain"`
	Id					int					`json:"id"`
	Popularity			int					`json:"popularity"`
	LatinName			string				`json:"latin_name"`
	IsRentals			bool				`json:"is_rentals"`
	HasRentals			bool				`json:"has_rentals"`
	PhotosCount			int					`json:"photos_count"`
	LatinFullFame		string				`json:"latin_full_name"`
	FullName			string				`json:"full_name"`
	PropertyTypeId		int					`json:"property_type_id"`
	PriceGroup			int					`json:"price_group"`
	MedianMinprice		int					`json:"median_minprice"`
}