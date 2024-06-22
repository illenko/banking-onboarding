package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	BadRequestCode    = "bad_request"
	BadRequestMessage = "invalid request"
)

func antifraudChecksHandler(c *gin.Context) {
	var request UserRequest

	err := c.Bind(&request)
	if err != nil {
		handleError(c, http.StatusBadRequest, BadRequestCode, BadRequestMessage)
		return
	}

	if strings.Contains(strings.ToLower(request.FirstName), "fraud") ||
		strings.Contains(strings.ToLower(request.LastName), "fraud") {
		c.JSON(http.StatusOK, AntifraudResponse{
			Passed:  false,
			Comment: "User is fraud",
		})
		return
	}

	c.JSON(http.StatusOK, AntifraudResponse{
		Passed:  true,
		Comment: "ok",
	})
}

func usersHandler(c *gin.Context) {
	var request UserRequest

	err := c.Bind(&request)
	if err != nil {
		handleError(c, http.StatusBadRequest, BadRequestCode, BadRequestMessage)
		return
	}

	user := UserResponse{
		ID:        uuid.New(),
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		City:      request.City,
	}

	c.JSON(http.StatusOK, user)
}

func agreementsHandler(c *gin.Context) {
	var request AgreementRequest

	err := c.Bind(&request)
	if err != nil {
		handleError(c, http.StatusBadRequest, BadRequestCode, BadRequestMessage)
		return
	}

	id := uuid.New()

	agreement := AgreementResponse{
		ID:   id,
		Link: "https://perfect-bank.ua/agreements/" + id.String(),
	}

	c.JSON(http.StatusOK, agreement)
}

func signaturesHandler(c *gin.Context) {
	var request SignatureRequest

	err := c.Bind(&request)

	if err != nil {
		handleError(c, http.StatusBadRequest, BadRequestCode, BadRequestMessage)
		return
	}

	if strings.Contains(request.Signature, "fraud") {
		c.JSON(http.StatusOK, SignatureResponse{
			ID:      uuid.New(),
			Valid:   false,
			Comment: "Signature is fraud",
		})
		return
	}

	c.JSON(http.StatusOK, SignatureResponse{
		ID:      uuid.New(),
		Valid:   true,
		Comment: "ok",
	})
}

func accountsHandler(c *gin.Context) {
	var request AccountRequest

	err := c.Bind(&request)
	if err != nil {
		handleError(c, http.StatusBadRequest, BadRequestCode, BadRequestMessage)
		return
	}

	account := AccountResponse{
		ID:       uuid.New(),
		UserID:   request.UserID,
		Currency: request.Currency,
		Type:     request.Type,
		Iban:     fmt.Sprintf("UA8937040044%016d", rand.Int63n(9999999999999999-1000000000000000)+1000000000000000),
		Balance:  0,
	}

	c.JSON(http.StatusOK, account)
}

func cardsHandler(c *gin.Context) {
	var request CardRequest

	err := c.Bind(&request)
	if err != nil {
		handleError(c, http.StatusBadRequest, BadRequestCode, BadRequestMessage)
		return
	}

	expiryDate := time.Now().AddDate(2, 0, 0)
	expiry := fmt.Sprintf("%02d/%02d", expiryDate.Month(), expiryDate.Year()%100)

	card := CardResponse{
		ID:        uuid.New(),
		AccountID: request.AccountID,
		Number:    fmt.Sprintf("%016d", rand.Int63n(9999999999999999-1000000000000000)+1000000000000000),
		Expire:    expiry,
		Cvv:       fmt.Sprintf("%03d", rand.Intn(1000)),
	}

	c.JSON(http.StatusOK, card)
}

func handleError(c *gin.Context, statusCode int, errorCode, errorMessage string) {
	c.JSON(statusCode, ErrorResponse{
		Code:    errorCode,
		Message: errorMessage,
	})
	c.Abort()
}
