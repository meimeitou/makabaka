package test

import (
	"fmt"
	"testing"

	"github.com/meimeitou/makabaka/db"
	"github.com/meimeitou/makabaka/db/postgres"
	"github.com/meimeitou/makabaka/model"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type ModelTestSuite struct {
	conn *db.Conn
	suite.Suite
}

func (suite *ModelTestSuite) SetupTest() {
	st := postgres.PostgresStorage{
		NetworkDB: db.NetworkDB{
			Host:     "localhost",
			Port:     5432,
			User:     "makabaka",
			Password: "makabaka",
			Database: "makabaka",
			Debug:    true,
		},
	}
	conn, err := st.Open(logrus.StandardLogger())
	if err != nil {
		panic(err)
	}
	suite.conn = conn
	suite.conn.SourceSQL("./teardown.sql")
	err = conn.AutoMigrate(&model.Apis{})
	suite.Nil(err)
}

func (suite *ModelTestSuite) TearDownSuite() {
	// suite.conn.SourceSQL("./teardown.sql")
}

func (suite *ModelTestSuite) TestApi() {
	api1 := &model.Apis{
		Name:        "test",
		Method:      "GET",
		SqlType:     model.SqlTypeTemplate,
		SqlTemplate: "",
		SqlTemplateParameters: model.JSONMap{
			"a": "aa",
			"b": 1,
		},
		// SqlTemplateResult: JSONMap{},
	}
	a := suite.conn.GetQuery().Apis
	err := a.Create(api1)
	suite.Nil(err)
	fmt.Println(api1)

	api2 := &model.Apis{
		Modelx:                model.Modelx{ID: api1.ID},
		Method:                "POST",
		Name:                  api1.Name,
		SqlType:               api1.SqlType,
		SqlTemplate:           api1.SqlTemplate,
		SqlTemplateParameters: api1.SqlTemplateParameters,
	}
	err = a.Save(api2)
	suite.Nil(err)
	tmp, err := a.First()
	suite.Nil(err)
	fmt.Println(tmp)
}

func TestDBCreateSuite(t *testing.T) {
	suite.Run(t, new(ModelTestSuite))
}
