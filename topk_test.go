package vector

import (
	"reflect"
	"testing"
)

func TestTopK_OrdersDescending(t *testing.T) {
	topk := NewTopK(3)
	topk.Offer(1, 0.5)
	topk.Offer(2, 0.9)
	topk.Offer(3, 0.2)
	topk.Offer(4, 0.9) // tie for highest score
	topk.Offer(5, 0.7)

	res := topk.Results(nil)
	want := []Match{
		{ID: 2, Score: 0.9},
		{ID: 4, Score: 0.9},
		{ID: 5, Score: 0.7},
	}

	if len(res) != len(want) {
		t.Fatalf("Results len = %d, want %d", len(res), len(want))
	}
	for i := range want {
		if res[i].Score != want[i].Score {
			t.Errorf("At %d: got Score %v, want %v", i, res[i].Score, want[i].Score)
		}
	}
}

func TestTopK_KLargerThanInput(t *testing.T) {
	topk := NewTopK(10)
	topk.Offer(1, 0.3)
	topk.Offer(2, 0.8)

	res := topk.Results(nil)
	if len(res) != 2 {
		t.Fatalf("Results len = %d, want 2", len(res))
	}
	if res[0].ID != 2 || res[1].ID != 1 {
		t.Errorf("Unexpected result order: %v", res)
	}
}

func TestTopK_Reset(t *testing.T) {
	topk := NewTopK(2)
	topk.Offer(1, 0.9)
	topk.Offer(2, 0.8)

	topk.Reset()
	topk.Offer(3, 0.4)

	res := topk.Results(nil)
	want := []Match{
		{ID: 3, Score: 0.4},
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("Results after reset = %v, want %v", res, want)
	}
}

func BenchmarkTopK_Offer(b *testing.B) {
	topk := NewTopK(10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		topk.Offer(int32(i), float32(i%100)*0.01)
	}
}
