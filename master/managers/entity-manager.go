package managers

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/vjftw/orchestrate/master/messages"
	"github.com/vjftw/orchestrate/master/models"
)

// EntityManager - Manages the lifecycle of entities
type EntityManager struct {
	// ORM persisters.IPersister `inject:"persister gorm"`
}

// Save - Persist a new or existing Entity. May be stored on multiple storage backends (PGSQL, Redis, etc.)
func (eM EntityManager) Save(entity models.IModel) {
	// eM.ORM.Save(entity)
}

// Validate - Validates a given Entity.
func (eM EntityManager) Validate(entity models.IModel) messages.ValidationMessage {
	result, _ := govalidator.ValidateStruct(entity)

	vM := messages.ValidationMessage{}
	vM.Valid = result

	if !vM.Valid {
		return vM
	}

	if entity, ok := entity.(*models.User); ok {
		fmt.Println("Do user specific validation")
		// eM.ORM.FindInto(entity, "email_address = ?", entity.EmailAddress)
		if entity.ID > 0 {
			vM.Valid = false
		}
	} else {
		fmt.Println("Do usual validation")
		fmt.Printf("%T", entity)
	}

	return vM

}
