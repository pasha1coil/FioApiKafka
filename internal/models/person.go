package models

type PersonOut struct {
	Name        string
	Surname     string
	Patronymic  string
	Gender      string
	Age         int64
	Nationality []Nation
}

type Nation struct {
	CountryId   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type Person struct {
	Name       string `json:"name" validate:"required"`
	Surname    string `json:"surname" validate:"required"`
	Patronymic string `json:"patronymic"`
}
