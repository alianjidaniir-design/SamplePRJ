package user

import (
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/userSchema"
	mainController "github.com/alianjidaniir-design/SamplePRJ/controllers/maincontroller"
	"github.com/alianjidaniir-design/SamplePRJ/models/repositories"
	"github.com/alianjidaniir-design/SamplePRJ/statics/constants/controllerBaseErrCode"
	"github.com/gofiber/fiber/v2"
)

func Info(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "22")
	defer mainController.FinishAPISpan(ctx)

	req := commonSchema.BaseRequest[userSchema.InfoRequest]{}

	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerBaseErrCode.UserErrCode, "01", errStr, code, err)
	}

	res, errStr, code, err := repositories.UserRepo.Info(spanCtx, req, mainController.GetUser(ctx))
	if err != nil {
		return mainController.Error(ctx, controllerBaseErrCode.UserErrCode, "02", errStr, code, err)
	}

	return mainController.Response(ctx, res)
}
