package bootstrap

import (
	"context"
	"fmt"
	"gae-backend-web/constant/dao"
	"gae-backend-web/infrastructure/hive"
	"gae-backend-web/infrastructure/mongo"
	"gae-backend-web/infrastructure/mysql"
	"gae-backend-web/infrastructure/redis"
	"log"
	"time"
)

func NewMongoDatabase(env *Env) mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbHost := env.MongoHost
	dbPort := env.MongoPort
	dbUser := env.MongoUser
	dbPass := env.MongoPass

	mongodbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", dbUser, dbPass, dbHost, dbPort)

	if dbUser == "" || dbPass == "" {
		mongodbURI = fmt.Sprintf("mongodb://%s:%s", dbHost, dbPort)
	}

	client, err := mongo.NewClient(mongodbURI)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func NewRedisDatabase(env *Env) redis.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbAddr := env.RedisAddr
	dbPassword := env.RedisPassword

	client, err := redis.NewRedisClient(redis.NewRedisApplication(dbAddr, dbPassword))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func NewMysqlDatabase(env *Env) mysql.Client {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		env.MysqlUser, env.MysqlPassword, env.MysqlHost, env.MysqlPort, env.MysqlDB)

	client, err := mysql.NewMysqlClient(dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func CloseMongoDBConnection(client mongo.Client) {
	if client == nil {
		return
	}

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to MongoDB closed.")
}

func NewHiveDBConnection(env *Env) hive.Client {
	client, err := hive.NewHiveClient(env.HiveUrl, env.HIvePort, dao.HiveNoneAuth)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connection to Hive Started")
	return client
}

func NewDatabases(env *Env) *Databases {
	return &Databases{
		Redis: NewRedisDatabase(env),
		Mysql: NewMysqlDatabase(env),
		Hive:  NewHiveDBConnection(env),
	}
}
