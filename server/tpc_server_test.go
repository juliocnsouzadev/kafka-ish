package server

import (
	"testing"
)

func TestServerStart(t *testing.T) {

	given, when, then := NewTcpServerTestStage(t)

	given.two_publish_commands().and().
		a_tcp_server()

	when.tcp_server_starts()

	then.no_erros_ocurr().and().
		interactions_are_correct()
}
