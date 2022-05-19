package database

import (
	"context"
	"errors"
	"log"

	"github.com/leksyking/go-ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct    = errors.New("Can't find the product")
	ErrCantDecodepPoducts = errors.New("Can't find the product")
	ErrUserIdIsNotValid   = errors.New("This user is not valid")
	ErrCantUpdateUser     = errors.New("Can't add this product to the cart")
	ErrCantRemoveItemCart = errors.New("Can't remove this item from the cart")
	ErrCantGetItem        = errors.New("Was unable to get the item from cart")
	ErrCantBuyCartItem    = errors.New("Can't update the purchase")
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	searchfromdb, err := prodCollection.Find(ctx, bson.M{"_id": productID})
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}
	var productCart []models.ProductUser
	if err := searchfromdb.All(ctx, &productCart); err != nil {
		log.Println(err)
		return ErrCantDecodepPoducts
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	ctx.Done()
	if err != nil {
		return ErrCantUpdateUser
	}
}

func RemoveItemFromCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userId string) error {
	//check whether product is valid
	var product models.ProductUser
	err := prodCollection.FindOne(ctx, bson.M{"_id": productID}).Decode(&product)
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$remove", Value: bson.D{primitive.E{Key: "usercart", Value: product}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	ctx.Done()
	if err != nil {
		return ErrCantUpdateUser
	}

}

func BuyItemFromCart() {

}

func InstantBuyer() {

}
