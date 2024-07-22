// Package storage allows to keep payment results in the database, and provides the handling of the /recall endpoint.
package storage

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cko-recruitment/payment-gateway-challenge-go/common"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Storage is the primary interface of this module.
type Storage interface {
	RecordPayment(payment *common.PaymentResponse) error
	Recall(c *gin.Context)
}

// A concrete implementation of the interface Storage.
type storage struct {
	// While this cache can speed up the lookup of current session's payment results up to some extent,
	// its primary purpose here is to provide some transient storage in absence of a real database.
	// TODO: in a real prod we would need to forget old payments to not run out of memory.
	cache map[common.ID]common.PaymentResponse
	db *gorm.DB
	logger *log.Logger
}

// New returns a fresh Storage object.
// NB: despite error being returned, currently no errors are provided.
//	db	a (possibly nil) pointer to a gorm object
//	returns (a fresh storage object, any error occurred)
func New(db *gorm.DB) (Storage, error) {
	logger := common.MakeLogger("STORAGE")
	if db == nil {
		logger.Print("no external DB connection, temporary session-limited payments cache only")
	} else {
		logger.Print("running in persistent storage mode")
	}
	return &storage{
		cache: make(map[common.ID]common.PaymentResponse),
		db: db,
		logger: logger,
	}, nil
}

// RecordPayment is called after receiving a successful payment (un)authorization from the bank.
// It stores the payment result in the session cache, and in the permanent storage, provided there is one.
//	payment	a payment response
//	returns	any error occurred
func (s *storage) RecordPayment(payment *common.PaymentResponse) error {
	s.logger.Printf("saving payment %v", payment.ID)
	s.cache[payment.ID] = *payment
	if s.db != nil {
		return s.db.Save(payment).Error
	}
	return nil
}

// Reca;; godoc
//	@Summary	The /recall endpoint
//	@Schemes
//	@Description	Request the result of a past payment.

//	@Param		payment_id	query	string	true	"id of the (un)authorized payment"
//	@Produce	json
//	@Success	200	{object}	common.PaymentResponse
//	@Failure	400	{object}	error
//	@Router		/recall [get]
func (s *storage) Recall(c *gin.Context) {
	queryParam, ok := c.GetQuery("payment_id")
	if !ok {
		s.logger.Printf("invalid recall request: missing payment_id")
		common.Refuse(c, fmt.Errorf("invalid recall request: missing payment_id"))
		return
	}

	id, err := common.ParseID(queryParam)
	if err != nil {
		common.Refuse(c, err)
		return
	}

	s.logger.Printf("recall request for payment %v", id)
	if resp, ok := s.cache[id]; ok {
		s.logger.Printf("Found cache entry, card no %s", resp.CardNumber)
		c.JSON(http.StatusOK, &resp)
		return
	}
	if s.db == nil {
		s.logger.Printf("no persinstent storage, cannot recognize payment")
		common.Refuse(c, fmt.Errorf("Payment %v not found", id))
		return
	}
	resp := common.PaymentResponse{ID: id}
	tx := s.db.First(&resp)
	if tx.RowsAffected == 0 {
		s.logger.Printf("no rows found that match ID: %s", tx.Error)
		common.Refuse(c, tx.Error)
		return
	}
	s.logger.Printf("found record in persistent DB")
	c.JSON(http.StatusOK, &resp)
}
