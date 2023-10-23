package graphqlhandler

import (
	"FioapiKafka/internal/services"
	"github.com/graphql-go/graphql"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type graphqlHandler struct {
	schema *graphql.Schema
}

func NewGraphQLHandler(service *services.Service) *graphqlHandler {
	gh := &graphqlHandler{}
	gh.schema = gh.createSchema(service)
	return gh
}

func (gh *graphqlHandler) ExecuteQuery(query string) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        *gh.schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		log.Errorf("Unexpected errors inside ExecuteQuery: %s", result.Errors)
	}
	return result
}

func (gh *graphqlHandler) createSchema(service *services.Service) *graphql.Schema {

	nationalityType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Nationality",
			Fields: graphql.Fields{
				"countryId": &graphql.Field{
					Type: graphql.String,
				},
				"probability": &graphql.Field{
					Type: graphql.Float,
				},
			},
		})

	personType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Person",
			Fields: graphql.Fields{
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"surname": &graphql.Field{
					Type: graphql.String,
				},
				"patronymic": &graphql.Field{
					Type: graphql.String,
				},
				"gender": &graphql.Field{
					Type: graphql.String,
				},
				"age": &graphql.Field{
					Type: graphql.Int,
				},
				"nationality": &graphql.Field{
					Type: graphql.NewList(nationalityType),
				},
			},
		})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"peopleAll": &graphql.Field{
				Type:        graphql.NewList(personType),
				Description: "Get list of people",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					people, err := service.GetPeople()
					return people, err
				},
			},
			"personById": &graphql.Field{
				Type:        personType,
				Description: "Get person by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(string)
					if ok {
						person, err := service.GetPersonByID(id)
						return person, err
					}
					return nil, nil
				},
			},
			"personsByName": &graphql.Field{
				Type:        graphql.NewList(personType),
				Description: "Get list of people by name",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					name, _ := p.Args["name"].(string)
					people, err := service.GetPersonsByName(name)
					return people, err
				},
			},
			"personsByAge": &graphql.Field{
				Type:        graphql.NewList(personType),
				Description: "Get list of people by age",
				Args: graphql.FieldConfigArgument{
					"age": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					age, _ := p.Args["age"].(int)
					people, err := service.GetPersonsByAge(strconv.Itoa(age))
					return people, err
				},
			},
			"personsByGender": &graphql.Field{
				Type:        graphql.NewList(personType),
				Description: "Get list of people by gender",
				Args: graphql.FieldConfigArgument{
					"gender": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					gender, _ := p.Args["gender"].(string)
					people, err := service.GetPersonsByGender(gender)
					return people, err
				},
			},
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"deletePerson": &graphql.Field{
				Type:        graphql.String,
				Description: "Delete a person by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(string)
					service.DeletePerson(id)
					return "Person deleted successfully.", nil
				},
			},
		},
	})
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	return &schema
}
