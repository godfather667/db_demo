package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

//
// Error Function
//
func testCheck(e error) {
	if e != nil {
		panic(e)
	}
}

//
// Database 'Data Set' Constants.
//   Name Code:  The first letter of each name, in order, followed by "_db".
//   Example 'c_db' string only contains name "Charles".
//   ** The combination of "jack" and "Jacky" luckily never caused a Name conflit - Whew!
//
const c_db = "[{\"Index\":0,\"Name\":\"Charles\",\"Body\":\"Q2hhcmxlcyBEYXRh\"}]"
const cajmj_db = "[{\"Index\":0,\"Name\":\"Charles\",\"Body\":\"Q2hhcmxlcyBEYXRh\"},{\"Index\":1,\"Name\":\"Ann\",\"Body\":\"QW5uIERhdGE=\"},{\"Index\":2,\"Name\":\"Jack\",\"Body\":\"SmFjayBEYXRh\"},{\"Index\":3,\"Name\":\"Mike\",\"Body\":\"TWlrZSBEYXRh\"},{\"Index\":4,\"Name\":\"Jacky\",\"Body\":\"SmFja3kgRGF0YQ==\"}]"
const cjj_db = "[{\"Name\":\"Charles\",\"Body\":\"Q2hhcmxlcyBEYXRh\"},{\"Name\":\"Jack\",\"Body\":\"SmFjayBEYXRh\"},{\"Name\":\"Jacky\",\"Body\":\"SmFja3kgRGF0YQ==\"}]"
const cjmj_db = "[{\"Index\":0,\"Name\":\"Charles\",\"Body\":\"Q2hhcmxlcyBEYXRh\"},{\"Index\":1,\"Name\":\"Jack\",\"Body\":\"SmFjayBEYXRh\"},{\"Index\":2,\"Name\":\"Mike\",\"Body\":\"TWlrZSBEYXRh\"},{\"Index\":3,\"Name\":\"Jacky\",\"Body\":\"SmFja3kgRGF0YQ==\"}]"
const cjmjh_db = "[{\"Index\":0,\"Name\":\"Charles\",\"Body\":\"Q2hhcmxlcyBEYXRh\"},{\"Index\":1,\"Name\":\"Jack\",\"Body\":\"SmFjayBEYXRh\"},{\"Index\":2,\"Name\":\"Mike\",\"Body\":\"TWlrZSBEYXRh\"},{\"Index\":3,\"Name\":\"Jacky\",\"Body\":\"SmFja3kgRGF0YQ==\"},{\"Index\":4,\"Name\":\"Henry\",\"Body\":\"\"}]"
const cmj_db = "[{\"Name\":\"Charles\",\"Body\":\"Q2hhcmxlcyBEYXRh\"},{\"Name\":\"Mike\",\"Body\":\"TWlrZSBEYXRh\"},{\"Name\":\"Jacky\",\"Body\":\"SmFja3kgRGF0YQ==\"}]"

//const cjmjha_db = "[{\"Name\":\"Charles\",\"Body\":\"Q2hhcmxlcyBEYXRh\"},{\"Name\":\"Jack\",\"Body\":\"SmFjayBEYXRh\"},{\"Name\":\"Mike\",\"Body\":\"TWlrZSBEYXRh\"},{\"Name\":\"Jacky\",\"Body\":\"SmFja3kgRGF0YQ==\"},{\"Name\":\"Henry\",\"Body\":\"\"},{\"Name\":\"Ann\",\"Body\":\"QW5uIE5ldyBWYWx1ZQ==\"}]"

const cjmjha_db = "[{\"Index\":0,\"Name\":\"Charles\",\"Body\":\"Q2hhcmxlcyBEYXRh\"},{\"Index\":1,\"Name\":\"Jack\",\"Body\":\"SmFjayBEYXRh\"},{\"Index\":2,\"Name\":\"Mike\",\"Body\":\"TWlrZSBEYXRh\"},{\"Index\":3,\"Name\":\"Jacky\",\"Body\":\"SmFja3kgRGF0YQ==\"},{\"Index\":4,\"Name\":\"Henry\",\"Body\":\"\"},{\"Index\":5,\"Name\":\"Ann\",\"Body\":\"QW5uIE5ldyBWYWx1ZQ==\"}]"
const null_db = "null"

//
// Test Database Loader
//
func TestLoadDatabase(t *testing.T) {
	fmt.Println("Starting Database Test")
	err := os.Remove("Data.db") // Remove Current Database
	if err != nil {
		fmt.Println("Database Not present - Creating")
	}

	loadDatabase() // Load Database

	data, err := json.Marshal(xMem) // Marshall Database
	testCheck(err)

	if !reflect.DeepEqual(data, []byte(cajmj_db)) {
		t.Error("\nExpected = ", cajmj_db, "\nReturned = ", string(data))
	}
}

//
//  Test "viewHandler" Function
//
func TestViewHandler(t *testing.T) {
	fmt.Println("Starting View Handler Test")

	viewRequest, err := http.NewRequest("GET", "/view/", nil)
	if err != nil {
		t.Fatal("View NewRequest error: ", err)
	}

	listRequest, err := http.NewRequest("GET", "/view/", nil)
	if err != nil {
		t.Fatal("View NewRequest error: ", err)
	}

	allRequest, err := http.NewRequest("GET", "/view/ALL", nil)
	if err != nil {
		t.Fatal("View/ALL NewRequest error: ", err)
	}

	nameRequest, err := http.NewRequest("GET", "/view/Charles", nil)
	if err != nil {
		t.Fatal("View/Charles request error: ", err)
	}

	badRequest, err := http.NewRequest("GET", "/view/Ann", nil)
	if err != nil {
		t.Fatal("View/Ann NewRequest error: ", err)
	}

	jackRequest, err := http.NewRequest("GET", "/view/Jac", nil)
	if err != nil {
		t.Fatal("View/Jac NewRequest error: ", err)
	}

	cases := []struct {
		w                    *httptest.ResponseRecorder
		r                    *http.Request
		expectedResponseCode int
		expectedResponseBody []byte
		initial_DB           []byte
		returnedDB           []byte
	}{
		{
			w:                    httptest.NewRecorder(),
			r:                    viewRequest,
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: []byte("<h1>View: Empty Database</h1>"),
			initial_DB:           []byte(null_db),
			returnedDB:           []byte(null_db),
		},
		{
			w:                    httptest.NewRecorder(),
			r:                    listRequest,
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: []byte("<h1>Database contains the following Names:</h1><textarea nameM=\"body\" rows=\"20\" cols=\"80\">Record  0 :  Charles \n</textarea><br></form>"),
			initial_DB:           []byte(c_db),
			returnedDB:           []byte(null_db),
		},
		{
			w:                    httptest.NewRecorder(),
			r:                    allRequest,
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: []byte("<h1>Database contains the following Names:</h1><textarea nameM=\"body\" rows=\"20\" cols=\"80\">Record  0 :  Charles \n</textarea><br></form>"),
			initial_DB:           []byte(c_db),
			returnedDB:           []byte(null_db),
		},
		{
			w:                    httptest.NewRecorder(),
			r:                    nameRequest,
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: []byte("<h1>View: Charles</h1><form action=\"/load/Charles\" method=\"POST\"><textarea nameM=\"body\" rows=\"20\" cols=\"80\">Charles Data</textarea><br></form>"),
			initial_DB:           []byte(c_db),
			returnedDB:           []byte(null_db),
		},
		{
			w:                    httptest.NewRecorder(),
			r:                    badRequest,
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: []byte("<h1>View: Charles</h1><form action=\"/load/Charles\" method=\"POST\"><textarea nameM=\"body\" rows=\"20\" cols=\"80\">Charles Data</textarea><br></form>"),
			initial_DB:           []byte(c_db),
			returnedDB:           []byte(null_db),
		},
		{
			w:                    httptest.NewRecorder(),
			r:                    jackRequest,
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: []byte("<h1>View: Name Matches > 1!</h1>"),
			initial_DB:           []byte(cjj_db),
			returnedDB:           []byte(null_db),
		},
	}

	for i, c := range cases {
		// Create Desired Database
		err = json.Unmarshal(c.initial_DB, &xMem) //Reload In-Memory Copy
		testCheck(err)

		viewHandler(c.w, c.r)

		if c.expectedResponseCode != c.w.Code {
			t.Errorf("Status Code didn't match:\nExpected Value:\t%d\nReturned Value:\t%d", c.expectedResponseCode, c.w.Code)
		}

		if i == 4 { // Special Case -- Expected to Fail - Not an Error!
			if bytes.Equal(c.expectedResponseBody, c.w.Body.Bytes()) {
				t.Errorf("Body Did match:\n\t%q\n\t%q", string(c.expectedResponseBody), c.w.Body.String())
			}
		} else {
			if !bytes.Equal(c.expectedResponseBody, c.w.Body.Bytes()) {
				t.Errorf("Body didn't match:\n\t%q\n\t%q", string(c.expectedResponseBody), c.w.Body.String())
			}
		}
	}
}

//
//  Test deleteHandler" Function
//
func TestDeleteHandler(t *testing.T) {
	fmt.Println("Starting Delete Handler Test")
	// Setup Requests needed later
	allRequest, err := http.NewRequest("GET", "/delete/ALL", nil)
	if err != nil {
		t.Fatal("Delete NewRequest error: ", err)
	}

	annRequest, err := http.NewRequest("GET", "/delete/Ann", nil)
	if err != nil {
		t.Fatal("Delete NewRequest error: ", err)
	}

	henryRequest, err := http.NewRequest("GET", "/delete/Henry", nil)
	if err != nil {
		t.Fatal("Delete NewRequest error: ", err)
	}

	dupRequest, err := http.NewRequest("GET", "/delete/Jac", nil)
	if err != nil {
		t.Fatal("Delete request error: ", err)
	}

	jackRequest, err := http.NewRequest("GET", "/delete/Jack", nil)
	if err != nil {
		t.Fatal("Delete NewRequest error: ", err)
	}

	cases := []struct {
		w                    *httptest.ResponseRecorder
		r                    *http.Request
		expectedResponseCode int
		expectedResponseBody []byte
		initial_DB           []byte
		returnedDB           []byte
	}{
		{
			w:                    httptest.NewRecorder(),
			r:                    allRequest,
			expectedResponseCode: http.StatusFound,
			expectedResponseBody: []byte("<a href=\"/view/\">Found</a>.\n\n"),
			initial_DB:           []byte(null_db),
			returnedDB:           []byte(null_db),
		},
		{
			w:                    httptest.NewRecorder(),
			r:                    allRequest,
			expectedResponseCode: http.StatusFound,
			expectedResponseBody: []byte("<a href=\"/view/\">Found</a>.\n\n"),
			initial_DB:           []byte(cajmj_db),
			returnedDB:           []byte(null_db),
		},
		{
			w:                    httptest.NewRecorder(),
			r:                    annRequest,
			expectedResponseCode: http.StatusFound,
			expectedResponseBody: []byte("<a href=\"/view/\">Found</a>.\n\n"),
			initial_DB:           []byte(cajmj_db),
			returnedDB:           []byte(cjmj_db),
		},
		{
			w:                    httptest.NewRecorder(),
			r:                    henryRequest,
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: []byte("<h1>Delete: 'Henry' not found!</h1>"),
			initial_DB:           []byte(cjmj_db),
			returnedDB:           []byte(cjmj_db),
		},
		{
			w:                    httptest.NewRecorder(),
			r:                    dupRequest,
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: []byte("<h1>Delete: 'Jac' not found!</h1>"),
			initial_DB:           []byte(cjmj_db),
			returnedDB:           []byte(cjmj_db),
		},
		{
			w:                    httptest.NewRecorder(),
			r:                    jackRequest,
			expectedResponseCode: http.StatusFound,
			expectedResponseBody: []byte("<a href=\"/view/\">Found</a>.\n\n"),
			initial_DB:           []byte(cjmj_db),
			returnedDB:           []byte(cmj_db),
		},
	}

	var rMem []Page
	for _, c := range cases {
		// Create Desired Database
		err = json.Unmarshal(c.initial_DB, &xMem) //Reload In-Memory Copy
		testCheck(err)

		deleteHandler(c.w, c.r)

		// Compare to returned Database
		err = json.Unmarshal(c.returnedDB, &rMem) //Reload In-Memory Copy
		testCheck(err)
		if !reflect.DeepEqual(xMem, rMem) {
			t.Error("\nExpected Data.db        = ", xMem, "\nReceived the following  = ", rMem)
		}

		if c.expectedResponseCode != c.w.Code {
			t.Errorf("Status Code didn't match:\nExpected Value:\t%d\nReturned Value:\t%d", c.expectedResponseCode, c.w.Code)
		}

		if !bytes.Equal(c.expectedResponseBody, c.w.Body.Bytes()) {
			t.Errorf("Body Didn't match:\n\tExpected:\t%q\n\tGot:\t%q", string(c.expectedResponseBody), c.w.Body.String())
		}
	}
}

//
//  Test "editHandler" Function
//
func TestEditHandler(t *testing.T) {
	fmt.Println("Starting Edit Handler Test")
	// Setup Requests needed later
	allRequest, err := http.NewRequest("GET", "/edit/ALL", nil)
	if err != nil {
		t.Fatal("Edit NewRequest error: ", err)
	}

	blankRequest, err := http.NewRequest("GET", "/edit/", nil)
	if err != nil {
		t.Fatal("Edit NewRequest error: ", err)
	}

	henryRequest, err := http.NewRequest("GET", "/edit/Henry", nil)
	if err != nil {
		t.Fatal("Edit NewRequest error: ", err)
	}

	jackRequest, err := http.NewRequest("GET", "/edit/Jack", nil)
	if err != nil {
		t.Fatal("Edit NewRequest error: ", err)
	}

	cases := []struct {
		w                    *httptest.ResponseRecorder
		r                    *http.Request
		expectedResponseCode int
		expectedResponseBody []byte
		initial_DB           []byte
		returnedDB           []byte
	}{
		{
			w:                    httptest.NewRecorder(),
			r:                    allRequest,
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: []byte("<h1>Edit Error: 'ALL' is a Command!</h1>"),
			initial_DB:           []byte(cajmj_db),
			returnedDB:           []byte(cajmj_db),
		},
		{
			w:                    httptest.NewRecorder(),
			r:                    blankRequest,
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: []byte("<h1>Edit Error: Blank Name</h1>"),
			initial_DB:           []byte(cajmj_db),
			returnedDB:           []byte(cajmj_db),
		},
		{
			w:                    httptest.NewRecorder(),
			r:                    henryRequest,
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: []byte("<h1>Editing Henry</h1><form action=\"/save/Henry\" method=\"POST\"><textarea name=\"body\" rows=\"20\" cols=\"80\"></textarea><br><input type=\"submit\" value=\"Save\"></form>"),
			initial_DB:           []byte(cjmj_db),
			returnedDB:           []byte(cjmjh_db),
		},
		{
			w:                    httptest.NewRecorder(),
			r:                    jackRequest,
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: []byte("<h1>Editing Jack</h1><form action=\"/save/Jack\" method=\"POST\"><textarea name=\"body\" rows=\"20\" cols=\"80\">Jack Data</textarea><br><input type=\"submit\" value=\"Save\"></form>"),
			initial_DB:           []byte(cjmj_db),
			returnedDB:           []byte(cjmj_db),
		},
	}

	var rMem []Page
	for _, c := range cases {
		// Create Desired Database
		err = json.Unmarshal(c.initial_DB, &xMem) //Reload In-Memory Copy
		testCheck(err)
		editHandler(c.w, c.r)

		// Compare to returned Database
		err = json.Unmarshal(c.returnedDB, &rMem) //Reload In-Memory Copy
		testCheck(err)
		if !reflect.DeepEqual(xMem, rMem) {
			t.Error("Expected Data.db   = ", xMem)
			t.Error("Got the following  = ", rMem)
		}

		if c.expectedResponseCode != c.w.Code {
			t.Errorf("Status Code didn't match:\nExpected Value:\t%d\nReturned Value:\t%d", c.expectedResponseCode, c.w.Code)
		}

		if !bytes.Equal(c.expectedResponseBody, c.w.Body.Bytes()) {
			t.Errorf("Body Didn't match:\n\tExpected:\t%q\n\tGot:\t%q", string(c.expectedResponseBody), c.w.Body.String())
		}
	}
}

//
//  Test "saveHandler" Function
//
func TestSaveHandler(t *testing.T) {
	fmt.Println("Starting Save Handler Test")
	// Setup Requests needed later
	annRequest, err := http.NewRequest("GET", "/save/Ann", nil)
	if err != nil {
		t.Fatal("Edit NewRequest error: ", err)
	}
	emptyRequest, err := http.NewRequest("GET", "/save/", nil)
	if err != nil {
		t.Fatal("Edit NewRequest error: ", err)
	}
	badRequest, err := http.NewRequest("GET", "/save/ZZZZ", nil)
	if err != nil {
		t.Fatal("Edit NewRequest error: ", err)
	}

	cases := []struct {
		w                    *httptest.ResponseRecorder
		r                    *http.Request
		expectedResponseCode int
		expectedResponseBody []byte
		initial_DB           []byte
		returnedDB           []byte
	}{
		{
			w:                    httptest.NewRecorder(),
			r:                    annRequest,
			expectedResponseCode: http.StatusFound,
			expectedResponseBody: []byte("<a href=\"/view/Ann\">Found</a>.\n\n"),
			initial_DB:           []byte(cjmjha_db),
			returnedDB:           []byte(cjmjha_db),
		},
		{
			w:                    httptest.NewRecorder(),
			r:                    emptyRequest,
			expectedResponseCode: http.StatusFound,
			expectedResponseBody: []byte("<a href=\"/view/\">Found</a>.\n\n"),
			initial_DB:           []byte(c_db),
			returnedDB:           []byte(c_db),
		},
		{
			w:                    httptest.NewRecorder(),
			r:                    badRequest,
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: []byte("<a href=\"/view/\">OK</a>.\n\n"),
			initial_DB:           []byte(c_db),
			returnedDB:           []byte(c_db),
		},
	}

	var rMem []Page
	for _, c := range cases[2:3] {
		// Create Desired Database
		err = json.Unmarshal(c.initial_DB, &xMem) //Reload In-Memory Copy
		testCheck(err)
		saveHandler(c.w, c.r)

		// Compare to returned Database
		err = json.Unmarshal(c.returnedDB, &rMem) //Reload In-Memory Copy
		testCheck(err)
		if !reflect.DeepEqual(xMem, rMem) {
			t.Error("\nExpected Data.db   = ", xMem, "\nReceived the following  = ", rMem)
		}

		if c.expectedResponseCode != c.w.Code {
			t.Errorf("Status Code didn't match:\nExpected Value:\t%d\nReturned Value:\t%d", c.expectedResponseCode, c.w.Code)
		}

		if !bytes.Equal(c.expectedResponseBody, c.w.Body.Bytes()) {
			t.Errorf("Body Didn't match:\n\tExpected:\t%q\n\tGot:\t%q", string(c.expectedResponseBody), c.w.Body.String())
		}
	}
}
