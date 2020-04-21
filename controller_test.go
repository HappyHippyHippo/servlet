package servlet

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_BaseController_Options(t *testing.T) {
	controller := Controller{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	context := NewMockContext(ctrl)
	context.EXPECT().String(405, "").Times(1)

	t.Run("return 405", func(t *testing.T) {
		controller.Options(context)
	})
}

func Test_BaseController_Head(t *testing.T) {
	controller := Controller{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	context := NewMockContext(ctrl)
	context.EXPECT().String(405, "").Times(1)

	t.Run("return 405", func(t *testing.T) {
		controller.Head(context)
	})
}

func Test_BaseController_Get(t *testing.T) {
	controller := Controller{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	context := NewMockContext(ctrl)
	context.EXPECT().String(405, "").Times(1)

	t.Run("return 405", func(t *testing.T) {
		controller.Get(context)
	})
}

func Test_BaseController_Post(t *testing.T) {
	controller := Controller{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	context := NewMockContext(ctrl)
	context.EXPECT().String(405, "").Times(1)

	t.Run("return 405", func(t *testing.T) {
		controller.Post(context)
	})
}

func Test_BaseController_Put(t *testing.T) {
	controller := Controller{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	context := NewMockContext(ctrl)
	context.EXPECT().String(405, "").Times(1)

	t.Run("return 405", func(t *testing.T) {
		controller.Put(context)
	})
}

func Test_BaseController_Patch(t *testing.T) {
	controller := Controller{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	context := NewMockContext(ctrl)
	context.EXPECT().String(405, "").Times(1)

	t.Run("return 405", func(t *testing.T) {
		controller.Patch(context)
	})
}

func Test_BaseController_Delete(t *testing.T) {
	controller := Controller{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	context := NewMockContext(ctrl)
	context.EXPECT().String(405, "").Times(1)

	t.Run("return 405", func(t *testing.T) {
		controller.Delete(context)
	})
}
