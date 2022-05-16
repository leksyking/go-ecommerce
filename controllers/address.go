package contollers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leksyking/go-ecommerce/database"
	"github.com/leksyking/go-ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid code"})
			c.Abort()
			return
		}
		address, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal server error")
		}
		var addresses models.Address
		addresses.Address_ID = primitive.NewObjectID()
		if err := c.BindJSON(&addresses); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		grouping := bson.D{{Key: "$group", Value: {primitive.E{Key: "_id", Value: "$_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}
		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, grouping})
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}
		var listing []bson.M
		if err := pointcursor.All(ctx, &listing); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		for _, json := range listing {
			c.IndentedJSON(200, json["count"])
		}
		defer cancel()
	}
}
func EditHomeAddress() gin.HandlerFunc {

}

func EditWorkAddress() gin.HandlerFunc {

}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid search index"})
			c.Abort()
			return
		}
		addresses := make([]models.Address, 0)
		userId, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal Server error")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: userId}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
		err = UserCollection.FindOneAndUpdate(ctx, filter, update).Err()
		if err != nil {
			c.IndentedJSON(400, "Unable to update")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Address deleted successfully")
	}
}
