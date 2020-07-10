package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mysqlConfig = ManagerConfig{
		"mysql1": &Config{
			Driver:       "mysql",
			DSN:          "root:Hyeasy88@tcp(127.0.0.1:3306)/stack_base_test?charset=utf8mb4&parseTime=True&loc=Local",
			DialTimeout:  5000,
			ReadTimeout:  5000,
			WriteTimeout: 5000,
			MaxOpenConns: 50,
			MaxIdleConns: 20,
			MaxLifeConns: 60,
			Debug:        true,
			Metrics:      true,
		},
		"mysql2": &Config{
			Driver:       "mysql",
			DSN:          "root:Hyeasy88@tcp(127.0.0.1:3306)/stack_base_test?charset=utf8mb4&parseTime=True&loc=Local",
			DialTimeout:  5000,
			ReadTimeout:  5000,
			WriteTimeout: 5000,
			MaxOpenConns: 50,
			MaxIdleConns: 20,
			MaxLifeConns: 60,
			Debug:        true,
			Metrics:      true,
		},
	}
)

func TestNewClient(t *testing.T) {
	var (
		it     = assert.New(t)
		client *Client
	)

	it.NotPanics(func() {
		client = NewClient(WithDBConfig(mysqlConfig))
	})

	errs := client.Ping()
	it.Len(errs, 0)

	db, err := client.GetDb("mysql2")
	if it.Nil(err) {
		var (
			count              int
			id                 int64
			hashId, ObjectType string
		)
		err := db.Table("resource").Count(&count).Select("id, hash_id, object_type").
			Where("id = ?", 3).Row().Scan(&id, &hashId, &ObjectType)
		it.NoError(err)
		it.Equal("1d7do64tYsy47W0101rWGoqICZ5", hashId)
	}
}

//func TestClient_GetDb(t *testing.T) {
//	var (
//		it     = assert.New(t)
//		client *Client
//	)
//
//	it.NotPanics(func() {
//		client = NewClient(WithDBConfig(mysqlConfig))
//	})
//}
