package storage

import (
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoConnTestStage struct {
	t      *testing.T
	client *mongo.Client
	err    error
}

func NewMongoConnTestStage(t *testing.T) (*MongoConnTestStage, *MongoConnTestStage, *MongoConnTestStage) {
	stage := &MongoConnTestStage{t: t}
	return stage, stage, stage
}

func (s *MongoConnTestStage) and() *MongoConnTestStage {
	return s
}

func (s *MongoConnTestStage) user_and_password_enviroment_variables_are_set() *MongoConnTestStage {
	if user := os.Getenv("MONGO_USER"); user == "" {
		s.err = errors.Errorf("MONGO_USER environment variable is not set")
		return s
	}
	if pswd := os.Getenv("MONGO_PASSWORD"); pswd == "" {
		s.err = errors.Errorf("MONGO_PASSWORD environment variable is not set")
		return s
	}
	return s
}

func (s *MongoConnTestStage) connect_to_mongodb() *MongoConnTestStage {
	if s.err != nil {
		s.client, s.err = ConnectToMongoDB()
	}
	return s
}

func (s *MongoConnTestStage) no_errors_happen() *MongoConnTestStage {
	require.NoError(s.t, s.err)
	return s
}

func (s *MongoConnTestStage) client_is_not_nil() *MongoConnTestStage {
	require.NotNil(s.t, s.client)
	return s
}
