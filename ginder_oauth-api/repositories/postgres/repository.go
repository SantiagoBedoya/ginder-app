package postgres

import (
	"fmt"

	"github.com/SantiagoBedoya/ginder_oauth-api/oauth"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
)

type postgresRepository struct {
	DB *gorm.DB
}

func newPostgresClient(postgresURL string, tables ...interface{}) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", postgresURL)
	if err != nil {
		return nil, errors.Wrap(err, "newPostgresClient")
	}
	db.AutoMigrate(tables...)
	return db, err
}

func NewPostgresRepository(postgresURL string, tables ...interface{}) (oauth.AccessTokenRepository, error) {
	repo := &postgresRepository{}
	db, err := newPostgresClient(postgresURL, tables...)
	if err != nil {
		return nil, errors.Wrap(err, "NewPostgresRepository")
	}
	repo.DB = db
	return repo, nil
}

func (r *postgresRepository) FindByToken(token string) (*oauth.AccessToken, error) {
	at := &oauth.AccessToken{}
	r.DB.Where("token = ?", token).First(at)
	fmt.Println(at)
	if at.ID == 0 {
		return nil, oauth.ErrAccessTokenNotFound
	}
	return at, nil
}

func (r *postgresRepository) FindByUserID(userID string) (*oauth.AccessToken, error) {
	at := &oauth.AccessToken{}
	err := r.DB.Where("user_id = ?", userID).First(at).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, oauth.ErrAccessTokenNotFound
		}
	}
	return at, nil
}

func (r *postgresRepository) Create(at *oauth.AccessToken) (*oauth.AccessToken, error) {
	r.DB.Create(at)
	return at, nil
}
