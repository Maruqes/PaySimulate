package PaySimulate

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
)

type CardDetails struct {
	CardNumber  string
	HoldersName string
	ExpDate     string
	CVV         string
	BillingAddr string
	Country     string
}

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var qrCodeMap = make(map[string]bool)

func generateID(length int) string {
	if length%4 != 0 {
		panic("length must be a multiple of 4 generateID")
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var sb strings.Builder

	for i := 0; i < length; i++ {
		sb.WriteByte(charset[rng.Intn(len(charset))])
		if (i+1)%4 == 0 && i < length-1 {
			sb.WriteByte('-')
		}
	}

	return sb.String()
}

func cardTest(card CardDetails) (string, error) {
	if len(card.CardNumber) != 16 {
		return "", errors.New("card number must be 16 digits")
	}
	if len(card.HoldersName) < 3 {
		return "", errors.New("name must be at least 3 characters")
	}
	if len(card.ExpDate) != 5 {
		return "", errors.New("invalid expiration date")
	}
	if len(card.CVV) != 3 {
		return "", errors.New("invalid CVV")
	}

	if len(card.BillingAddr) < 5 {
		return "", errors.New("invalid billing address")
	}

	if len(card.Country) != 2 {
		return "", errors.New("invalid country code")
	}

	return generateID(32), nil
}

// Example function in PaySimulate package
func PaymentSuccess(card CardDetails, approved bool) (string, error) {

	id, err := cardTest(card)
	if err != nil {
		return "", err
	}

	return id, nil
}

func apiPostPayment(c *fiber.Ctx) error {
	card := new(CardDetails)
	if err := c.BodyParser(card); err != nil {
		return err
	}

	id, err := PaymentSuccess(*card, true)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"id": id,
	})
}

func apiPostPayQRcode(c *fiber.Ctx) error {

	qrId := generateID(16)

	fullURL := c.Protocol() + "://" + c.Hostname() + c.OriginalURL()

	fullURL += "Pay?qrId=" + qrId

	png, err := qrcode.Encode(fullURL, qrcode.Medium, 256)
	if err != nil {
		log.Println("Failed to generate QR code:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error generating QR code")
	}

	c.Set("Content-Type", "image/png")

	qrCodeMap[qrId] = false
	fmt.Println(qrId)
	fmt.Println(fullURL)
	return c.Send(png)
}

func qrCodePay(c *fiber.Ctx) error {

	qrId := c.Query("qrId", "")
	if qrId == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing 'qrId' query parameter")
	}

	if _, ok := qrCodeMap[qrId]; !ok {
		return c.Status(fiber.StatusNotFound).SendString("QR code not found")
	}

	qrCodeMap[qrId] = true
	return c.SendString("Payment successful")
}

func checkQRCodeStatus(c *fiber.Ctx) error {
	qrId := c.Query("qrId", "")
	if qrId == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing 'qrId' query parameter")
	}

	if _, ok := qrCodeMap[qrId]; !ok {
		return c.Status(fiber.StatusNotFound).SendString("QR code not found")
	}

	if qrCodeMap[qrId] {
		return c.SendString("Payment successful")
	}

	return c.SendString("Payment pending")
}

func StartWebServer(port string) {
	//test if port is like :"port"
	if port[0] != ':' {
		port = ":" + port
	}

	app := fiber.New()

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "pong",
		})
	})

	app.Post("/payment", apiPostPayment)
	app.Get("/qrcode", apiPostPayQRcode)
	app.Get("/qrcodePay", qrCodePay)
	app.Get("/checkQRCodeStatus", checkQRCodeStatus)

	app.Listen(port)
}
