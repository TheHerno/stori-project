package email

import (
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/utils/constant"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetTransactionByMonth(t *testing.T) {
	t.Run("Should success on", func(t *testing.T) {
		t.Run("Getting list by month", func(t *testing.T) {
			// fixture
			movements := []entity.Movement{
				{
					Date: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					Date: time.Date(2020, time.March, 2, 0, 0, 0, 0, time.UTC),
				},
			}
			// action
			list := getTransactionByMonth(movements)

			// assert
			assert.Equal(t, "Number of transactions in January: 1<br>Number of transactions in March: 1<br>", list)
		})
	})
}

func TestGetAvgCredit(t *testing.T) {
	t.Run("Should success on", func(t *testing.T) {
		t.Run("Getting average credit", func(t *testing.T) {
			// fixture
			movements := []entity.Movement{
				{
					Date:     time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
					Type:     constant.IncomeType,
					Quantity: 100.00,
				},
				{
					Date:     time.Date(2020, time.March, 2, 0, 0, 0, 0, time.UTC),
					Type:     constant.IncomeType,
					Quantity: 200.00,
				},
			}
			// action
			avgCredit := getAvgCredit(movements)

			// assert
			assert.Equal(t, float64(150), avgCredit)
		})
	})
}

func TestGetAvgDebit(t *testing.T) {
	t.Run("Should success on", func(t *testing.T) {
		t.Run("Getting average debit", func(t *testing.T) {
			// fixture
			movements := []entity.Movement{
				{
					Date:     time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
					Type:     constant.OutcomeType,
					Quantity: 100.00,
				},
				{
					Date:     time.Date(2020, time.March, 2, 0, 0, 0, 0, time.UTC),
					Type:     constant.OutcomeType,
					Quantity: 200.00,
				},
			}
			// action
			avgDebit := getAvgDebit(movements)

			// assert
			assert.Equal(t, float64(150), avgDebit)
		})
	})
}
