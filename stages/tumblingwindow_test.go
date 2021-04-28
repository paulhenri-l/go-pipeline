package stages

import (
	"github.com/golang/mock/gomock"
	"github.com/paulhenri-l/go-pipeline/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewTumblingWindow(t *testing.T) {
	m := newFakeTumblingWindowStage(t)
	w := NewTumblingWindow(m, time.Second, 1)

	assert.IsType(t, &TumblingWindow{}, w)
}

func newFakeTumblingWindowStage(t *testing.T) *mocks.MockTumblingWindowStage {
	ctl := gomock.NewController(t)

	return mocks.NewMockTumblingWindowStage(ctl)
}
