package servlet

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_BaseController_Options(t *testing.T) {
	t.Run("should return 405", func(t *testing.T) {
		controller := controller{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		context := NewMockContext(ctrl)
		context.EXPECT().String(405, "").Times(1)

		controller.Options(context)
	})
}

func Test_BaseController_Head(t *testing.T) {
	t.Run("should return 405", func(t *testing.T) {
		controller := controller{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		context := NewMockContext(ctrl)
		context.EXPECT().String(405, "").Times(1)

		controller.Head(context)
	})
}

func Test_BaseController_Get(t *testing.T) {
	t.Run("should return 405", func(t *testing.T) {
		controller := controller{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		context := NewMockContext(ctrl)
		context.EXPECT().String(405, "").Times(1)

		controller.Get(context)
	})
}

func Test_BaseController_Post(t *testing.T) {
	t.Run("should return 405", func(t *testing.T) {
		controller := controller{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		context := NewMockContext(ctrl)
		context.EXPECT().String(405, "").Times(1)

		controller.Post(context)
	})
}

func Test_BaseController_Put(t *testing.T) {
	t.Run("should return 405", func(t *testing.T) {
		controller := controller{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		context := NewMockContext(ctrl)
		context.EXPECT().String(405, "").Times(1)

		controller.Put(context)
	})
}

func Test_BaseController_Patch(t *testing.T) {
	t.Run("should return 405", func(t *testing.T) {
		controller := controller{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		context := NewMockContext(ctrl)
		context.EXPECT().String(405, "").Times(1)

		controller.Patch(context)
	})
}

func Test_BaseController_Delete(t *testing.T) {
	t.Run("should return 405", func(t *testing.T) {
		controller := controller{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		context := NewMockContext(ctrl)
		context.EXPECT().String(405, "").Times(1)

		controller.Delete(context)
	})
}
