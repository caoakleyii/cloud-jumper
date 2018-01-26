package handler

import (
	"time"

	"github.com/caoakleyii/cloud-jumper/src/cache"

	"github.com/caoakleyii/cloud-jumper/src/hasher"
)

/*
	2. Hash and Encode Passwords over HTTP

	Program listens for a POST to /hash with a body form including password
	Response is a base64 encoded hash of the password
	http.Serve creates new threads for each incoming
	request allowing our app to still handle other request while this does work
*/

// PostPassword handler for the POST "/hash" endpoint
// Hashes a password and returns the value
func postPassword(ctx *Context) {
	time.Sleep(time.Second * 5)

	p := ctx.Request.FormValue("password")

	// Check "validation" on the incoming password
	if p == "" {
		ctx.String(400, "Bad Request")
		return
	}

	// Always generate a secure salt for your hash
	salt, err := hasher.GenerateRandomString(32)

	if err != nil {
		ctx.String(500, "Internal Server Error")
		return
	}

	// Hash and Encode the password
	p = hasher.Sha512HashBase64Encode(p, salt)
	ctx.String(201, p)
}

/*
	4. Hash End-Point Returns Identifier

	Modified version of function "postPassword"
	hashes the password, creates a subroutine that takes ~5 seconds to
	store the password, meanwhile returning a CREATED response with the
	password id.
*/

// PostPassword handler for the POST "/hash" endpoint
// Hashes a password and returns the id and after 5 seconds stores it
// in-memory
func PostPassword(ctx *Context) {
	p := ctx.Request.FormValue("password")

	// Check "validation" the incoming password
	if p == "" {
		ctx.String(400, "Bad Request")
		return
	}

	// Always generate a secure salt for your hash
	salt, err := hasher.GenerateRandomString(32)

	if err != nil {
		ctx.String(500, "Internal Server Error")
		return
	}

	// generate a short id
	id, err := hasher.GenerateRandomString(8)

	if err != nil {
		ctx.String(500, "Internal Server Error")
		return
	}

	// store the data on a seperate thread 5 seconds later
	go func(id string, m map[string]string) {
		time.Sleep(time.Second * 5)

		// hash, encode and store the password
		p = hasher.Sha512HashBase64Encode(p, salt)
		m[id] = p
	}(id, cache.InMemoryPasswordStorage)

	// return 201 created with the id
	ctx.String(201, id)
	return
}

/*
	5. GET a Hashed Password

	Returns a hashed password, using the /hash/{id} request
	retrieves the password from an in-memory map
*/

// GetPassword handler function that returns an
// OK response with the hashed password
func GetPassword(ctx *Context) {
	id := ctx.Param("id")

	p := cache.InMemoryPasswordStorage[id]

	if p == "" {
		ctx.String(404, "Password Not Found")
		return
	}

	ctx.String(200, p)
	return
}
