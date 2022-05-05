package database

import "errors"

var (
	ErrCantFindProduct    = errors.New("Can't find the product")
	ErrCantDecodepPoducts = errors.New("Can't find the product")
	ErrUserIdIsNotValid   = errors.New("This user is not valid")
	ErrCantUpdateUser     = errors.New("Can't add this product to the cart")
	ErrCantRemoveItemCart = errors.New("Can't remove this item from the cart")
	ErrCantGetItem        = errors.New("Was unable to get the item from cart")
	ErrCantBuyCartItem    = errors.New("Can't update the purchase")
)

func AddProductToCart() {

}

func RemoveItemFromCart() {

}

func BuyItemFromCart() {

}

func InstantBuyer() {

}
