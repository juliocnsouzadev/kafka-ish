package server

import (
	"testing"
)

func TestHttpServerStart(t *testing.T) {

	given, when, then := NewHttpServerTestStage(t)

	given.two_publish_commands().and().
		a_http_server()

	when.http_server_starts()

	then.no_erros_ocurr().and().
		interactions_are_correct()
}
