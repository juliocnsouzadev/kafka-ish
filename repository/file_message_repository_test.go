package repository

import (
	"testing"
)

func TestInsert_File_Success(t *testing.T) {

	given, when, then := NewFileMessageTestStage(t)

	given.a_message()

	when.the_message_is_created()

	then.no_errors_happen()
}

func TestFindByTopic_File_Success(t *testing.T) {
	given, when, then := NewFileMessageTestStage(t)

	given.some_existing_messages_in_the_topic("tests_topic")

	when.find_messages_for_topic("tests_topic")

	then.no_errors_happen().and().
		found_messages_matches_existing_messages()
}
