package handler

import (
	"FioapiKafka/internal/models"
	"FioapiKafka/internal/services"
	"context"
	"encoding/json"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type Handlers struct {
	service *services.Service
}

func NewHandlers(service *services.Service) *Handlers {
	return &Handlers{service: service}
}

// GetPersons возвращает всех людей
func (h *Handlers) GetPersons(c *fiber.Ctx) error {
	log.Infoln("Start Handler - GetPersons")
	people, err := h.service.GetPeople()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(people)
}

// GetPersonByID возвращает человека по ID
func (h *Handlers) GetPersonByID(c *fiber.Ctx) error {
	log.Infoln("Start Handler - GetPersonByID")
	id := c.Params("id")
	person, err := h.service.GetPersonByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(person)
}

// GetPersonsByName возвращает людей по имени
func (h *Handlers) GetPersonsByName(c *fiber.Ctx) error {
	log.Infoln("Start Handler - GetPersonsByName")
	name := c.Params("name")
	people, err := h.service.GetPersonsByName(name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if len(people) == 0 {
		return c.Status(fiber.StatusOK).JSON("There are no people with that name")
	}
	return c.Status(fiber.StatusOK).JSON(people)
}

// GetPersonsByAge возвращает людей по возрасту
func (h *Handlers) GetPersonsByAge(c *fiber.Ctx) error {
	log.Infoln("Start Handler - GetPersonsByAge")
	age := c.Params("age")
	people, err := h.service.GetPersonsByAge(age)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if len(people) == 0 {
		return c.Status(fiber.StatusOK).JSON("There are no people that age")
	}
	return c.Status(fiber.StatusOK).JSON(people)
}

// GetPersonsByGender возвращает людей по полу
func (h *Handlers) GetPersonsByGender(c *fiber.Ctx) error {
	log.Infoln("Start Handler - GetPersonsByGender")
	gender := c.Params("gender")
	people, err := h.service.GetPersonsByGender(gender)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if len(people) == 0 {
		return c.Status(fiber.StatusOK).JSON("There are no people with this gender")
	}
	return c.Status(fiber.StatusOK).JSON(people)
}

// CreatePerson создает нового человека
func (h *Handlers) CreatePerson(c *fiber.Ctx) error {
	log.Infoln("Start Handler - CreatePerson")
	person := new(models.Person)

	if err := c.BodyParser(person); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	personJSON, _ := json.Marshal(person)

	writer := &kafka.Writer{
		Addr:     kafka.TCP(os.Getenv("KAFKA_BROKER")),
		Topic:    os.Getenv("KAFKA_FIO_TOPIC"),
		Balancer: &kafka.LeastBytes{},
	}
	log.Infoln("Recording message to KAFKA FIO_TOPIC")
	err := writer.WriteMessages(context.Background(), kafka.Message{Value: personJSON})
	if err != nil {
		log.Errorf("Error writing a message to KAFKA FIO_TOPIC: %s", err.Error())
	}
	writer.Close()
	log.Infoln("Writing the message to KAFKA FIO_TOPIC was successful")
	return c.Status(fiber.StatusCreated).JSON("OK")
}

// UpdatePerson обновляет человека по ID
func (h *Handlers) UpdatePerson(c *fiber.Ctx) error {
	log.Infoln("Start Handler - UpdatePerson")
	id := c.Params("id")
	person := new(models.PersonOut)

	if err := c.BodyParser(person); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err := h.service.UpdatePerson(id, person)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(person)
}

// DeletePerson удаляет человека по ID
func (h *Handlers) DeletePerson(c *fiber.Ctx) error {
	log.Infoln("Start Handler - DeletePerson")
	id := c.Params("id")

	err := h.service.DeletePerson(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).SendString("Person deleted successfully")
}
