package repository

import (
	"testing"
)

func TestInsert_MongoDB_Success(t *testing.T) {

	given, when, then := NewMongoMessageTestStage(t)

	given.a_message(ExpectedMessageOK)

	when.the_message_is_created()

	then.no_errors_happen()
}

func TestInsert_MongoDB_Fail_DulicatedKey(t *testing.T) {
	given, when, then := NewMongoMessageTestStage(t)

	given.a_message(ExpectedMessageDuplicated)

	when.the_message_is_created()

	then.fail_duplicated_key()

}

func TestFindByTopic_MongoDB_Success(t *testing.T) {
	given, when, then := NewMongoMessageTestStage(t)

	given.some_existing_messages_in_the_topic("tests_topic")

	when.find_messages_for_topic("tests_topic")

	then.no_errors_happen().and().
		found_messages_matches_existing_messages()
}
