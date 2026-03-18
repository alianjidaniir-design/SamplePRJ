package student

import (
	"Fiber/API/2/apiSchema/commonSchema"
	"Fiber/API/2/apiSchema/studentsSchema"
	"Fiber/API/2/controllers/mainController"
	"Fiber/API/2/models/repositories"
	"Fiber/API/2/statics/constants/controllerBaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func Create(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "11")

	defer mainController.FinishAPISpan(ctx)

	req := commonSchema.BaseRequest[studentsSchema.CreateUserRequest]{}

	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerBaseErrCode.UserErrCode, "01", errStr, code, err)
	}
	res, errStr, code, err := repositories.UserRepo.Create(spanCtx, req)
	if err != nil {
		return mainController.Error(ctx, controllerBaseErrCode.UserErrCode, "02", errStr, code, err)

	}
	return mainController.Response(ctx, res)
}
