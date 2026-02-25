package task

import (
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/taskSchema"
	mainController "github.com/alianjidaniir-design/SamplePRJ/controllers/maincontroller"
	"github.com/alianjidaniir-design/SamplePRJ/models/repositories"
	"github.com/alianjidaniir-design/SamplePRJ/statics/constants/controllerBaseErrCode"
	"github.com/gofiber/fiber/v2"
)

func Create(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "11")
	defer mainController.FinishAPISpan(ctx)

	req := commonSchema.BaseRequest[taskSchema.CreateRequest]{}

	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerBaseErrCode.TaskErrCode, "01", errStr, code, err)
	}

	res, errStr, code, err := repositories.TaskRepo.Create(spanCtx, req, mainController.GetUser(ctx))
	if err != nil {
		return mainController.Error(ctx, controllerBaseErrCode.TaskErrCode, "02", errStr, code, err)
	}

	return mainController.Response(ctx, res)
}
