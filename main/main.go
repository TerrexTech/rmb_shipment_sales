package main

import (
	"log"
	"os"

	"github.com/TerrexTech/go-mongoutils/mongo"

	"github.com/TerrexTech/rmb-shipment-sales/model"

	"github.com/TerrexTech/go-commonutils/commonutil"
	"github.com/TerrexTech/go-eventspoll/poll"
	"github.com/TerrexTech/rmb-shipment-sales/connutil"
	"github.com/TerrexTech/rmb-shipment-sales/event"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func validateEnv() error {
	log.Println("Reading environment file")
	err := godotenv.Load("./.env")
	if err != nil {
		err = errors.Wrap(err,
			".env file not found, env-vars will be read as set in environment",
		)
		log.Println(err)
	}

	missingVar, err := commonutil.ValidateEnv(
		"SERVICE_NAME",

		"KAFKA_BROKERS",

		"KAFKA_CONSUMER_EVENT_GROUP",
		"KAFKA_CONSUMER_EVENT_QUERY_GROUP",

		"KAFKA_CONSUMER_EVENT_TOPIC",
		"KAFKA_CONSUMER_EVENT_QUERY_TOPIC",

		"KAFKA_PRODUCER_EVENT_QUERY_TOPIC",

		"KAFKA_END_OF_STREAM_TOKEN",

		"MONGO_HOSTS",
		"MONGO_DATABASE",
		"MONGO_ITEMS_COLLECTION",
		"MONGO_SALES_COLLECTION",
		"MONGO_META_COLLECTION",
		"MONGO_CONNECTION_TIMEOUT_MS",
	)

	if err != nil {
		err = errors.Wrapf(err, "Env-var %s is required for testing, but is not set", missingVar)
		return err
	}
	return nil
}

func saleColl() (*mongo.Collection, error) {
	conn, err := connutil.GetMongoConn()
	if err != nil {
		err = errors.Wrap(err, "Error getting shipment-collection")
		return nil, err
	}

	database := os.Getenv("MONGO_DATABASE")

	indexConfigs := []mongo.IndexConfig{
		mongo.IndexConfig{
			ColumnConfig: []mongo.IndexColumnConfig{
				mongo.IndexColumnConfig{
					Name: "saleID",
				},
			},
			IsUnique: true,
			Name:     "saleID_index",
		},
		mongo.IndexConfig{
			ColumnConfig: []mongo.IndexColumnConfig{
				mongo.IndexColumnConfig{
					Name: "timestamp",
				},
			},
			IsUnique: true,
			Name:     "timestamp_index",
		},
	}

	collection := os.Getenv("MONGO_SALES_COLLECTION")
	c := &mongo.Collection{
		Connection:   conn,
		Database:     database,
		Name:         collection,
		SchemaStruct: &model.Sale{},
		Indexes:      indexConfigs,
	}
	saleColl, err := mongo.EnsureCollection(c)
	if err != nil {
		err = errors.Wrap(err, "Error creating MongoCollection")
		return nil, err
	}
	return saleColl, nil
}

func main() {
	err := validateEnv()
	if err != nil {
		log.Fatalln(err)
	}

	kc, err := connutil.LoadKafkaConfig(model.AggregateID)
	if err != nil {
		err = errors.Wrap(err, "Error in KafkaConfig")
		log.Fatalln(err)
	}
	mc, err := connutil.LoadMongoConfig(model.AggregateID)
	if err != nil {
		err = errors.Wrap(err, "Error in MongoConfig")
		log.Fatalln(err)
	}
	ioConfig := poll.IOConfig{
		KafkaConfig: *kc,
		MongoConfig: *mc,
	}

	eventPoll, err := poll.Init(ioConfig)
	if err != nil {
		err = errors.Wrap(err, "Error creating EventPoll service")
		log.Fatalln(err)
	}

	saleColl, err := saleColl()
	if err != nil {
		err = errors.Wrap(err, "Error creating Sales collection")
		log.Fatalln(err)
	}

	eosToken := os.Getenv("KAFKA_END_OF_STREAM_TOKEN")
	for {
		select {
		case <-eventPoll.Context().Done():
			err = errors.New("service-context closed")
			log.Fatalln(err)

		case eventResp := <-eventPoll.Events():
			go func(eventResp *poll.EventResponse) {
				if eventResp == nil {
					return
				}
				err := eventResp.Error
				if err != nil {
					err = errors.Wrap(err, "Error in Query-EventResponse")
					log.Println(err)
					return
				}

				err = event.Handle(
					mc.AggCollection,
					saleColl,
					eosToken,
					&eventResp.Event,
				)
				if err != nil {
					err = errors.Wrap(err, "Error handling event")
					log.Println(err)
				}
			}(eventResp)
		}
	}
}
