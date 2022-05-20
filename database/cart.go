package database

import (
	"context"
	"errors"
	"log"
	"time"

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
	return nil
}

func RemoveItemFromCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userId string) error {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	ctx.Done()
	if err != nil {
		return ErrCantRemoveItemCart
	}
	return nil
}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userId string) error {
	//fetch items from cart
	//find the cart total
	//create an order with the items
	//add order to user collectin
	//add items in the cart to order list
	// empty up cart
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	var getcartitems models.User
	var ordercart models.Order

	ordercart.Order_ID = primitive.NewObjectID()
	ordercart.Ordered_At = time.Now()
	ordercart.Order_Cart = make([]models.ProductUser, 0)
	ordercart.Payment_Method.COD = true
	//calculate the aggregate of price
	match := bson.D{primitive.E{Key: "_id", Value: id}}
	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: bson.D{{Key: "usercart"}}}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

	currentResults, err := userCollection.Aggregate(ctx, mongo.Pipeline{match, unwind, grouping})
	ctx.Done()
	if err != nil {
		panic(err)
	}
	var getUserCart []bson.M

	if err := currentResults.All(ctx, &getUserCart); err != nil {
		panic(err)
	}
	var total_price int32

	for _, user_item := range getUserCart {
		price := user_item["total"]
		total_price = price.(int32)
	}

	ordercart.Price = int(total_price)
	//add orders
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: ordercart}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}
	//find user
	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getcartitems)
	if err != nil {
		log.Println(err)
	}
	//add usercart to orders
	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each": getcartitems.UserCart}}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}
	//empty the cart
	usercart_empty := make([]models.ProductUser, 0)
	filter3 := bson.D{primitive.E{Key: "_id", Value: id}}
	update3 := bson.M{"$set": bson.M{"usercart": usercart_empty}}
	_, err = userCollection.UpdateOne(ctx, filter3, update3)
	if err != nil {
		return ErrCantBuyCartItem
	}
	return nil
}

func InstantBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userId string) error {

}
