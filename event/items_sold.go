package event

import (
	"encoding/json"
	"log"

	cmodel "github.com/TerrexTech/go-common-models/model"
	"github.com/TerrexTech/go-mongoutils/mongo"
	"github.com/TerrexTech/rmb-shipment-sales/model"
	"github.com/pkg/errors"
)

func itemsSold(
	itemColl *mongo.Collection,
	saleColl *mongo.Collection,
	event *cmodel.Event,
) error {
	sale := &model.Sale{}
	err := json.Unmarshal(event.Data, sale)
	if err != nil {
		err = errors.Wrap(err, "Error unmarshalling event-data to sale")
		return err
	}

	for _, item := range sale.Items {
		filter := map[string]interface{}{
			"itemID": item.ItemID,
		}
		update := map[string]interface{}{
			"soldWeight": item.SoldWeight,
		}
		result, err := itemColl.FindOne(filter)
		if err != nil {
			err = errors.Wrap(err, "Error getting item from database")
			return err
		}
		dbItem, assertOK := result.(*model.Item)
		if !assertOK {
			err = errors.New("Error asserting dbItem")
			log.Println(err)
		}

		updateResult, err := itemColl.UpdateMany(filter, update)
		if err != nil {
			err = errors.Wrap(err, "Error getting item from database")
			return err
		}
		if updateResult.MatchedCount == 0 || updateResult.ModifiedCount == 0 {
			err = errors.New("error updating item in database")
			return err
		}

		item.DateArrived = dbItem.DateArrived
		item.Lot = dbItem.Lot
		item.Name = dbItem.Name
		item.Price = dbItem.Price
		item.SKU = dbItem.SKU
		item.TotalWeight = dbItem.TotalWeight
		item.UPC = dbItem.UPC
	}

	_, err = saleColl.InsertOne(sale)
	if err != nil {
		err = errors.Wrap(err, "Error inserting sale into database")
		return err
	}
	return nil
}
