package main

import (
	"fmt"
	"time"

	"github.com/kataras/iris"
	"gopkg.in/mgo.v2/bson"

	"github.com/kataras/iris/context"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/proci/network"
)

type User struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	Firstname  string        `json:"firstname"`
	Lastname   string        `json:"lastname"`
	Age        int           `json:"age"`
	Msisdn     string        `json:"msisdn"`
	InsertedAt time.Time     `json:"inserted_at" bson:"inserted_at"`
	LastUpdate time.Time     `json:"last_update" bson:"last_update"`
}

func newApp() *iris.Application {
	app := iris.New()
	app.Logger().SetLevel("debug")
	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())

	// Method:   GET Default Endpoint
	// Resource: http://localhost:8080
	app.Handle("GET", "/", func(ctx context.Context) {

		network.GenerateConfigTx("./crypto-config", "nvxtien.com")
		network.GenerateCryptoCfg()
		network.ExecuteCryptogen()
		network.CreateOrderGenesisBlock()

		ctx.JSON(context.Map{"message": "Welcome User Micro Service"})
	})

	// Gets all users
	// Method:   GET
	// Resource: this to get all all users
	app.Handle("GET", "/users", func(ctx context.Context) {
		results := []User{}
		ctx.JSON(context.Map{"response": results})
	})

	// Gets a single user
	// Method:   GET
	// Resource: this to get all all users
	app.Handle("GET", "/users/{msisdn: string}", func(ctx context.Context) {
		msisdn := ctx.Params().Get("msisdn")
		fmt.Println(msisdn)
		if msisdn == "" {
			ctx.JSON(context.Map{"response": "please pass a valid msisdn"})
		}
		result := User{}
		ctx.JSON(context.Map{"response": result})
	})

	// Method:   POST
	// Resource: This is to create a new user
	app.Handle("POST", "/users", func(ctx context.Context) {
		params := &User{}
		err := ctx.ReadJSON(params)
		if err != nil {
			ctx.JSON(context.Map{"response": err.Error()})
		} else {
			params.LastUpdate = time.Now()

			//err := c.Insert(params)
			//if err != nil {
			//	ctx.JSON(context.Map{"response": err.Error()})
			//} else {
			//	fmt.Println("Successfully inserted into database")
			//	result := User{}
			//	err = c.Find(bson.M{"msisdn": params.Msisdn}).One(&result)
			//	if err != nil {
			//		ctx.JSON(context.Map{"response": err.Error()})
			//	}
			//	ctx.JSON(context.Map{"response": "User succesfully created", "message": result})
			//}
			result := User{}
			ctx.JSON(context.Map{"response": "User succesfully created", "message": result})
		}

	})

	// Method:   PATCH
	// Resource: This is to update a user record
	app.Handle("PATCH", "/users/{msisdn: string}", func(ctx context.Context) {
		msisdn := ctx.Params().Get("msisdn")
		fmt.Println(msisdn)
		if msisdn == "" {
			ctx.JSON(context.Map{"response": "please pass a valid msisdn"})
		}
		params := &User{}
		err := ctx.ReadJSON(params)
		if err != nil {
			ctx.JSON(context.Map{"response": err.Error()})
		} else {
			params.InsertedAt = time.Now()
			//query := bson.M{"msisdn": msisdn}
			//err = c.Update(query, params)
			if err != nil {
				ctx.JSON(context.Map{"response": err.Error()})
			} else {
				result := User{}
				//err = c.Find(bson.M{"msisdn": params.Msisdn}).One(&result)
				//if err != nil {
				//	ctx.JSON(context.Map{"response": err.Error()})
				//}
				ctx.JSON(context.Map{"response": "user record successfully updated", "data": result})
			}
		}

	})

	// Method:   DELETE
	// Resource: This is to delete a user record
	app.Handle("DELETE", "/users/{msisdn: string}", func(ctx context.Context) {
		msisdn := ctx.Params().Get("msisdn")
		fmt.Println(msisdn)
		if msisdn == "" {
			ctx.JSON(context.Map{"response": "please pass a valid msisdn"})
		}
		params := &User{}
		err := ctx.ReadJSON(params)
		if err != nil {
			ctx.JSON(context.Map{"response": err.Error()})
		} else {
			params.InsertedAt = time.Now()
			//query := bson.M{"msisdn": msisdn}
			//err = c.Remove(query)
			if err != nil {
				ctx.JSON(context.Map{"response": err.Error()})
			} else {
				ctx.JSON(context.Map{"response": "user record successfully deleted"})
			}
		}

	})

	return app
}

func main() {
	app := newApp()
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}