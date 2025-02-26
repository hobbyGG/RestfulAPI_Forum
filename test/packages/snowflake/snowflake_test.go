package snowflake_test

import (
	"testing"
	"time"

	"github.com/hobbyGG/RestfulAPI_forum/packages/snowflake"
)

func TestGetID(t *testing.T) {
	if err := snowflake.Init(time.Now().Format("2006-01-02"), 1); err != nil {
		t.Error("snowflake.Init error", err)
		return
	}
	get, err := snowflake.GetID()
	if err != nil {
		t.Errorf("snowflake.GetID error: %v\n", err)
	}
	expected := int64(0)
	if get == expected {
		t.Errorf("snowflake.GetID fail want != 0, get %v\n", get)
		return
	}
	t.Log("uid:", get)
}
