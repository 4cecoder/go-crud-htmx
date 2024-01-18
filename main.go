// https://www.loom.com/share/ca639d29fad54600b1478957e9afc9ff?sid=8a9de4e1-091a-4fdb-a796-e665bd326f3a
package main // This file is part of the main package

import (
	"fmt"                                      // The fmt package implements formatted I/O with functions analogous to C's printf and scanf.
	"github.com/gin-gonic/gin"                 // Gin is a HTTP web framework written in Go. It features a Martini-like API with much better performance.
	"github.com/jinzhu/gorm"                   // GORM is a developer friendly ORM library for handling relational databases. It is based on the concept of Object-Relational Mapping (ORM).
	_ "github.com/jinzhu/gorm/dialects/sqlite" // This is the SQLite dialect for GORM. The underscore means we're importing this package for its side-effects only (in this case, registering the SQLite driver).
	"log"                                      // The log package provides simple logging services. It defines a type, Logger, with methods for formatting output.
	"net/http"                                 // The net/http package provides HTTP client and server implementations.
	"os"                                       // The os package provides a platform-independent interface to operating system functionality.
)

/*
Go is a statically typed, compiled language developed by Google.

It's simple and efficient, ideal for fast and reliable software. Its syntax is clean and easy to learn.

Go has a ":=" operator for declaring and initializing variables in one line.
For example, "name := 'John'" declares a variable "name" and assigns it the value 'John'.
The type is inferred by the compiler, so "name" is of type string.

Go also has a "_" identifier, known as the blank identifier. It's used to ignore values.
For example, in "result, _ := fetchData()", "fetchData" is a function that returns two values.
If we're only interested in the first value, we can use "_" to ignore the second one.
So "result" will hold the first return value, and the second one is ignored.
*/

// main is the entry point of the application.
func main() {
	// Create a default Gin engine. Gin is a web framework written in Go.
	r := gin.Default()

	// Initialize the database connection by calling the initDB function.
	db := initDB()

	// Set up the routes for the application by calling the setupRoutes function.
	setupRoutes(r, db)

	// Print a message to the console indicating that the server is starting.
	fmt.Println("Starting server on port 8080...")

	// Print a message to the console indicating where the frontend and API can be accessed.
	fmt.Println("Go to http://localhost:8080/frontend to see the frontend.")
	fmt.Println("Go to http://localhost:8080/users to see the API.")

	// Start the server on port 8080. If there's an error starting the server, log the error and terminate the program.
	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// initDB is a function that initializes the database connection.
// It takes no parameters and returns a pointer to a gorm.DB object (which represents the database connection).
func initDB() *gorm.DB {
	// Call gorm.Open to open a new database connection.
	// The first argument is the name of the database driver ("sqlite3" in this case),
	// and the second argument is the name of the database file ("test.db" in this case).
	// If there's an error opening the database connection, log the error and terminate the program.
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Call db.AutoMigrate with the User struct to automatically create the users table in the database
	// based on the User struct. If the table already exists, it will be updated to include any new fields
	// in the User struct.
	db.AutoMigrate(&User{})

	// Log a success message and return the database connection.
	fmt.Println("Database initialized")
	return db
}

// User is a struct that represents a user in the application.
// It is also a gorm.Model, which means it includes fields for ID, CreatedAt, UpdatedAt, and DeletedAt.
// These fields are automatically managed by GORM when you perform database operations.
// The json tags are used to define the JSON encoding and decoding behavior for the struct fields.
// The gorm tags are used to define the database schema details for the struct fields.
type User struct {
	gorm.Model        // This embeds the gorm.Model struct, which includes fields ID, CreatedAt, UpdatedAt, and DeletedAt.
	Name       string `json:"name"`                                        // The Name field represents the user's name. It is a string and will be encoded/decoded as "name" in JSON.
	Email      string `gorm:"type:varchar(100);unique_index" json:"email"` // The Email field represents the user's email. It is a string and will be encoded/decoded as "email" in JSON. In the database, it is a varchar of length 100 and must be unique.
	Password   string `json:"password"`                                    // The Password field represents the user's password. It is a string and will be encoded/decoded as "password" in JSON.
}

// UserResponse is a struct that represents the response sent to the client when they request user data.
// It includes the user's ID, name, and email.
// The json tags are used to define the JSON encoding and decoding behavior for the struct fields.
type UserResponse struct {
	ID    uint   `json:"id"`    // The ID field represents the user's ID. It is an unsigned integer and will be encoded/decoded as "id" in JSON.
	Name  string `json:"name"`  // The Name field represents the user's name. It is a string and will be encoded/decoded as "name" in JSON.
	Email string `json:"email"` // The Email field represents the user's email. It is a string and will be encoded/decoded as "email" in JSON.
}

// fetchUsers is a function that fetches all users from the database.
// It takes one parameter: a pointer to a gorm.DB object (which represents the database connection).
func fetchUsers(db *gorm.DB) []UserResponse {
	// Declare a slice of User to hold the users data from the database.
	var users []User

	// Fetch all users from the database.
	db.Find(&users)

	// Declare a slice of UserResponse to hold the response data.
	var userResponses []UserResponse

	// Loop over the users slice.
	for _, user := range users {
		// For each user, create a UserResponse with the user's ID, name, and email,
		// and append it to the userResponses slice.
		userResponses = append(userResponses, UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	// Return the userResponses slice.
	return userResponses
}

// createUser is a function that creates a new user in the database.
// It takes two parameters: a pointer to a gorm.DB object (which represents the database connection)
// and a pointer to a gin.Context object (which represents the current HTTP request context).
func createUser(db *gorm.DB, c *gin.Context) {
	// Declare a variable of type User to hold the user data from the request.
	var user User

	// Get the Content-Type header from the request.
	contentType := c.GetHeader("Content-Type")

	// Check if the Content-Type is not application/json.
	// If it's not, log an error message and return a 400 Bad Request response with an error message.
	if contentType != "application/json" {
		log.Println("Invalid content type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type. Only application/json is accepted"})
		return
	}

	// Bind the JSON body of the request to the user variable.
	// If there's an error (e.g., the request body is not valid JSON), log the error and return a 400 Bad Request response with an error message.
	err := c.BindJSON(&user)
	if err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if any of the required fields (name, email, password) are empty.
	// If any of them are, log an error message and return a 400 Bad Request response with an error message.
	if user.Name == "" || user.Email == "" || user.Password == "" {
		log.Println("All fields are required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	// Create a new record in the database with the user data.
	// If there's an error (e.g., the email is not unique), log the error and return a 400 Bad Request response with an error message.
	db.Create(&user)
	if db.Error != nil {
		log.Printf("Failed to create user: %v", db.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}

	// If everything went well, log a success message and return a 200 OK response with a success message.
	fmt.Println("User created")
	c.JSON(http.StatusOK, gin.H{"message": "User created"})
}

// updateUser is a function that updates an existing user in the database.
// It takes three parameters: a pointer to a gorm.DB object (which represents the database connection),
// a pointer to a gin.Context object (which represents the current HTTP request context), and the ID of the user to update.
func updateUser(db *gorm.DB, c *gin.Context, id string) UserResponse {
	// Declare a variable of type User to hold the existing user data from the database.
	var user User

	// Fetch the user with the given ID from the database.
	// If the user does not exist, gorm will return an error.
	db.First(&user, id)

	// Declare a variable of type User to hold the new user data from the request.
	var newUser User

	// Bind the JSON body of the request to the newUser variable.
	// If there's an error (e.g., the request body is not valid JSON), log the error and return an empty UserResponse.
	err := c.BindJSON(&newUser)
	if err != nil {
		log.Printf("Error binding JSON: %v", err)
		return UserResponse{}
	}

	// Update the user's name and email with the new data.
	user.Name = newUser.Name
	user.Email = newUser.Email

	// Save the updated user data to the database.
	// If there's an error (e.g., the email is not unique), log the error and return an empty UserResponse.
	db.Save(&user)
	if db.Error != nil {
		log.Printf("Failed to update user: %v", db.Error)
		return UserResponse{}
	}

	// If everything went well, log a success message and return a UserResponse with the updated user data.
	fmt.Println("User updated")
	return UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

// deleteUser is a function that deletes an existing user from the database.
// It takes two parameters: a pointer to a gorm.DB object (which represents the database connection),
// and the ID of the user to delete.
func deleteUser(db *gorm.DB, id string) {
	// Declare a variable of type User to hold the existing user data from the database.
	var user User

	// Fetch the user with the given ID from the database.
	// If the user does not exist, gorm will return an error.
	db.First(&user, id)

	// Delete the user from the database.
	// If there's an error (e.g., the user does not exist), log the error and return.
	db.Delete(&user)
	if db.Error != nil {
		log.Printf("Failed to delete user: %v", db.Error)
		return
	}

	// If everything went well, log a success message.
	fmt.Println("User deleted")
}

// getUsersHandler is a function that returns a Gin HandlerFunc for fetching all users.
// It takes a pointer to a gorm.DB object (which represents the database connection) as a parameter.
// The returned HandlerFunc calls the fetchUsers function with the database connection,
// then sends a JSON response with the fetched users data.
func getUsersHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		users := fetchUsers(db)
		c.JSON(200, users)
	}
}

// createUserHandler is a function that returns a Gin HandlerFunc for creating a user.
// It takes a pointer to a gorm.DB object (which represents the database connection) as a parameter.
// The returned HandlerFunc calls the createUser function with the database connection and the Gin Context object.
func createUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		createUser(db, c)
	}
}

// getUserByIDHandler is a function that returns a Gin HandlerFunc for fetching a user by their ID.
// It takes a pointer to a gorm.DB object (which represents the database connection) as a parameter.
// The returned HandlerFunc calls the fetchUserByID function with the database connection and the user ID from the request parameters.
// It then sends a JSON response with the fetched user data.
func getUserByIDHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := fetchUserByID(db, c.Params.ByName("id"))
		c.JSON(200, user)
	}
}

// updateUserHandler is a function that returns a Gin HandlerFunc for updating a user.
// It takes a pointer to a gorm.DB object (which represents the database connection) as a parameter.
// The returned HandlerFunc calls the updateUser function with the database connection, the Gin Context object, and the user ID from the request parameters.
// It then sends a JSON response with the updated user data.
func updateUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := updateUser(db, c, c.Params.ByName("id"))
		c.JSON(200, user)
	}
}

// Handler for DELETE /users/:id
func deleteUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		deleteUser(db, c.Params.ByName("id"))
		c.JSON(200, gin.H{"message": "User deleted"})
	}
}

// setupRoutes is a function that sets up the routes for the application.
// It takes two parameters: a pointer to a gin.Engine object (which represents the HTTP router),
// and a pointer to a gorm.DB object (which represents the database connection).
func setupRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/users", getUsersHandler(db))
	r.POST("/users", createUserHandler(db))
	r.GET("/users/:id", getUserByIDHandler(db))
	r.PUT("/users/:id", updateUserHandler(db))
	r.DELETE("/users/:id", deleteUserHandler(db))

	// Call the frontend function to set up the route for the frontend of the application.
	frontend(r)
}

// frontend is a function that sets up the route for the frontend of the application.
// It takes one parameter: a pointer to a gin.Engine object (which represents the HTTP router).
func frontend(r *gin.Engine) {
	// Define a GET route for "/frontend".
	r.GET("/frontend", func(c *gin.Context) {
		// Set the Content-Type header of the response to "text/html".
		c.Header("Content-Type", "text/html")

		// Read the contents of the "index.html" file.
		// If there's an error (e.g., the file does not exist), log the error and return.
		html, err := os.ReadFile("index.html")
		if err != nil {
			log.Fatalf("Error reading index.html: %v", err)
		}

		// Send a 200 OK response with the contents of the "index.html" file.
		c.String(http.StatusOK, string(html))
	})
}

// fetchUserByID is a function that fetches a user from the database by their ID.
// It takes two parameters: a pointer to a gorm.DB object (which represents the database connection),
// and the ID of the user to fetch.
func fetchUserByID(db *gorm.DB, id string) UserResponse {
	// Declare a variable of type User to hold the user data from the database.
	var user User

	// Fetch the user with the given ID from the database.
	// If the user does not exist, gorm will return an error.
	db.First(&user, id)

	// Check if there was an error fetching the user.
	// If there was, log the error and return an empty UserResponse.
	if db.Error != nil {
		log.Printf("Failed to fetch user by ID: %v", db.Error)
		return UserResponse{}
	}

	// If everything went well, log a success message and return a UserResponse with the user data.
	fmt.Println("User fetched by ID")
	return UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}
