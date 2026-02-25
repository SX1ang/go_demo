package snowflake

import (
	"testing"
)

func TestGenID(t *testing.T) {
	// 初始化雪花算法节点
	Init("2020-01-01", 1)

	// 生成一个 ID
	id1 := GenID()
	if id1 == 0 {
		t.Errorf("expected non-zero ID, got %d", id1)
	}

	// 再生成一个 ID
	id2 := GenID()
	if id2 == id1 {
		t.Errorf("expected unique IDs, got duplicate %d", id1)
	}

	// 批量生成，检查唯一性
	ids := make(map[int64]bool)
	for i := 0; i < 1000; i++ {
		id := GenID()
		if ids[id] {
			t.Errorf("duplicate ID found: %d", id)
		}
		ids[id] = true
	}
}
