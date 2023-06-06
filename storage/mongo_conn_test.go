package storage

import (
	"testing"
)

func TestConnectToMongoDB(t *testing.T) {

	given, when, then := NewMongoConnTestStage(t)

	given.user_and_password_enviroment_variables_are_set()

	when.connect_to_mongodb()

	then.no_errors_happen().
		and().
		client_is_not_nil()

}
