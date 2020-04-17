package log

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewLogger(t *testing.T) {
	t.Run("create a new logger", func(t *testing.T) {
		if logger := NewLogger(); logger == nil {
			t.Errorf("didn't return a valid reference")
		}
	})
}

func Test_Logger_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := NewLogger()

	id1 := "stream.1"
	stream1 := NewMockStream(ctrl)
	stream1.EXPECT().Close().Times(1)
	logger.AddStream(id1, stream1)

	id2 := "stream.2"
	stream2 := NewMockStream(ctrl)
	stream2.EXPECT().Close().Times(1)
	logger.AddStream(id2, stream2)

	t.Run("execute close process", func(t *testing.T) {
		logger.Close()

		if logger.HasStream(id1) {
			t.Errorf("didn't removed the stream")
		}
		if logger.HasStream(id2) {
			t.Errorf("didn't removed the stream")
		}
	})
}

func Test_Logger_Signal(t *testing.T) {
	id1 := "stream.1"
	id2 := "stream.2"
	channel := "channel"
	level := WARNING
	fields := F{"field": "value"}
	message := "message"

	t.Run("propagate to all streams", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLogger()
		defer logger.Close()

		stream1 := NewMockStream(ctrl)
		stream1.EXPECT().Signal(channel, level, message, fields).Return(nil).Times(1)
		stream1.EXPECT().Close().Return(nil).Times(1)
		logger.AddStream(id1, stream1)

		stream2 := NewMockStream(ctrl)
		stream2.EXPECT().Signal(channel, level, message, fields).Return(nil).Times(1)
		stream2.EXPECT().Close().Return(nil).Times(1)
		logger.AddStream(id2, stream2)

		if err := logger.Signal(channel, level, message, fields); err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("return on the first error", func(t *testing.T) {
		expectedError := "error"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLogger()
		defer logger.Close()

		stream1 := NewMockStream(ctrl)
		stream1.EXPECT().Signal(channel, level, message, fields).Return(fmt.Errorf(expectedError)).AnyTimes()
		stream1.EXPECT().Close().Return(nil).Times(1)
		logger.AddStream(id1, stream1)

		stream2 := NewMockStream(ctrl)
		stream2.EXPECT().Signal(channel, level, message, fields).Return(nil).AnyTimes()
		stream2.EXPECT().Close().Return(nil).Times(1)
		logger.AddStream(id2, stream2)

		if err := logger.Signal(channel, level, message, fields); err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})
}

func Test_Logger_Broadcast(t *testing.T) {
	id1 := "stream.1"
	id2 := "stream.2"
	level := WARNING
	fields := F{"field": "value"}
	message := "message"

	t.Run("propagate to all streams", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLogger()
		defer logger.Close()

		stream1 := NewMockStream(ctrl)
		stream1.EXPECT().Broadcast(level, message, fields).Return(nil).Times(1)
		stream1.EXPECT().Close().Return(nil).Times(1)
		logger.AddStream(id1, stream1)

		stream2 := NewMockStream(ctrl)
		stream2.EXPECT().Broadcast(level, message, fields).Return(nil).Times(1)
		stream2.EXPECT().Close().Return(nil).Times(1)
		logger.AddStream(id2, stream2)

		if err := logger.Broadcast(level, message, fields); err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("return on the first error", func(t *testing.T) {
		expectedError := "error"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLogger()
		defer logger.Close()

		stream1 := NewMockStream(ctrl)
		stream1.EXPECT().Broadcast(level, message, fields).Return(fmt.Errorf(expectedError)).AnyTimes()
		stream1.EXPECT().Close().Return(nil).Times(1)
		logger.AddStream(id1, stream1)

		stream2 := NewMockStream(ctrl)
		stream2.EXPECT().Broadcast(level, message, fields).Return(nil).AnyTimes()
		stream2.EXPECT().Close().Return(nil).Times(1)
		logger.AddStream(id2, stream2)

		if err := logger.Broadcast(level, message, fields); err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})
}

func Test_Logger_HasStream(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := NewLogger()
	defer logger.Close()

	id1 := "stream.1"
	stream1 := NewMockStream(ctrl)
	stream1.EXPECT().Close().Return(nil).Times(1)
	logger.AddStream(id1, stream1)

	id2 := "stream.2"
	stream2 := NewMockStream(ctrl)
	stream2.EXPECT().Close().Return(nil).Times(1)
	logger.AddStream(id2, stream2)

	id3 := "stream.3"

	t.Run("check the registration of a stream", func(t *testing.T) {
		if !logger.HasStream(id1) {
			t.Errorf("returned false")
		}
		if !logger.HasStream(id2) {
			t.Errorf("returned false")
		}
		if logger.HasStream(id3) {
			t.Errorf("returned true")
		}
	})
}

func Test_Logger_AddStream(t *testing.T) {
	t.Run("error if nil stream", func(t *testing.T) {
		logger := NewLogger()
		defer logger.Close()

		if err := logger.AddStream("id", nil); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'stream' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error if id is duplicate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLogger()
		defer logger.Close()

		id := "stream"
		stream1 := NewMockStream(ctrl)
		stream1.EXPECT().Close().Return(nil).Times(1)

		stream2 := NewMockStream(ctrl)
		logger.AddStream(id, stream1)

		if err := logger.AddStream(id, stream2); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Duplicate id : stream" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register a new stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLogger()
		defer logger.Close()

		id := "stream"
		stream := NewMockStream(ctrl)
		stream.EXPECT().Close().Return(nil).Times(1)

		if err := logger.AddStream(id, stream); err != nil {
			t.Errorf("resulted the (%v) error", err)
		} else if check := logger.Stream(id); !reflect.DeepEqual(check, stream) {
			t.Errorf("didn0t stored the stream")
		}
	})
}

func Test_Logger_RemoveStream(t *testing.T) {
	t.Run("unregister a stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLogger()
		defer logger.Close()

		id := "stream"
		stream := NewMockStream(ctrl)
		stream.EXPECT().Close().Return(nil).Times(1)

		logger.AddStream(id, stream)
		logger.RemoveStream(id)

		if logger.HasStream(id) {
			t.Errorf("dnd't removed the stream")
		}
	})
}

func Test_Logger_Stream(t *testing.T) {
	t.Run("nil on a non-existing stream", func(t *testing.T) {
		logger := NewLogger()
		defer logger.Close()

		if result := logger.Stream("invalid id"); result != nil {
			t.Errorf("returned a valid stream")
		}
	})

	t.Run("retrieve the requested stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLogger()
		defer logger.Close()

		id := "stream"
		stream := NewMockStream(ctrl)
		stream.EXPECT().Close().Return(nil).Times(1)
		logger.AddStream(id, stream)

		if check := logger.Stream(id); !reflect.DeepEqual(check, stream) {
			t.Errorf("didn0t retrieved the stored stream")
		}
	})
}
