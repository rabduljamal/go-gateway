package metabase_handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/rabduljamal/gateway-snip/config"
)

type DataInput struct {
	Question int         `json:"question"`
	Params   interface{} `json:"params"`
}

func GetMetabases(c *fiber.Ctx) error {

	METABASE_SECRET_KEY := config.Config("METABASE_SECRET_KEY")
	METABASE_SITE_URL := config.Config("METABASE_SITE_URL")

	dataInput := new(DataInput)
	if err := c.BodyParser(dataInput); err != nil {
		return err // Handle parsing error
	}

	// Create the payload
	expiration := time.Now().Add(10 * time.Minute).Unix()
	payload := jwt.MapClaims{
		"resource": map[string]interface{}{
			"question": dataInput.Question,
		},
		"params": dataInput.Params,
		"exp":    expiration,
	}

	fmt.Println(payload)

	// Generate the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte(METABASE_SECRET_KEY))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("JWT signing failed")
	}

	// Build the URL
	url := METABASE_SITE_URL + "/api/embed/card/" + tokenString + "/query/json"

	// Perform the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{"Error fetching data": err.Error()})
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{"Error fetching data": "Non-OK response"})
	}

	// Parse the response body
	var data interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{"Error decoding JSON": err.Error()})
	}

	return c.JSON(data)
}
