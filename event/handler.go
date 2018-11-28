package event

import (
	"fmt"

	"github.com/TerrexTech/go-common-models/model"
	"github.com/TerrexTech/go-mongoutils/mongo"
	"github.com/pkg/errors"
)

// Handle handles the provided event.
func Handle(
	itemColl *mongo.Collection,
	saleColl *mongo.Collection,
	eosToken string,
	event *model.Event,
) error {
	if itemColl == nil {
		return errors.New("itemColl cannot be nil")
	}
	if saleColl == nil {
		return errors.New("saleColl cannot be nil")
	}
	if event == nil || event.Action == eosToken {
		return nil
	}

	switch event.Action {
	case "ItemsSold":
		err := itemsSold(itemColl, saleColl, event)
		if err != nil {
			err = errors.Wrap(err, "Error processing ItemsSold-event")
			return err
		}
		return nil

	default:
		return fmt.Errorf("unregistered Action: %s", event.Action)
	}
}
