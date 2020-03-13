package manage

import (
	"context"

	"github.com/leeif/go-web-template/datatype/error"
	"github.com/leeif/go-web-template/models"
)

func (m *Manager) Healthcheck() *error.Error {
	_, err := models.Dummies().All(context.Background(), m.db)
	if err != nil {
		return error.ServerError.Wrapper(err)
	}
	return nil
}
