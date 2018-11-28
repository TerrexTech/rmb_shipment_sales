package connutil

import (
	"log"
	"os"
	"strconv"

	"github.com/TerrexTech/go-commonutils/commonutil"
	"github.com/TerrexTech/go-eventspoll/poll"
	"github.com/TerrexTech/go-mongoutils/mongo"
	"github.com/TerrexTech/rmb-shipment-sales/model"
	"github.com/pkg/errors"
)

// LoadMongoConfig is a convenient function to load Mongoconfig for EventsPoll.
func LoadMongoConfig(aggID int8) (*poll.MongoConfig, error) {
	conn, err := GetMongoConn()
	if err != nil {
		log.Fatalln(err)
	}

	database := os.Getenv("MONGO_DATABASE")
	aggCollection := os.Getenv("MONGO_ITEMS_COLLECTION")
	metaCollection := os.Getenv("MONGO_META_COLLECTION")

	aggMongoCollection, err := createMongoCollection(conn, database, aggCollection)
	if err != nil {
		err = errors.Wrap(err, "Error creating MongoCollection")
		return nil, err
	}

	return &poll.MongoConfig{
		AggregateID:        aggID,
		AggCollection:      aggMongoCollection,
		Connection:         conn,
		MetaDatabaseName:   database,
		MetaCollectionName: metaCollection,
	}, nil
}

// GetMongoConn reutrns a Mongo=client using env-vars.
func GetMongoConn() (*mongo.ConnectionConfig, error) {
	hosts := *commonutil.ParseHosts(
		os.Getenv("MONGO_HOSTS"),
	)

	itemname := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")
	connTimeoutStr := os.Getenv("MONGO_CONNECTION_TIMEOUT_MS")
	connTimeout, err := strconv.Atoi(connTimeoutStr)
	if err != nil {
		err = errors.Wrap(err, "Error converting MONGO_CONNECTION_TIMEOUT_MS to integer")
		log.Println(err)
		log.Println("A defalt value of 3000 will be used for MONGO_CONNECTION_TIMEOUT_MS")
		connTimeout = 3000
	}

	mongoConfig := mongo.ClientConfig{
		Hosts:               hosts,
		Username:            itemname,
		Password:            password,
		TimeoutMilliseconds: uint32(connTimeout),
	}

	// MongoDB Client
	client, err := mongo.NewClient(mongoConfig)
	if err != nil {
		err = errors.Wrap(err, "Error creating MongoClient")
		return nil, err
	}

	resTimeoutStr := os.Getenv("MONGO_CONNECTION_TIMEOUT_MS")
	resTimeout, err := strconv.Atoi(resTimeoutStr)
	if err != nil {
		err = errors.Wrap(err, "Error converting MONGO_RESOURCE_TIMEOUT_MS to integer")
		log.Println(err)
		log.Println("A defalt value of 5000 will be used for MONGO_RESOURCE_TIMEOUT_MS")
		resTimeout = 5000
	}
	conn := &mongo.ConnectionConfig{
		Client:  client,
		Timeout: uint32(resTimeout),
	}

	return conn, nil
}

func createMongoCollection(
	conn *mongo.ConnectionConfig, db string, coll string,
) (*mongo.Collection, error) {
	// Index Configuration
	indexConfigs := []mongo.IndexConfig{
		mongo.IndexConfig{
			ColumnConfig: []mongo.IndexColumnConfig{
				mongo.IndexColumnConfig{
					Name: "itemID",
				},
			},
			IsUnique: true,
			Name:     "itemID_index",
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

	// Create New Collection
	c := &mongo.Collection{
		Connection:   conn,
		Database:     db,
		Name:         coll,
		SchemaStruct: &model.Item{},
		Indexes:      indexConfigs,
	}
	collection, err := mongo.EnsureCollection(c)
	if err != nil {
		err = errors.Wrap(err, "Error creating MongoCollection")
		return nil, err
	}
	return collection, nil
}
