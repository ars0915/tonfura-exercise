package usecase

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ars0915/tonfura-exercise/entity"
)

func (h TaskHandler) ListFlights(ctx context.Context, param entity.ListFlightParam) (flights []entity.Flight, count int64, err error) {
	if flights, err = h.db.ListFlights(param); err != nil {
		return flights, 0, errors.Wrap(err, "list flights")
	}

	count, err = h.db.GetFlightsCount(param)
	if err != nil {
		return nil, 0, errors.Wrap(err, "count flights")
	}

	return
}
