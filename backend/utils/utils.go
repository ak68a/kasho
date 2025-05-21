package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var Currencies = map[string]string {
	"USD": "USD",
	"NGN": "NGN",
	"ZAR": "ZAR",
}

func IsSupportedCurrency(currency string) bool {
	if _, ok := Currencies[currency]; ok {
		return true
	}
	return false
}
	
func GetActiveUser(c *gin.Context) (int64, error) {
	value, exists := c.Get("user_id")	
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return 0, fmt.Errorf("Error occured")
	}

	userId, ok := value.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get user id"})
		return 0, fmt.Errorf("Error occured")
	}

	return userId, nil
}
