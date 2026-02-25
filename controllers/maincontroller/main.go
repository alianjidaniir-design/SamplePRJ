package maincontroller

import (
	"context"
	"fmt"
	"reflect"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/models/datamodel"
	"github.com/alianjidaniir-design/SamplePRJ/statics/constants/status"
	"github.com/gofiber/fiber/v2"
)

type errorResponse struct {
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}

type responseEnvelope struct {
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
		return "00", status.StatusBadRequest, err
	}

	fillHeaders(ctx, req)
	if errStr, code, err := validateBody(req); err != nil {
		return errStr, code, err
	}

	return "", status.StatusOK, nil
}

func Error(ctx *fiber.Ctx, baseErrCode string, section string, errStr string, code int, err error) error {
	return ctx.Status(code).JSON(errorResponse{
		ErrorCode: fmt.Sprintf("%s%s%s", baseErrCode, section, errStr),
		Message:   err.Error(),
	})
}

func Response(ctx *fiber.Ctx, res any) error {
	return ctx.Status(status.StatusOK).JSON(responseEnvelope{Data: res})
}

func GetUser(ctx *fiber.Ctx) datamodel.User {
	_ = ctx
	return datamodel.User{ID: 11, Username: "demo-user"}
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
		headers[key] = value
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
