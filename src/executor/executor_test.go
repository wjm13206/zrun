package executor

import (
	"testing"
	"zrun/src/types"
)

func TestResolveOrderDepsFirst(t *testing.T) {
	a := &types.Task{Name: "a", Deps: []string{"b"}}
	b := &types.Task{Name: "b"}
	byName := map[string]*types.Task{"a": a, "b": b}
	ordered, err := resolveOrder(byName, []string{"a"})
	if err != nil {
		t.Fatal(err)
	}
	if len(ordered) != 2 || ordered[0].Name != "b" || ordered[1].Name != "a" {
		t.Errorf("依赖应先执行，得到 %+v", ordered)
	}
}

func TestResolveOrderCycle(t *testing.T) {
	a := &types.Task{Name: "a", Deps: []string{"b"}}
	b := &types.Task{Name: "b", Deps: []string{"a"}}
	byName := map[string]*types.Task{"a": a, "b": b}
	if _, err := resolveOrder(byName, []string{"a"}); err == nil {
		t.Error("成环应报错")
	}
}
