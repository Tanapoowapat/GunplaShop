package entities

import (
	"github.com/Tanapoowapat/GunplaShop/pkg/gunplalogger"
	"github.com/gofiber/fiber/v2"
)

type IResponse interface {
	Sucess(code int, data any) IResponse
	Error(code int, tractId, message string) IResponse
	Res() error
}

type Response struct {
	StatusCode int
	Data       any
	ErrorRes   *ErrorResponse
	Context    *fiber.Ctx
	IsError    bool
}

type ErrorResponse struct {
	TractId string `json:"tractId"`
	Message string `json:"message"`
}

func NewResponse(context *fiber.Ctx) IResponse {
	return &Response{
		Context: context,
	}
}

func (r *Response) Sucess(code int, data any) IResponse {
	r.StatusCode = code
	r.Data = data
	r.IsError = false
	// save log
	gunplalogger.NewGunplaLogger(r.Context, &r.Data, r.StatusCode).Print().Save()
	return r
}

func (r *Response) Error(code int, tractId, message string) IResponse {
	r.StatusCode = code
	r.ErrorRes = &ErrorResponse{
		TractId: tractId,
		Message: message,
	}
	r.IsError = true
	// save log
	gunplalogger.NewGunplaLogger(r.Context, r.ErrorRes, r.StatusCode).Print().Save()
	return r
}

func (r *Response) Res() error {
	if r.IsError {
		return r.Context.Status(r.StatusCode).JSON(r.ErrorRes)
	}
	return r.Context.Status(r.StatusCode).JSON(r.Data)
}

type PaginateRes struct {
	Data       any `json:"data"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPage  int `json:"total_page"`
	TotalItems int `json:"total_item"`
}
