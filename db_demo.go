// db_demo - Demostrates a minimal "NOSQL" database and Web Interface.
// This initial program was done for a job interview and is documented
// SEE: http://www.hawthorne-press.com/GO_Short_Interview_Questions_Explained_In_Color.pdf
// Added additional features and Annotations
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var xMem []Page // In Memory Database File

type Page struct { // Database Page
	Index int    // Index of Database Page
	Name  string // KEY: Name as Search Key
	Body  []byte // VALUE: Data associated with the Key
}

func main() {
	fmt.Println("Starting Database Server")
	//	http.HandleFunc("/", slashHandler) // Display Help Commands

	loadDatabase()                         // Load Database
	http.HandleFunc("/", slashHandler)     // Display Help Commands
	http.HandleFunc("/view/", viewHandler) // Setup Handler Functions
	http.HandleFunc("/exit/", exitHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/delete/", deleteHandler)
	http.ListenAndServe(":8080", nil) // Setup up Server to listen on port 8080
}

//
// Load Database
//
func loadDatabase() {
	//	var xMem []Page // Temporary Blank in-memory Database
	data, err := ioutil.ReadFile("Data.db") // Load Database
	if err != nil {                         // If missing - Create
		//
		// Create the Initial Database with Test Data - Remove appends below for empty database
		//
		xMem = append(xMem, Page{Index: 0, Name: "Charles", Body: []byte("Charles Data")})
		xMem = append(xMem, Page{Index: 1, Name: "Ann", Body: []byte("Ann Data")})
		xMem = append(xMem, Page{Index: 2, Name: "Jack", Body: []byte("Jack Data")})
		xMem = append(xMem, Page{Index: 3, Name: "Mike", Body: []byte("Mike Data")})
		xMem = append(xMem, Page{Index: 4, Name: "Jacky", Body: []byte("Jacky Data")})

		data, err = json.Marshal(xMem) // Marshall Database
		check("Marshalling Failed", err)
		_, err := os.Create("Data.db") // Create Database
		check("Create File Failed", err)
		writeData(data)                   // Write Database
		err = json.Unmarshal(data, &xMem) //Reload In-Memory Copy
		check("Unmarshal Failed", err)
	}
}

//
// check error("error msg", err)
//
func check(s string, e error) {
	if e != nil {
		fmt.Println(s)
		panic(e)
	}
}

//
// Write Data Set to Disk
//
func writeData(data []byte) { // Write "Mashalled" data to external device
	err := ioutil.WriteFile("Data.db", data, 0644)
	check("Write File Failed", err) // Error Check -- Panic if write fails
}

// Find Name Function - Locates by string.Contains.
// If more than one name matches, it searches again with string.Compare Function (equality match)
//
func findName(xMem []Page, name string) (Page, bool) {
	var retPg Page // Create Variables

	dup := 0
	for i, key := range xMem { // Search using string.Contains method
		if strings.Contains(key.Name, name) {
			dup++ // Increment variable dup for each match
			retPg = xMem[i]
		}
	}
	if dup != 1 { // Duplicate Search Result Check!
		for i, key := range xMem { // Find exact match!
			if strings.Compare(key.Name, name) == 0 {
				retPg = xMem[i]    // If found - Store returned page
				return retPg, true // Return Result and True (exact match found)
			}
		}
		return retPg, false // If no exact match return False
	}
	return retPg, true // Return Result and True (only one "contains" Result found)
}

//
// Find Name Function - Locates by string.Compare "equality" match
//
func findExactName(xMem []Page, name string) (Page, bool) {
	var retPg Page             // Create Variables
	for i, key := range xMem { // Search in-memory database for an exact match
		if strings.Compare(key.Name, name) == 0 {
			retPg = xMem[i]    // Match Found
			return retPg, true // Return Page and bool return value of True
		}
	}
	return retPg, false // No match found -- Return bool return value of False
}

//
// Slash Handler --
//
// localhost:8080            -- Lists Help Information
//
func slashHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>This Database responds to the following commands</h1>"+
		"<h2>localhost:8080/&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;This Message<br>"+
		"localhost:8080/help/&emsp;&emsp;&emsp;&emsp;This Message<br>"+
		"localhost:8080/exit/&emsp;<br>"+
		"localhost:8080/view/name/&emsp;(name Optional)<br>"+
		"localhost:8080/edit/name/&emsp;<br>"+
		"localhost:8080/delete/name/&emsp;  <br></h2>")
	return
}

//
// Exit Handler --
//
//
func exitHandler(w http.ResponseWriter, r *http.Request) {
	os.Exit(-1)
}

//
// View Handler --
//
// localhost:8080/view/     -- Lists all names in the database
//
// localhost:8080/view/name  -- Displays the Page for "name"
//
func viewHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/view/"):]
	// Extract "name" from URL path
	var body string

	// Create Variables
	if len(name) <= 0 {
		// Check for /view without name
		if len(xMem) <= 0 {
			// Empty Database Check
			body = fmt.Sprintln("")
			// Create Empty TextArea
		} else {
			// Page Validation Check --
			for i := 0; i < len(xMem); i++ {
				if i != xMem[i].Index {
					// Make sure in-memory index equals internal page index
					// Database corruption error!
					error := fmt.Sprintln("Database Indexing Failure: Name = ", xMem[i].Name)
					fmt.Fprintf(w, "<h1>View: %s</h1>", error)
					time.Sleep(time.Second)
					return
				}
				// Display&emsp; Page Elements in Testarea of <form>
				body += fmt.Sprintln("Record ", xMem[i].Index, ": ", xMem[i].Name, string(""))
			}
			// Send constructed Display to the Client!
			fmt.Fprintf(w, "<h1>Database contains the following Names:</h1>"+
				"<textarea nameM=\"body\" rows=\"20\" cols=\"80\">%s</textarea><br>"+
				"</form>", body)
			return
		}
	}
	// Display Empty Database Message
	if len(xMem) <= 0 {
		fmt.Fprintf(w, "<h1>View: %s</h1>", "Empty Database")
		return
	}
	// Display "ALL" Error (Only for /delete/ALL)
	if name == "ALL" {
		fmt.Fprintf(w, "<h1>View Error: %s</h1>", "'ALL' is a Command!")
		return
	}
	// Handle Display of a "Named" Page
	p, ok := findName(xMem, name)
	if !ok {
		// Too many matches or Name not found
		if len(p.Name) > 0 {
			fmt.Fprintf(w, "<h1>View: %s</h1>", "Name Matches > 1!")
			return
		}
		fmt.Fprintf(w, "<h1>View: %s</h1>", "Name not found!")
		return
	}

	fmt.Fprintf(w, "<h1>View: %s</h1>"+
		// Create <form> for viable “name” result
		"<form action=\"/load/%s\" method=\"POST\">"+
		"<textarea nameM=\"body\" rows=\"20\" cols=\"80\">%s</textarea><br>"+
		"</form>",
		p.Name, p.Name, p.Body)
}

////
// SaveHandler helper to create and store a new page in the database
//
func (p *Page) save() {
	np := Page{Index: p.Index, Name: p.Name, Body: p.Body} // Create a database Page
	xMem[p.Index] = np                                     // Append it to the in-memory database
	data, err := json.Marshal(xMem)                        // Marshal the database
	check("Marshalling Failed", err)                       // Check for error
	writeData(data)                                        // Write database to the disk
	return                                                 // Return
}

// Save Handler Function -- Should not be used by the Client
//
func saveHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/save/"):] // Get "name" value if present
	if len(name) <= 0 {                // If no name - redirect to /view
		http.Redirect(w, r, "/view/", http.StatusFound)
	}
	pg, ok := findName(xMem, name) // Find "name"
	if !ok {                       // If error - report it and panic
		fmt.Println("Name Lookup Failed")
		panic("error")
	}
	body := r.FormValue("body") // Get <form> value for "body"
	p := &Page{Index: pg.Index, Name: name, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+name, http.StatusFound) // Redirect to /view/name
}

//
// Edit Handler
//
// localhost:8080/edit/name
//
func editHandler(w http.ResponseWriter, r *http.Request) {
	var np Page
	// Create Variables
	name := r.URL.Path[len("/edit/"):]
	// Extract Name Portion
	if name == "ALL" {
		// Name = ALL -- Print Error
		fmt.Fprintf(w, "<h1>Edit Error: %s</h1>", "'ALL' is a Command!")
		return
	}
	p, ok := findName(xMem, string(name))
	if !ok { // Find Name
		// If no name -
		//Create name with empty body
		np = Page{Index: len(xMem), Name: name, Body: []byte("")}
		xMem = append(xMem, np)         // Append the in-memory database
		data, err := json.Marshal(xMem) // Marshal xMem ==> Data set
		check("Marshalling Failed", err)
		writeData(data)                      // Write Data Set to disk
		p, ok = findName(xMem, string(name)) // Find newly created name!
		if !ok {
			fmt.Println("Update Failure")
			// Notify of failure
			return
		}
	}
	fmt.Fprintf(w, "<h1>Editing %s</h1>"+
		// Build Form and send to client
		"<form action=\"/save/%s\" method=\"POST\">"+
		"<textarea name=\"body\" rows=\"20\" cols=\"80\">%s</textarea><br>"+
		"<input type=\"submit\" value=\"Save\">"+
		"</form>",
		p.Name, p.Name, p.Body)
}

//
// Delete Handler
//     Delete/ALL
//     Delete/Name
//
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	var zMem []Page // Empty Database
	name := r.URL.Path[len("/delete/"):]
	if name == "ALL" {
		// Process ALL
		data, err := json.Marshal(zMem)
		// Marshal empty Database
		check("Marshalling Failed", err)
		writeData(data)
		// Write Data Set to disk
		err = json.Unmarshal(data, &xMem) //Reload In-Memory Copy
		check("Unmarshal Failed", err)
		http.Redirect(w, r, "/view/", http.StatusFound)
		// Redirect to /view
	}
	p, ok := findExactName(xMem, string(name))
	// Not ALL - Find Name
	if !ok {
		// Report Failure
		fmt.Fprintf(w, "<h1>View: '%s' %s</h1>", name, "not found!")
		return
	} else {
		// Deletion has to be done this way to insure that internal index updated!
		i := 0
		for j, v := range xMem {
			// Walk old Database and Create new Database
			if j != p.Index {
				zMem = append(zMem, Page{Index: i, Name: v.Name, Body: v.Body})
				i++
			}
		}
		data, err := json.Marshal(zMem)
		// Marshal new Database
		check("Marshalling Failed", err)
		writeData(data)
		// Write to disk
		err = json.Unmarshal(data, &xMem)
		//Reload In-Memory Copy
		check("Unmarshal Failed", err)
		http.Redirect(w, r, "/view/", http.StatusFound)
		// Redirect to /view
	}
}
