package stages

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/paulhenri-l/go-pipeline/mocks"
	"github.com/paulhenri-l/go-pipeline/repo"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewTumblingWindow(t *testing.T) {
	m := newFakeTumblingWindowStage(t)
	w := NewTumblingWindow(m, time.Second, 1)

	assert.IsType(t, &TumblingWindow{}, w)
}

func TestTumblingWindow_Name(t *testing.T) {
	m := newFakeTumblingWindowStage(t)
	w := NewTumblingWindow(m, time.Second, 1)

	assert.Equal(t, "Tumbling window", w.Name())
}

func TestTumblingWindow_Start(t *testing.T) {
	m := newFakeTumblingWindowStage(t)
	w := NewTumblingWindow(m, time.Second, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	items := make(chan interface{})

	out := w.Start(ctx, items)
	close(items)

	_, ok := <-out
	assert.False(t, ok)
}

func TestTumblingWindow_Start_DatapointsSentDownwards(t *testing.T) {
	m := newFakeTumblingWindowStage(t)
	w := NewTumblingWindow(m, time.Second, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	items := make(chan interface{})
	fakeEvent := &someEvent{timestamp: time.Now().Unix()}
	out := w.Start(ctx, items)

	m.EXPECT().Process(gomock.Any(), gomock.Any())

	m.EXPECT().ToDataPoints(gomock.Any(), gomock.Any()).Return([]interface{}{
		"my_first_data_point",
		"my_second_data_point",
	})

	items <- fakeEvent
	close(items)
	assert.Equal(t, "my_first_data_point", <-out)
	assert.Equal(t, "my_second_data_point", <-out)

	_, ok := <-out
	assert.False(t, ok)
}

func TestTumblingWindow_Start_TimestampedArePassedWithWindow(t *testing.T) {
	m := newFakeTumblingWindowStage(t)
	w := NewTumblingWindow(m, time.Second, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	items := make(chan interface{})
	fakeEvent := &someEvent{timestamp: time.Now().Unix()}
	out := w.Start(ctx, items)

	m.EXPECT().Process(gomock.Any(), gomock.Any()).DoAndReturn(func(w *repo.Window, i interface{}) {
		assert.Equal(t, fakeEvent, i.(*someEvent))
		(*w)["hello"] = int64(0)
	})

	m.EXPECT().ToDataPoints(gomock.Any(), gomock.Any()).DoAndReturn(func(tb int64, data *repo.Window) []interface{} {
		assert.Equal(t, int64(0), (*data)["hello"])
		assert.Equal(t, fakeEvent.timestamp, tb)
		return []interface{}{}
	})

	items <- fakeEvent
	close(items)

	_, ok := <-out
	assert.False(t, ok)
}

func TestTumblingWindow_Start_CorrectTimeBin(t *testing.T) {
	m := newFakeTumblingWindowStage(t)
	w := NewTumblingWindow(m, time.Minute, time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	items := make(chan interface{})
	now := time.Now()
	fakeEvent := &someEvent{timestamp: now.Unix()}
	out := w.Start(ctx, items)

	m.EXPECT().Process(gomock.Any(), gomock.Any())

	m.EXPECT().ToDataPoints(gomock.Any(), gomock.Any()).DoAndReturn(func(tb int64, data *repo.Window) []interface{} {
		bin := now.Add(time.Duration(now.Second()) * time.Second * -1).Unix()
		assert.Equal(t, bin, tb)
		return []interface{}{}
	})

	items <- fakeEvent
	close(items)
	<-out
}

func TestTumblingWindow_Start_OldEventsAreDiscarded(t *testing.T) {
	m := newFakeTumblingWindowStage(t)
	w := NewTumblingWindow(m, time.Minute, time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	items := make(chan interface{})
	now := time.Now()
	out := w.Start(ctx, items)

	m.EXPECT().Process(gomock.Any(), gomock.Any()).Times(1)
	m.EXPECT().ToDataPoints(gomock.Any(), gomock.Any()).Return([]interface{}{})

	items <- &someEvent{timestamp: now.Add(time.Minute * -1).Unix()}
	items <- &someEvent{timestamp: now.Add((time.Minute + time.Second) * -1).Unix()}
	close(items)
	<-out
}

func TestTumblingWindow_Start_WindowOutputOnceTooOld(t *testing.T) {
	m := newFakeTumblingWindowStage(t)
	w := NewTumblingWindow(m, time.Second, time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	items := make(chan interface{})
	fakeEvent := &someEvent{timestamp: time.Now().Unix()}
	out := w.Start(ctx, items)

	m.EXPECT().Process(gomock.Any(), gomock.Any())

	m.EXPECT().ToDataPoints(gomock.Any(), gomock.Any()).DoAndReturn(func(tb int64, data *repo.Window) []interface{} {
		return []interface{}{"some_data_point"}
	})

	items <- fakeEvent
	_ = <-out
	receiveTimestamp := time.Now().Unix()

	// Window of one second, with late window wait of 1s = 2 seconds total wait time
	assert.GreaterOrEqual(t, receiveTimestamp, fakeEvent.timestamp+int64(3))
	close(items)
}

func TestTumblingWindow_Start_DatapointsDrainedOnClose(t *testing.T) {
	m := newFakeTumblingWindowStage(t)
	w := NewTumblingWindow(m, time.Hour, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	items := make(chan interface{})
	fakeEvent := &someEvent{timestamp: time.Now().Unix()}
	out := w.Start(ctx, items)

	m.EXPECT().Process(gomock.Any(), gomock.Any())

	m.EXPECT().ToDataPoints(gomock.Any(), gomock.Any()).Return([]interface{}{
		"some_data_point",
	})

	items <- fakeEvent
	close(items)

	// Windows are one hour long but the datapoint has been emitted early
	// because of closure.
	assert.Equal(t, "some_data_point", <-out)

	_, ok := <-out
	assert.False(t, ok)
}

func TestTumblingWindow_Start_StopsWhenInClose(t *testing.T) {
	m := newFakeTumblingWindowStage(t)
	w := NewTumblingWindow(m, time.Hour, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	items := make(chan interface{})
	out := w.Start(ctx, items)

	close(items)

	_, ok := <-out
	assert.False(t, ok)
}

func TestTumblingWindow_Start_StopsWhenCancel(t *testing.T) {
	m := newFakeTumblingWindowStage(t)
	w := NewTumblingWindow(m, time.Hour, 1)
	ctx, cancel := context.WithCancel(context.Background())
	items := make(chan interface{})
	out := w.Start(ctx, items)

	cancel()

	_, ok := <-out
	assert.False(t, ok)
}

func TestTumblingWindow_Start_WrongTypeReported(t *testing.T) {
	m := newFakeTumblingWindowStage(t)
	w := NewTumblingWindow(m, time.Hour, 1)
	ctx, cancel := context.WithCancel(context.Background())
	items := make(chan interface{})
	out := w.Start(ctx, items)

	m.EXPECT().HandleError(gomock.Any()).Do(func(e error) {
		assert.Error(t, e)
		assert.Contains(t, e.Error(), "wrong_type")
	})

	items<-"wrong_type"
	cancel()

	_, ok := <-out
	assert.False(t, ok)
}

func newFakeTumblingWindowStage(t *testing.T) *mocks.MockTumblingWindowStage {
	ctl := gomock.NewController(t)

	return mocks.NewMockTumblingWindowStage(ctl)
}

type someEvent struct {
	timestamp int64
}

func (s *someEvent) GetTimestamp() int64 {
	return s.timestamp
}
