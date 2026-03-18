package mainController

import (
	"Fiber/API/2/apiSchema/commonSchema"
	"Fiber/API/2/statics/constants/status"
	"context"
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

type ResponseEnvelope struct {
	Data any `json:"data"`
}

func InitAPI(ctx *fiber.Ctx, sectionErrCode string) context.Context {
	_ = ctx
	_ = sectionErrCode
	return context.Background()
}

func FinishAPISpan(ctx *fiber.Ctx) {
	_ = ctx
}

func ParseBody(ctx *fiber.Ctx, req any) (string, int, error) {
	if err := ctx.BodyParser(req); err != nil {
		return "01", status.StatusBadRequest, err
	}

	fillHeaders(ctx, req)
	if errStr, code, err := validateBody(req); err != nil {
		return errStr, code, err
	}

	return "", status.StatusOK, nil
}

func Error(ctx *fiber.Ctx, baseErrCode string, section string, errStr string, code int, err error) error {
	return ctx.Status(code).JSON(ErrorResponse{
		ErrorCode: fmt.Sprintf("%s%s%s", baseErrCode, section, errStr),
		Message:   err.Error(),
	})
}

func Response(ctx *fiber.Ctx, res any) error {
	return ctx.Status(status.StatusOK).JSON(ResponseEnvelope{Data: res})
}

func fillHeaders(ctx *fiber.Ctx, req any) {
	refValue := reflect.ValueOf(req)
	if refValue.Kind() != reflect.Ptr || refValue.Elem().Kind() != reflect.Struct {
		return
	}

	headersField := refValue.Elem().FieldByName("Headers")
	if !headersField.IsValid() || !headersField.CanSet() || headersField.Kind() != reflect.Map {
		return
	}

	headers := map[string]string{}
	for key, value := range ctx.GetReqHeaders() {
		headers[key] = value[0]
	}
	headersField.Set(reflect.ValueOf(headers))
}

func validateBody(req any) (string, int, error) {
	refValue := reflect.ValueOf(req)
	if refValue.Kind() != reflect.Ptr || refValue.Elem().Kind() != reflect.Struct {
		return "", status.StatusOK, nil
	}

	bodyField := refValue.Elem().FieldByName("Body")
	if !bodyField.IsValid() || !bodyField.CanAddr() {
		return "", status.StatusOK, nil
	}

	validator, ok := bodyField.Addr().Interface().(interface {
		Validate(validateExtraData commonSchema.ValidateExtraData) (string, int, error)
	})
	if !ok {
		return "", status.StatusOK, nil
	}

	headers := map[string]string{}
	headersField := refValue.Elem().FieldByName("Headers")
	if headersField.IsValid() {
		if value, castOK := headersField.Interface().(map[string]string); castOK {
			headers = value
		}
	}

	return validator.Validate(commonSchema.ValidateExtraData{Headers: headers})
}
