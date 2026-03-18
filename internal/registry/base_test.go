package registry

import "testing"

func TestBaseCRUD(t *testing.T) {
	r := NewBase[string, int]()

	if err := r.Register("a", 1); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	if err := r.Register("a", 2); err != ErrDuplicateKey {
		t.Fatalf("expected duplicate key error, got %v", err)
	}

	v, err := r.Get("a")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if v != 1 {
		t.Fatalf("expected value 1, got %d", v)
	}

	r.Upsert("a", 3)
	v, err = r.Get("a")
	if err != nil {
		t.Fatalf("get failed after upsert: %v", err)
	}
	if v != 3 {
		t.Fatalf("expected value 3, got %d", v)
	}

	if err := r.Delete("a"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	if _, err := r.Get("a"); err != ErrNotFound {
		t.Fatalf("expected not found after delete, got %v", err)
	}
}
