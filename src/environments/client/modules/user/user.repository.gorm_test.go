package user

import (
	"os"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/database"
	"stori-service/src/libs/errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	// setup
	database.SetupStoriGormDB()
	code := m.Run()
	os.Exit(code)
}

type result struct {
	User *entity.User
	Err  error
	Date time.Time
}

var users = []entity.User{
	{
		UserID: 1,
		Name:   "User 1",
	},
	{
		UserID: 2,
		Name:   "User 2",
	},
	{
		UserID: 3,
		Name:   "User 3",
	},
	{
		UserID: 4,
		Name:   "User 4",
	},
}

/*
	Fixtures: four users
*/
func addFixtures(tx *gorm.DB) {
	tx.Unscoped().Where("1=1").Delete(&entity.User{}) // cleaning users
	tx.Create(users)
}

func TestGormRepository(t *testing.T) {
	t.Run("FindByUserID", func(t *testing.T) {
		//Fixture
		userID := users[0].UserID
		t.Run("Should success on", func(t *testing.T) {
			tx := database.GetStoriGormConnection().Begin()
			addFixtures(tx)
			rUser := NewUserGormRepo(tx)

			got, err := rUser.FindByUserID(userID)

			//Data Assertion
			assert.NoError(t, err)
			assert.True(t, cmp.Equal(got, &users[0], cmpopts.IgnoreFields(entity.User{}, "UpdatedAt", "CreatedAt")))

			t.Cleanup(func() {
				tx.Rollback()
			})
		})

		t.Run("Should fail on", func(t *testing.T) {
			t.Run("User not found", func(t *testing.T) {
				//Fixture
				userID := 9999999
				tx := database.GetStoriGormConnection().Begin()
				addFixtures(tx)

				rUser := NewUserGormRepo(tx)

				got, err := rUser.FindByUserID(userID)

				//Data Assertion
				assert.EqualError(t, err, errors.ErrNotFound.Error())
				assert.Nil(t, got)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
			t.Run("Table doesn't exist", func(t *testing.T) {
				connection := database.GetStoriGormConnection()
				tx := connection.Begin()
				rUser := NewUserGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.User{}) //Cleaning users
				tx.Migrator().DropTable(&entity.User{})

				got, err := rUser.FindByUserID(userID)

				//Data Assertion
				assert.Nil(t, got)
				assert.Error(t, err)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
		})
	})
	t.Run("FindAndLockByID", func(t *testing.T) {
		t.Run("Should success on", func(t *testing.T) {
			// fixture
			userIDToFind := users[1].UserID
			channelResult1 := make(chan result)
			channelResult2 := make(chan result)

			connection := database.GetStoriGormConnection()
			addFixtures(connection)
			tx := connection.Begin()
			rUser := NewUserGormRepo(tx)

			go func() {
				got, err := rUser.FindAndLockByUserID(userIDToFind)
				channelResult1 <- result{got, err, time.Now()}
				time.Sleep(500 * time.Millisecond)
				tx.Commit()
			}()

			go func() {
				err2 := connection.Table("user").Where("user_id = ?", userIDToFind).
					Update("name", "Useritooo").Error
				channelResult2 <- result{nil, err2, time.Now()}
			}()

			// assertions
			result1 := <-channelResult1
			assert.NoError(t, result1.Err)
			assert.True(t, cmp.Equal(result1.User, &users[1], cmpopts.IgnoreTypes(time.Time{})))
			// tries to update but still nil because the lock is not released
			result2 := <-channelResult2
			assert.Nil(t, result2.User)
			assert.NoError(t, result2.Err)
			assert.True(t, result2.Date.After(result1.Date))

			t.Cleanup(func() {
				tx.Rollback()
				connection.Unscoped().Where("1=1").Delete(&entity.User{})
			})
		})
		t.Run("Should fail on", func(t *testing.T) {
			t.Run("User doesn't exists", func(t *testing.T) {
				// fixture
				userIDToFind := 78
				connection := database.GetStoriGormConnection()
				addFixtures(connection)
				tx := connection.Begin()
				rUser := NewUserGormRepo(tx)

				// action
				got, err := rUser.FindAndLockByUserID(userIDToFind)

				// assertions
				assert.Nil(t, got)
				assert.Error(t, err)
				assert.EqualError(t, err, errors.ErrNotFound.Error())

				t.Cleanup(func() {
					tx.Rollback()
					connection.Unscoped().Where("1=1").Delete(&entity.User{})
				})
			})
			t.Run("Table doesn't exist", func(t *testing.T) {
				connection := database.GetStoriGormConnection()
				tx := connection.Begin()
				rUser := NewUserGormRepo(tx)
				tx.Unscoped().Where("1=1").Delete(&entity.User{}) //Cleaning users
				tx.Migrator().DropTable(&entity.User{})

				got, err := rUser.FindAndLockByUserID(1)

				//Data Assertion
				assert.Nil(t, got)
				assert.Error(t, err)

				t.Cleanup(func() {
					tx.Rollback()
				})
			})
		})
	})
	t.Run("Clone", func(t *testing.T) {
		db := database.GetStoriGormConnection()
		rUser := NewUserGormRepo(db)

		clone := rUser.Clone()

		assert.NotNil(t, clone)
		assert.Equal(t, rUser, clone)
	})
}
