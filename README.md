# goprimepayer
[![GoDoc](https://godoc.org/github.com/kikree/goprimepayer?status.png)](https://godoc.org/github.com/kikree/goprimepayer)
<br />
primepayer.com api wrapper written in golang


## create payment:
```go
package main

import (
	"fmt"
	"math/big"
	
	"github.com/kikree/goprimepayer"
	"github.com/shopspring/decimal"
)

var pp = goprimepayer.New(&goprimepayer.Config{
        ShopID: 1, // Your shop id
        Secret: "secret", // Your shop secret
})

func main() {
    payment := pp.NewPayment(
        big.NewInt(1), // Payment id
        3, // Currency (3 - RUB)
        decimal.New(100, 0), // Amount
        fmt.Sprintf("Something for %s", "User"), // Description
    )
    // Set user variable
    payment.Set("user_id", 1)
    
    // Make sign
    sign := payment.Sign()
    
    fmt.Println(sign)
}
```

## verify notification:
```go
package main

import (
	"fmt"
	
	"github.com/kikree/goprimepayer"
)

var pp = goprimepayer.New(&goprimepayer.Config{
        ShopID: 1, // Your shop id
        Secret: "secret", // Your shop secret
})

func main() {
	postParams := map[string]string{}
	
   	if notification, err := pp.VerifyNotification(postParams); err != nil {
   		// error
   	} else {
   		// ok
   		fmt.Println(notification.Get("user_id"))
   	}
}
```