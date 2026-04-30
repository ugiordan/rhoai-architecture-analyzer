package testdata

import (
	"database/sql"
	"net/http"
	"os/exec"
)

// TaintDirect: user input flows directly to SQL sink.
func TaintDirect(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	db, _ := sql.Open("postgres", "")
	db.Query(query)
}

// TaintViaCall: user input flows through a helper to SQL sink.
func TaintViaCall(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("cmd")
	runQuery(input)
}

func runQuery(q string) {
	db, _ := sql.Open("postgres", "")
	db.Exec(q)
}

// TaintUnreachable: user input with length check.
func TaintUnreachable(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("x")
	if len(input) > 100 {
		return
	}
	_ = input
}

// NoTaint: no user input, no taint path.
func NoTaint(x int, y int) int {
	result := x + y
	return result
}

// TaintToCommand: user input flows to command execution.
func TaintToCommand(w http.ResponseWriter, r *http.Request) {
	cmd := r.URL.Query().Get("cmd")
	exec.Command(cmd).Run()
}
