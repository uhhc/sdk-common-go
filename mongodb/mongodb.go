package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/uhhc/sdk-common-go/log"
)

// MongoClient represents the struct of mongodb client
type MongoClient struct {
	dbname      string
	client      *mongo.Client
	logger      log.Logger
	loggerClone log.Logger
}

// NewMongoClient to get mongodb instance
func NewMongoClient(logger log.Logger) (*MongoClient, error) {
	loggerClone := logger
	logger.SugaredLogger = logger.With("method", "NewMongoClient")

	viper.AutomaticEnv()

	// Get config
	user := viper.GetString("MONGODB_USER")
	password := viper.GetString("MONGODB_PASSWORD")
	host := viper.GetString("MONGODB_HOST")
	port := viper.GetString("MONGODB_PORT")
	ssl := viper.GetString("MONGODB_SSL")
	timeout := viper.GetInt64("MONGODB_CONN_TIMEOUT")

	// Set default value
	if ssl == "" {
		ssl = "false"
	}
	if timeout == 0 {
		timeout = 10
	}

	// Set URI
	// mongodb://[username:password@]host1[:port1][,...hostN[:portN]][/[database][?options]]
	// See https://docs.mongodb.com/manual/reference/connection-string/
	var uri string
	if user != "" && password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s/?ssl=%s", user, password, host, port, ssl)
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s/?ssl=%s", host, port, ssl)
	}
	// Do not log the password
	logger.Debugw("", "user", user, "host", host, "port", port, "ssl", ssl)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		logger.Errorw("connect to mongodb error", "error", err, "uri", uri)
		return nil, err
	}

	// Check the connection
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		logger.Errorw("ping mongodb error", "error", err, "uri", uri)
		return nil, err
	}

	return &MongoClient{
		client:      client,
		logger:      logger,
		loggerClone: loggerClone,
	}, nil
}

// SetDatabase to set default database
func (mc *MongoClient) SetDatabase(dbname string) *MongoClient {
	mc.dbname = dbname
	return mc
}

func (mc *MongoClient) getDbHandler() *mongo.Database {
	if mc.dbname == "" {
		mc.dbname = viper.GetString("MONGODB_DBNAME")
	}
	if mc.dbname == "" {
		mc.logger.Fatalw("you have not set mongodb database")
	}
	return mc.client.Database(mc.dbname)
}

func (mc *MongoClient) getContext() context.Context {
	timeout := viper.GetInt64("MONGODB_OP_TIMEOUT")
	if timeout == 0 {
		timeout = 10
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	return ctx
}

// GetCollectionHandler to get a collection handler
func (mc *MongoClient) GetCollectionHandler(name string) *mongo.Collection {
	if name == "" {
		mc.logger.Fatalw("you have not set mongodb collection name")
	}
	return mc.getDbHandler().Collection(name)
}

// InsertOne to insert one document into mongodb collection
// See https://godoc.org/go.mongodb.org/mongo-driver/mongo#Collection.InsertOne
// Example:
//
// 		collectionName := "info_data"
//
// 		data1 := bson.D{
// 			{"test_id", "id001"},
// 			{"name", "name001"},
// 		}
// 		// Insert bson.D data
// 		res, err := mongo.InsertOne(collectionName, data1)
//
// 		type test struct {
// 			TestId string `bson:"test_id"`
// 			Name string `bson:"name"`
// 		}
// 		data2 := test{
// 			TestId: "id002",
// 			Name: "name002",
// 		}
// 		// Insert struct data
// 		res, err := mongo.InsertOne(collectionName, data2)
//
func (mc *MongoClient) InsertOne(collectionName string, data interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	mc.logger = mc.loggerClone
	mc.logger.SugaredLogger = mc.logger.With("method", "InsertOne")

	collection := mc.GetCollectionHandler(collectionName)
	res, err := collection.InsertOne(mc.getContext(), data, opts...)
	if err != nil {
		mc.logger.Errorw("insert one data error", "error", err)
		return res, err
	}
	return res, nil
}

// InsertMany to insert many documents into mongodb collection
// See https://godoc.org/go.mongodb.org/mongo-driver/mongo#Collection.InsertMany
// Example:
//
// 		data1 := []interface{}{
// 			bson.D{{"test_id", "id001"}, {"name", "name001"}},
// 			bson.D{{"test_id", "id002"}, {"name", "name002"}},
// 		}
// 		// Insert bson.D data
// 		res, err := mongo.InsertMany(collectionName, data1)
//
// 		type test struct {
// 			TestId string `bson:"test_id"`
// 			Name string `bson:"name"`
// 		}
// 		data2 := []interface{}{
// 			test{
// 				TestId: "id003",
// 				Name: "name003",
// 			},
// 			test{
// 				TestId: "id004",
// 				Name: "name004",
// 			},
// 		}
// 		// Insert struct data
// 		res, err := mongo.InsertMany(collectionName, data2)
//
func (mc *MongoClient) InsertMany(collectionName string, data []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	mc.logger = mc.loggerClone
	mc.logger.SugaredLogger = mc.logger.With("method", "InsertMany")

	collection := mc.GetCollectionHandler(collectionName)
	res, err := collection.InsertMany(mc.getContext(), data, opts...)
	if err != nil {
		mc.logger.Errorw("insert many data error", "error", err)
		return res, err
	}
	return res, nil
}

// GetOne to get one document
// See https://godoc.org/go.mongodb.org/mongo-driver/mongo#Collection.FindOne
// Example:
//
// 		type test struct {
// 			TestId string `bson:"test_id"`
// 			Name string `bson:"name"`
// 		}
// 		var data test
// 		err := mongo.GetOne(collectionName, bson.D{{"name", "name001"}}, &data)
// 		fmt.Printf("data: %+v", data)
//
func (mc *MongoClient) GetOne(collectionName string, filter interface{}, result interface{}, opts ...*options.FindOneOptions) error {
	mc.logger = mc.loggerClone
	mc.logger.SugaredLogger = mc.logger.With("method", "GetOne")

	collection := mc.GetCollectionHandler(collectionName)
	err := collection.FindOne(mc.getContext(), filter, opts...).Decode(result)
	if err != nil {
		mc.logger.Errorw("get one data error", "error", err)
		return err
	}
	return nil
}

// GetManyWithBsonFmt to get many documents with bson.M format
// See https://godoc.org/go.mongodb.org/mongo-driver/mongo#Collection.Find
// Example:
//
// 		type test struct {
// 			TestId string `bson:"test_id"`
// 			Name string `bson:"name"`
// 		}
// 		var (
// 			item test
// 			results []test
// 		)
// 		opts := options.Find().SetSkip(1).SetLimit(2)
// 		res, err := mongo.GetManyWithBsonFmt(collectionName, bson.D{{"name", "name001"}}, opts)
//		// convert bson to struct
// 		for _, v := range *res {
// 			bsonBytes, _ := bson.Marshal(v)
// 			_ = bson.Unmarshal(bsonBytes, &item)
// 			results = append(results, item)
// 		}
// 		fmt.Printf("results: %+v\n", results)
//
func (mc *MongoClient) GetManyWithBsonFmt(collectionName string, filter interface{}, opts ...*options.FindOptions) (*[]bson.M, error) {
	mc.logger = mc.loggerClone
	mc.logger.SugaredLogger = mc.logger.With("method", "GetManyWithBsonFmt")

	ctx := mc.getContext()
	collection := mc.GetCollectionHandler(collectionName)
	cur, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		mc.logger.Errorw("find collection data error", "error", err)
		return nil, err
	}
	var results []bson.M
	if err := cur.All(ctx, &results); err != nil {
		mc.logger.Errorw("get many data error", "error", err)
		return nil, err
	}

	return &results, nil
}

// UpdateOne to update one document
// See https://godoc.org/go.mongodb.org/mongo-driver/mongo#Collection.UpdateOne
// Example:
//
// 		type test struct {
// 			TestId string `bson:"test_id"`
// 			Name string `bson:"name"`
// 		}
// 		filter := bson.D{{"test_id", "id001"}}
// 		update := bson.D{{"$set", bson.D{{"name", "newname001"}}}}
// 		res, err := mongo.UpdateOne(collectionName, filter, update)
// 		fmt.Printf("res: %+v\n", res)
//
func (mc *MongoClient) UpdateOne(collectionName string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	mc.logger = mc.loggerClone
	mc.logger.SugaredLogger = mc.logger.With("method", "UpdateOne")

	collection := mc.GetCollectionHandler(collectionName)
	res, err := collection.UpdateOne(mc.getContext(), filter, update, opts...)
	if err != nil {
		mc.logger.Errorw("update one data error", "error", err)
		return res, err
	}
	return res, nil
}

// UpdateMany to update many document
// See https://godoc.org/go.mongodb.org/mongo-driver/mongo#Collection.UpdateMany
// Example:
//
// 		type test struct {
// 			TestId string `bson:"test_id"`
// 			Name string `bson:"name"`
// 		}
// 		filter := bson.D{{"name", "name002"}}
// 		update := bson.D{{"$set", bson.D{{"name", "newname002"}}}}
// 		res, err := mongo.UpdateMany(collectionName, filter, update)
// 		fmt.Printf("res: %+v\n", res)
//
func (mc *MongoClient) UpdateMany(collectionName string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	mc.logger = mc.loggerClone
	mc.logger.SugaredLogger = mc.logger.With("method", "UpdateMany")

	collection := mc.GetCollectionHandler(collectionName)
	res, err := collection.UpdateMany(mc.getContext(), filter, update, opts...)
	if err != nil {
		mc.logger.Errorw("update many data error", "error", err)
		return res, err
	}
	return res, nil
}

// DeleteOne to delete one document
// See https://godoc.org/go.mongodb.org/mongo-driver/mongo#Collection.DeleteOne
// Example:
//
// 		type test struct {
// 			TestId string `bson:"test_id"`
// 			Name string `bson:"name"`
// 		}
// 		filter := bson.D{{"test_id", "id0001"}}
// 		res, err := mongo.DeleteOne(collectionName, filter)
// 		fmt.Printf("res: %+v\n", res)
//
func (mc *MongoClient) DeleteOne(collectionName string, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	mc.logger = mc.loggerClone
	mc.logger.SugaredLogger = mc.logger.With("method", "DeleteOne")

	collection := mc.GetCollectionHandler(collectionName)
	res, err := collection.DeleteOne(mc.getContext(), filter, opts...)
	if err != nil {
		mc.logger.Errorw("delete one data error", "error", err)
		return res, err
	}
	return res, nil
}

// DeleteMany to delete many documents
// See https://godoc.org/go.mongodb.org/mongo-driver/mongo#Collection.DeleteMany
// Example:
//
// 		type test struct {
// 			TestId string `bson:"test_id"`
// 			Name string `bson:"name"`
// 		}
// 		filter := bson.D{{"name", "name003"}}
// 		res, err := mongo.DeleteMany(collectionName, filter)
// 		fmt.Printf("res: %+v\n", res)
//
func (mc *MongoClient) DeleteMany(collectionName string, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	mc.logger = mc.loggerClone
	mc.logger.SugaredLogger = mc.logger.With("method", "DeleteMany")

	collection := mc.GetCollectionHandler(collectionName)
	res, err := collection.DeleteMany(mc.getContext(), filter, opts...)
	if err != nil {
		mc.logger.Errorw("delete many data error", "error", err)
		return res, err
	}
	return res, nil
}

// Distinct to get unique values of documents
// See https://godoc.org/go.mongodb.org/mongo-driver/mongo#Collection.Distinct
// Example:
//
// 		type test struct {
// 			TestId string `bson:"test_id"`
// 			Name string `bson:"name"`
// 		}
// 		res, err := mongo.Distinct(collectionName, "test_id", bson.D{{"name", "name001"}})
// 		fmt.Printf("res: %+v\n", res)
//
func (mc *MongoClient) Distinct(collectionName string, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error) {
	mc.logger = mc.loggerClone
	mc.logger.SugaredLogger = mc.logger.With("method", "Distinct")

	collection := mc.GetCollectionHandler(collectionName)
	res, err := collection.Distinct(mc.getContext(), fieldName, filter, opts...)
	if err != nil {
		mc.logger.Errorw("get distinct data error", "error", err)
		return res, err
	}
	return res, nil
}

// CountDocumentsByFilter to get the amount of documents in some condition
// See https://godoc.org/go.mongodb.org/mongo-driver/mongo#Collection.CountDocuments
// Example:
//
// 		type test struct {
// 			TestId string `bson:"test_id"`
// 			Name string `bson:"name"`
// 		}
// 		res, err := mongo.CountDocumentsByFilter(collectionName, bson.D{{"name", "name001"}})
// 		fmt.Printf("res: %+v\n", res)
//
func (mc *MongoClient) CountDocumentsByFilter(collectionName string, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	mc.logger = mc.loggerClone
	mc.logger.SugaredLogger = mc.logger.With("method", "CountDocuments")

	collection := mc.GetCollectionHandler(collectionName)
	total, err := collection.CountDocuments(mc.getContext(), filter, opts...)
	mc.logger.Infow("", "opts", opts)
	if err != nil {
		mc.logger.Errorw("count documents by filter error", "error", err)
		return total, err
	}
	return total, nil
}

// CountDocumentsTotal to get the total amount of documents
// See https://godoc.org/go.mongodb.org/mongo-driver/mongo#Collection.EstimatedDocumentCount
// Example:
//
// 		type test struct {
// 			TestId string `bson:"test_id"`
// 			Name string `bson:"name"`
// 		}
// 		res, err := mongo.CountDocumentsTotal(collectionName)
// 		fmt.Printf("res: %+v\n", res)
//
func (mc *MongoClient) CountDocumentsTotal(collectionName string, opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	mc.logger = mc.loggerClone
	mc.logger.SugaredLogger = mc.logger.With("method", "CountDocuments")

	collection := mc.GetCollectionHandler(collectionName)
	total, err := collection.EstimatedDocumentCount(mc.getContext(), opts...)
	if err != nil {
		mc.logger.Errorw("count documents total error", "error", err)
		return total, err
	}
	return total, nil
}

// CreateOneIndex to create one index
// See https://godoc.org/go.mongodb.org/mongo-driver/mongo#IndexView.CreateOne
// Example:
//
// 		model := mongo.IndexModel{
// 			Keys:    bson.D{{"name", 1}, {"age", 1}},
// 			Options: options.Index().SetName("nameAge"),
// 		}
// 		res, err := mongo.CreateOneIndex(collectionName, model)
// 		fmt.Printf("res: %+v\n", res)
//
func (mc *MongoClient) CreateOneIndex(collectionName string, model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error) {
	mc.logger = mc.loggerClone
	mc.logger.SugaredLogger = mc.logger.With("method", "CreateOneIndex")

	collection := mc.GetCollectionHandler(collectionName)
	res, err := collection.Indexes().CreateOne(mc.getContext(), model, opts...)
	if err != nil {
		mc.logger.Errorw(
			"create one index error",
			"error", err,
			"collectionName", collectionName,
			"model", model,
			"opts", opts,
		)
		return res, err
	}
	return res, nil
}

// CreateManyIndexes to create many indexes
// See https://godoc.org/go.mongodb.org/mongo-driver/mongo#IndexView.CreateMany
// Example:
//
// 		models := []mongo.IndexModel{
// 			{
// 				Keys: bson.D{{"name", 1}, {"email", 1}},
// 			},
// 			{
// 				Keys:    bson.D{{"name", 1}, {"age", 1}},
// 				Options: options.Index().SetName("nameAge"),
// 			},
// 		}
// 		res, err := mongo.CreateOneIndex(collectionName, model)
// 		fmt.Printf("res: %+v\n", res)
//
func (mc *MongoClient) CreateManyIndexes(collectionName string, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	mc.logger = mc.loggerClone
	mc.logger.SugaredLogger = mc.logger.With("method", "CreateManyIndexes")

	collection := mc.GetCollectionHandler(collectionName)
	res, err := collection.Indexes().CreateMany(mc.getContext(), models, opts...)
	if err != nil {
		mc.logger.Errorw(
			"create many indexes error",
			"error", err,
			"collectionName", collectionName,
			"models", models,
			"opts", opts,
		)
		return res, err
	}
	return res, nil
}

// DecodeDocument decodes a mongo bson.D doc to struct
func DecodeDocument(doc interface{}, val interface{}) error {
	bsonBytes, err := bson.Marshal(doc)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(bsonBytes, val)
	return err
}
