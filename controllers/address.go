package contollers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leksyking/go-ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAddress() gin.HandlerFunc {

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
