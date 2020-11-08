package servlet

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func Test_NewLog(t *testing.T) {
	t.Run("new logger", func(t *testing.T) {
		if logger := NewLog(); logger == nil {
			t.Error("didn't returned a valid reference")
		}
	})
}

func Test_Log_Close(t *testing.T) {
	t.Run("nil pointer receiver", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("didn't panic")
			} else {
				switch e := r.(type) {
				case error:
					if e.Error() != "nil pointer receiver" {
						t.Errorf("panic with the (%v) error", e)
					}
				default:
					t.Error("didn't panic with an error")
				}
			}
		}()

		var log *Log
		_ = log.Close()
	})

	t.Run("execute close process", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()

		id1 := "stream.1"
		stream1 := NewMockLogStream(ctrl)
		stream1.EXPECT().Close().Times(1)
		_ = logger.AddStream(id1, stream1)

		id2 := "stream.2"
		stream2 := NewMockLogStream(ctrl)
		stream2.EXPECT().Close().Times(1)
		_ = logger.AddStream(id2, stream2)

		_ = logger.Close()

		if logger.HasStream(id1) {
			t.Error("didn't removed the stream")
		}
		if logger.HasStream(id2) {
			t.Error("didn't removed the stream")
		}
	})
}

func Test_Log_Signal(t *testing.T) {
	t.Run("propagate to all streams", func(t *testing.T) {
		id1 := "stream.1"
		id2 := "stream.2"
		channel := "channel"
		level := WARNING
		fields := map[string]interface{}{"field": "value"}
		message := "message"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		defer func() { _ = logger.Close() }()

		stream1 := NewMockLogStream(ctrl)
		stream1.EXPECT().Signal(channel, level, message, fields).Return(nil).Times(1)
		stream1.EXPECT().Close().Return(nil).Times(1)
		_ = logger.AddStream(id1, stream1)

		stream2 := NewMockLogStream(ctrl)
		stream2.EXPECT().Signal(channel, level, message, fields).Return(nil).Times(1)
		stream2.EXPECT().Close().Return(nil).Times(1)
		_ = logger.AddStream(id2, stream2)

		if err := logger.Signal(channel, level, message, fields); err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("return on the first error", func(t *testing.T) {
		id1 := "stream.1"
		id2 := "stream.2"
		channel := "channel"
		level := WARNING
		fields := map[string]interface{}{"field": "value"}
		message := "message"
		expectedError := "error"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		defer func() { _ = logger.Close() }()

		stream1 := NewMockLogStream(ctrl)
		stream1.EXPECT().Signal(channel, level, message, fields).Return(fmt.Errorf(expectedError)).AnyTimes()
		stream1.EXPECT().Close().Return(nil).Times(1)
		_ = logger.AddStream(id1, stream1)

		stream2 := NewMockLogStream(ctrl)
		stream2.EXPECT().Signal(channel, level, message, fields).Return(nil).AnyTimes()
		stream2.EXPECT().Close().Return(nil).Times(1)
		_ = logger.AddStream(id2, stream2)

		if err := logger.Signal(channel, level, message, fields); err == nil {
			t.Error("didn't returned the expected  error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})
}

func Test_Log_Broadcast(t *testing.T) {
	t.Run("propagate to all streams", func(t *testing.T) {
		id1 := "stream.1"
		id2 := "stream.2"
		level := WARNING
		fields := map[string]interface{}{"field": "value"}
		message := "message"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		defer func() { _ = logger.Close() }()

		stream1 := NewMockLogStream(ctrl)
		stream1.EXPECT().Broadcast(level, message, fields).Return(nil).Times(1)
		stream1.EXPECT().Close().Return(nil).Times(1)
		_ = logger.AddStream(id1, stream1)

		stream2 := NewMockLogStream(ctrl)
		stream2.EXPECT().Broadcast(level, message, fields).Return(nil).Times(1)
		stream2.EXPECT().Close().Return(nil).Times(1)
		_ = logger.AddStream(id2, stream2)

		if err := logger.Broadcast(level, message, fields); err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("return on the first error", func(t *testing.T) {
		id1 := "stream.1"
		id2 := "stream.2"
		level := WARNING
		fields := map[string]interface{}{"field": "value"}
		message := "message"
		expectedError := "error"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		defer func() { _ = logger.Close() }()

		stream1 := NewMockLogStream(ctrl)
		stream1.EXPECT().Broadcast(level, message, fields).Return(fmt.Errorf(expectedError)).AnyTimes()
		stream1.EXPECT().Close().Return(nil).Times(1)
		_ = logger.AddStream(id1, stream1)

		stream2 := NewMockLogStream(ctrl)
		stream2.EXPECT().Broadcast(level, message, fields).Return(nil).AnyTimes()
		stream2.EXPECT().Close().Return(nil).Times(1)
		_ = logger.AddStream(id2, stream2)

		if err := logger.Broadcast(level, message, fields); err == nil {
			t.Error("didn't returned the expected  error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})
}

func Test_Log_HasStream(t *testing.T) {
	t.Run("check the registration of a stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		defer func() { _ = logger.Close() }()

		id1 := "stream.1"
		stream1 := NewMockLogStream(ctrl)
		stream1.EXPECT().Close().Return(nil).Times(1)
		_ = logger.AddStream(id1, stream1)

		id2 := "stream.2"
		stream2 := NewMockLogStream(ctrl)
		stream2.EXPECT().Close().Return(nil).Times(1)
		_ = logger.AddStream(id2, stream2)

		id3 := "stream.3"

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

func Test_Log_AddStream(t *testing.T) {
	t.Run("nil pointer receiver", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("didn't panic")
			} else {
				switch e := r.(type) {
				case error:
					if e.Error() != "nil pointer receiver" {
						t.Errorf("panic with the (%v) error", e)
					}
				default:
					t.Error("didn't panic with an error")
				}
			}
		}()

		var log *Log
		_ = log.AddStream("id", nil)
	})

	t.Run("error if nil stream", func(t *testing.T) {
		logger := NewLog()
		defer func() { _ = logger.Close() }()

		if err := logger.AddStream("id", nil); err == nil {
			t.Errorf("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'stream' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error if id is duplicate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		defer func() { _ = logger.Close() }()

		id := "stream"
		stream1 := NewMockLogStream(ctrl)
		stream1.EXPECT().Close().Return(nil).Times(1)

		stream2 := NewMockLogStream(ctrl)
		_ = logger.AddStream(id, stream1)

		if err := logger.AddStream(id, stream2); err == nil {
			t.Errorf("didn't returned the expected error")
		} else if err.Error() != "duplicate id : stream" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register a new stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		defer func() { _ = logger.Close() }()

		id := "stream"
		stream := NewMockLogStream(ctrl)
		stream.EXPECT().Close().Return(nil).Times(1)

		if err := logger.AddStream(id, stream); err != nil {
			t.Errorf("resulted the (%v) error", err)
		} else if check := logger.Stream(id); !reflect.DeepEqual(check, stream) {
			t.Errorf("didn't stored the stream")
		}
	})
}

func Test_Log_RemoveStream(t *testing.T) {
	t.Run("nil pointer receiver", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("didn't panic")
			} else {
				switch e := r.(type) {
				case error:
					if e.Error() != "nil pointer receiver" {
						t.Errorf("panic with the (%v) error", e)
					}
				default:
					t.Error("didn't panic with an error")
				}
			}
		}()

		var log *Log
		log.RemoveStream("id")
	})

	t.Run("unregister a stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		defer func() { _ = logger.Close() }()

		id := "stream"
		stream := NewMockLogStream(ctrl)
		stream.EXPECT().Close().Return(nil).Times(1)

		_ = logger.AddStream(id, stream)
		logger.RemoveStream(id)

		if logger.HasStream(id) {
			t.Errorf("dnd't removed the stream")
		}
	})
}

func Test_Log_Stream(t *testing.T) {
	t.Run("nil on a non-existing stream", func(t *testing.T) {
		logger := NewLog()
		defer func() { _ = logger.Close() }()

		if result := logger.Stream("invalid id"); result != nil {
			t.Errorf("returned a valid stream")
		}
	})

	t.Run("retrieve the requested stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		defer func() { _ = logger.Close() }()

		id := "stream"
		stream := NewMockLogStream(ctrl)
		stream.EXPECT().Close().Return(nil).Times(1)
		_ = logger.AddStream(id, stream)

		if check := logger.Stream(id); !reflect.DeepEqual(check, stream) {
			t.Errorf("didn0t retrieved the stored stream")
		}
	})
}
