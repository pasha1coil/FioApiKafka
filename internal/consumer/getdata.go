package consumer

import (
	"FioapiKafka/internal/models"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type GenderURL struct {
	Gender string `json:"gender"`
}

type AgeURL struct {
	Age int64 `json:"age"`
}

type NationURL struct {
	Country []models.Nation `json:"country"`
}

func GetData(data []byte) (*models.PersonOut, error) {
	var person models.Person
	err := json.Unmarshal(data, &person)
	if err != nil {
		return nil, fmt.Errorf("error unmarshal data to Person struct:%s", err.Error())
	}
	genderURL := "https://api.genderize.io?name=" + person.Name
	agURL := "https://api.agify.io?name=" + person.Name
	nationalURL := "https://api.nationalize.io/?name=" + person.Name

	resp, err := http.Get(genderURL)
	if err != nil {
		return nil, fmt.Errorf("error obtaining gender: %s", err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading gender response: %s", err.Error())
	}
	gender := GenderURL{}
	err = json.Unmarshal(body, &gender)
	if err != nil {
		return nil, fmt.Errorf("error parsing gender response: %s", err.Error())
	}

	resp, err = http.Get(agURL)
	if err != nil {
		return nil, fmt.Errorf("error obtaining age: %s", err.Error())
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading age response: %s", err.Error())
	}
	age := AgeURL{}
	err = json.Unmarshal(body, &age)
	if err != nil {
		return nil, fmt.Errorf("error parsing age response: %s", err.Error())
	}

	resp, err = http.Get(nationalURL)
	if err != nil {
		return nil, fmt.Errorf("error obtaining nationality: %s", err.Error())
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading nationality response: %s", err.Error())
	}
	nation := NationURL{}
	err = json.Unmarshal(body, &nation)
	if err != nil {
		return nil, fmt.Errorf("error parsing nationality response: %s", err.Error())
	}
	if gender.Gender == "" || age.Age == 0 || len(nation.Country) == 0 {
		log.Error("Failed FIO data, no such data - Gender/Age/Country")
		return nil, fmt.Errorf("failed FIO")
	}

	alldata := models.PersonOut{
		Name:        person.Name,
		Surname:     person.Surname,
		Patronymic:  person.Patronymic,
		Gender:      gender.Gender,
		Age:         age.Age,
		Nationality: nation.Country,
	}

	return &alldata, nil
}
