package gunplalogger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Tanapoowapat/GunplaShop/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type IGunplaLogger interface {
	Print() IGunplaLogger
	Save()
	SetQuery(c *fiber.Ctx)
	SetBody(c *fiber.Ctx)
	SetResponse(res any)
}

type gunplaLogger struct {
	Time       string `json:"time"`
	Ip         string `json:"ip"`
	Method     string `json:"method"`
	StatusCode int    `json:"status_code"`
	Path       string `json:"path"`
	Query      any    `json:"query"`
	Body       any    `json:"body"`
	Response   any    `json:"response"`
}

func NewGunplaLogger(c *fiber.Ctx, res any, code int) IGunplaLogger {
	log := &gunplaLogger{
		Time:       time.Now().Local().Format("2006-01-02 15:04:05"),
		Ip:         c.IP(),
		Method:     c.Method(),
		Path:       c.Path(),
		StatusCode: code,
	}
	log.SetQuery(c)
	log.SetBody(c)
	log.SetResponse(res)
	return log
}

func (l *gunplaLogger) Print() IGunplaLogger {
	utils.Debug(l)
	return l
}

func (l *gunplaLogger) Save() {
	data := utils.Output(l)
	filename := fmt.Sprintf("./assets/logs/gunplaLogerr-%s.json", strings.ReplaceAll(time.Now().Local().Format("2006-01-02"), "-", ""))
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	file.WriteString(string(data) + "\n")
}

func (l *gunplaLogger) SetQuery(c *fiber.Ctx) {
	var body any
	if err := c.QueryParser(&body); err != nil {
		log.Printf("query parser error: %v", err)
	}
	l.Query = body
}

func (l *gunplaLogger) SetBody(c *fiber.Ctx) {
	var body any
	if err := c.BodyParser(&body); err != nil {
		log.Printf("body parser error: %v", err)
	}
	switch l.Path {
	case "v1/users/signup":
		l.Body = "sensitive data..."
	default:
		l.Body = body
	}

}

func (l *gunplaLogger) SetResponse(res any) {
	l.Response = res
}
