package main

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestGetContacts(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	require.NoError(t, db.Ping())
	defer db.Close()

	_, err = db.Exec("create table contacts (id integer primary key, name varchar(100), email varchar(100))")
	require.NoError(t, err)

	_, err = db.Exec("insert into contacts (name, email) values ($1, $2)", "foo", "foo@gmail.com")
	require.NoError(t, err)
	_, err = db.Exec("insert into contacts (name, email) values ($1, $2)", "bar", "bar@gmail.com")
	require.NoError(t, err)

	server := NewServer(db)

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)
	resp := httptest.NewRecorder()

	server.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	require.JSONEq(t, `[
		{
			"id": 1,
			"name": "foo",
			"email": "foo@gmail.com"
		},
		{
			"id": 2,
			"name": "bar",
			"email": "bar@gmail.com"
		}
	]`, resp.Body.String())
}

func TestPostContacts(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	require.NoError(t, db.Ping())
	defer db.Close()

	_, err = db.Exec("create table contacts (id integer primary key, name varchar(100), email varchar(100))")
	require.NoError(t, err)

	server := NewServer(db)

	var body bytes.Buffer
	body.WriteString(`{
		"name": "foo",
		"email": "foo@gmail.com"
	}`)
	req, err := http.NewRequest("POST", "/", &body)
	require.NoError(t, err)
	resp := httptest.NewRecorder()

	server.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)
	require.JSONEq(t, `{
			"id": 1,
			"name": "foo",
			"email": "foo@gmail.com"
		}`, resp.Body.String())

	rows, err := db.Query("select id, name, email from contacts")
	require.NoError(t, err)
	contacts := make([]Contact, 0)
	for rows.Next() {
		var contact Contact
		err = rows.Scan(&contact.ID, &contact.Name, &contact.Email)
		require.NoError(t, rows.Err())
		contacts = append(contacts, contact)
	}
	require.NoError(t, rows.Err())
	require.Equal(t, []Contact{
		{ID: 1, Name: "foo", Email: "foo@gmail.com"},
	}, contacts)
}
