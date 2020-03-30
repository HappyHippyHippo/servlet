package log

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewLogger(t *testing.T) {
	t.Run("create a new logger", func(t *testing.T) {
		action := "Creating a new logger"

		if logger := NewLogger(); logger == nil {
			t.Errorf("%s didn't return a valid reference to a new logger", action)
		}
	})
}

func Test_Logger_Close(t *testing.T) {
	t.Run("should call the close method on all stored streams", func(t *testing.T) {
		action := "Closing the logger"

		id1 := "stream.1"
		id2 := "stream.2"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		stream1 := NewMockStream(ctrl)
		stream1.EXPECT().Close().Times(1)

		stream2 := NewMockStream(ctrl)
		stream2.EXPECT().Close().Times(1)

		logger := NewLogger()
		logger.AddStream(id1, stream1)
		logger.AddStream(id2, stream2)
		logger.Close()

		if logger.HasStream(id1) {
			t.Errorf("%s didn't removed the stored (%s) stream", action, id1)
		}
		if logger.HasStream(id2) {
			t.Errorf("%s didn't removed the stored (%s) stream", action, id2)
		}
	})
}

func Test_Logger_Signal(t *testing.T) {
	t.Run("should propagate to all stored streams", func(t *testing.T) {
		action := "Calling a signal type of message to the logger"

		id1 := "stream.1"
		id2 := "stream.2"
		channel := "channel"
		level := WARNING
		fields := F{"field": "value"}
		message := "message"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		stream1 := NewMockStream(ctrl)
		stream1.EXPECT().Signal(channel, level, fields, message).Return(nil).Times(1)
		stream1.EXPECT().Close().Return(nil).Times(1)

		stream2 := NewMockStream(ctrl)
		stream2.EXPECT().Signal(channel, level, fields, message).Return(nil).Times(1)
		stream2.EXPECT().Close().Return(nil).Times(1)

		logger := NewLogger()
		defer logger.Close()

		logger.AddStream(id1, stream1)
		logger.AddStream(id2, stream2)

		if result := logger.Signal(channel, level, fields, message); result != nil {
			t.Errorf("%s returned the unexpected error : %v", action, result)
		}
	})

	t.Run("should return on the first error", func(t *testing.T) {
		action := "Calling a signal type of message to the logger when the first stream return an error"

		id1 := "stream.1"
		id2 := "stream.2"
		channel := "channel"
		level := WARNING
		fields := F{"field": "value"}
		message := "message"
		err := fmt.Errorf("__dummy_error__")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		stream1 := NewMockStream(ctrl)
		stream1.EXPECT().Signal(channel, level, fields, message).Return(err).AnyTimes()
		stream1.EXPECT().Close().Return(nil).Times(1)

		stream2 := NewMockStream(ctrl)
		stream2.EXPECT().Signal(channel, level, fields, message).Return(nil).AnyTimes()
		stream2.EXPECT().Close().Return(nil).Times(1)

		logger := NewLogger()
		defer logger.Close()

		logger.AddStream(id1, stream1)
		logger.AddStream(id2, stream2)

		if result := logger.Signal(channel, level, fields, message); result != err {
			t.Errorf("%s return (%v), expected (%v)", action, result, err)
		}
	})
}

func Test_Logger_Broadcast(t *testing.T) {
	t.Run("should propagate to all stored streams", func(t *testing.T) {
		action := "Calling a broadcast type of message to the logger"

		id1 := "stream.1"
		id2 := "stream.2"
		level := WARNING
		fields := F{"field": "value"}
		message := "message"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		stream1 := NewMockStream(ctrl)
		stream1.EXPECT().Broadcast(level, fields, message).Return(nil).Times(1)
		stream1.EXPECT().Close().Return(nil).Times(1)

		stream2 := NewMockStream(ctrl)
		stream2.EXPECT().Broadcast(level, fields, message).Return(nil).Times(1)
		stream2.EXPECT().Close().Return(nil).Times(1)

		logger := NewLogger()
		defer logger.Close()

		logger.AddStream(id1, stream1)
		logger.AddStream(id2, stream2)

		if result := logger.Broadcast(level, fields, message); result != nil {
			t.Errorf("%s returned the unexpected error %v", action, result)
		}
	})

	t.Run("should return on the first error", func(t *testing.T) {
		action := "Calling a broadcast type of message to the logger when the first stream return an error"

		id1 := "stream.1"
		id2 := "stream.2"
		level := WARNING
		fields := F{"field": "value"}
		message := "message"
		err := fmt.Errorf("__dummy_error__")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		stream1 := NewMockStream(ctrl)
		stream1.EXPECT().Broadcast(level, fields, message).Return(err).AnyTimes()
		stream1.EXPECT().Close().Return(nil).Times(1)

		stream2 := NewMockStream(ctrl)
		stream2.EXPECT().Broadcast(level, fields, message).Return(nil).AnyTimes()
		stream2.EXPECT().Close().Return(nil).Times(1)

		logger := NewLogger()
		defer logger.Close()

		logger.AddStream(id1, stream1)
		logger.AddStream(id2, stream2)

		if result := logger.Broadcast(level, fields, message); result != err {
			t.Errorf("%s return (%v), expected (%v)", action, result, err)
		}
	})
}

func Test_Logger_HasStream(t *testing.T) {
	t.Run("should correctly check the registration of a stream", func(t *testing.T) {
		action := "Checking if a stream is registed in the logger"

		id1 := "stream.1"
		id2 := "stream.2"
		id3 := "stream.3"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		stream1 := NewMockStream(ctrl)
		stream1.EXPECT().Close().Return(nil).Times(1)

		stream2 := NewMockStream(ctrl)
		stream2.EXPECT().Close().Return(nil).Times(1)

		logger := NewLogger()
		defer logger.Close()

		logger.AddStream(id1, stream1)
		logger.AddStream(id2, stream2)

		if !logger.HasStream(id1) {
			t.Errorf("%s didn't find the expected stream ; %s", action, id1)
		}
		if !logger.HasStream(id2) {
			t.Errorf("%s didn't find the expected stream : %s", action, id2)
		}
		if logger.HasStream(id3) {
			t.Errorf("%s found a unexpected stream : %s", action, id3)
		}
	})
}

func Test_Logger_AddStream(t *testing.T) {
	t.Run("should return a error if registration method recieves a nil reference", func(t *testing.T) {
		action := "Adding a nil stream to the logger"

		expected := "Invalid nil 'stream' argument"

		logger := NewLogger()
		defer logger.Close()

		err := logger.AddStream("id", nil)
		if err == nil {
			t.Errorf("%s returned an unexpected nil error", action)
			return
		}
		if err.Error() != expected {
			t.Errorf("%s returned (%s) error, expected (%s)", action, err.Error(), expected)
		}
	})

	t.Run("should return a error if there is a stream with the same id", func(t *testing.T) {
		action := "Adding a duplicate id to the logger"

		id := "stream"
		expected := fmt.Sprintf("Duplicate id : %s", id)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		stream1 := NewMockStream(ctrl)
		stream1.EXPECT().Close().Return(nil).Times(1)

		stream2 := NewMockStream(ctrl)

		logger := NewLogger()
		defer logger.Close()

		logger.AddStream(id, stream1)

		if err := logger.AddStream(id, stream2); err == nil {
			t.Errorf("%s didn't result in a error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should correctly register a new stream", func(t *testing.T) {
		action := "Adding a stream to the logger"

		id := "stream"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		stream := NewMockStream(ctrl)
		stream.EXPECT().Close().Return(nil).Times(1)

		logger := NewLogger()
		defer logger.Close()

		if err := logger.AddStream(id, stream); err != nil {
			t.Errorf("%s adding resulted in a unexpected error : %v", action, err)
		}
		if result := logger.Stream(id); !reflect.DeepEqual(result, stream) {
			t.Errorf("%s retrieved (%v) when requesting the (%s) stream, expected (%v)", action, result, id, stream)
		}
	})
}

func Test_Logger_RemoveStream(t *testing.T) {
	t.Run("should correctly unregister a stream", func(t *testing.T) {
		action := "Removing a stream from the logger"

		id := "stream"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		stream := NewMockStream(ctrl)
		stream.EXPECT().Close().Return(nil).Times(1)

		logger := NewLogger()
		defer logger.Close()

		logger.AddStream(id, stream)
		logger.RemoveStream(id)

		if logger.HasStream(id) {
			t.Errorf("%s dnd't removed the requested stream", action)
		}
	})
}

func Test_Logger_Stream(t *testing.T) {
	t.Run("should return nil when trying to retrieve a non-existing stream", func(t *testing.T) {
		action := "Retrieving a non-existing stream from the logger"

		logger := NewLogger()
		defer logger.Close()

		if result := logger.Stream("invalid id"); result != nil {
			t.Errorf("%s returned a non nil value when retrieving a invalid id : %v", action, result)
		}
	})

	t.Run("should correctly retrieve a registed stream", func(t *testing.T) {
		action := "Retrieving a stored stream from the logger"

		id := "stream"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		stream := NewMockStream(ctrl)
		stream.EXPECT().Close().Return(nil).Times(1)

		logger := NewLogger()
		defer logger.Close()

		logger.AddStream(id, stream)

		if result := logger.Stream(id); !reflect.DeepEqual(result, stream) {
			t.Errorf("%s retrieved (%v) when requesting the (%v) stream, expected (%v)", action, result, id, stream)
		}
	})
}
