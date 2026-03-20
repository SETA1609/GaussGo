package controllers

import (
	"errors"
	"testing"

	"gaussgo/internal/apperrors"
	"gaussgo/internal/contracts"
	"gaussgo/internal/types"
)

func TestAppControllerLoadActiveStateUpdatesContext(t *testing.T) {
	_, store, created := prepareStateAndMods(t)
	recorder := &eventRecorder{}
	ctx := &types.RuntimeContext{}
	c := NewAppController(store, ctx, noopLogger{}, recorder)

	loaded, err := c.LoadActiveState()
	if err != nil {
		t.Fatalf("load active state failed: %v", err)
	}
	if loaded.StateID != created.StateID {
		t.Fatalf("expected loaded state %s, got %s", created.StateID, loaded.StateID)
	}
	if ctx.ActiveStateID != created.StateID {
		t.Fatalf("expected ctx active state %s, got %s", created.StateID, ctx.ActiveStateID)
	}
	if ctx.CurrentLocale != "en" {
		t.Fatalf("expected locale en, got %s", ctx.CurrentLocale)
	}
	if len(recorder.events) != 1 || recorder.events[0].Name != EventStateLoaded {
		t.Fatalf("expected one state.loaded event, got %#v", recorder.events)
	}
}

func TestAppControllerLoadActiveStateUsesPrepopulatedContextStateID(t *testing.T) {
	store := &stubStateStore{
		loadFn: func(stateID string) (contracts.State, error) {
			if stateID != "preset" {
				t.Fatalf("expected load for preset, got %s", stateID)
			}
			return contracts.State{StateID: "preset", UI: contracts.UIState{Locale: "en"}, EnabledMods: []string{"core"}}, nil
		},
		getActiveFn: func() (string, error) {
			t.Fatal("GetActive should not be called when ctx.ActiveStateID is prepopulated")
			return "", nil
		},
	}
	ctx := &types.RuntimeContext{ActiveStateID: "preset"}
	c := NewAppController(store, ctx, noopLogger{}, &eventRecorder{})

	if _, err := c.LoadActiveState(); err != nil {
		t.Fatalf("load active state failed: %v", err)
	}
}

func TestAppControllerLoadActiveStatePropagatesMissingActivePointer(t *testing.T) {
	store := &stubStateStore{
		getActiveFn: func() (string, error) {
			return "", apperrors.New(apperrors.CodeNotFound, apperrors.ErrorTypeInput, "read active pointer")
		},
		loadFn: func(string) (contracts.State, error) {
			return contracts.State{}, errors.New("unexpected load call")
		},
	}
	c := NewAppController(store, &types.RuntimeContext{}, noopLogger{}, &eventRecorder{})

	if _, err := c.LoadActiveState(); err == nil {
		t.Fatal("expected load active state error")
	}
}

type stubStateStore struct {
	createFn    func(profileName string) (contracts.State, error)
	loadFn      func(stateID string) (contracts.State, error)
	saveFn      func(state contracts.State) error
	setActiveFn func(stateID string) error
	getActiveFn func() (string, error)
}

func (s *stubStateStore) Create(profileName string) (contracts.State, error) {
	if s.createFn == nil {
		return contracts.State{}, errors.New("create not implemented")
	}
	return s.createFn(profileName)
}

func (s *stubStateStore) Load(stateID string) (contracts.State, error) {
	if s.loadFn == nil {
		return contracts.State{}, errors.New("load not implemented")
	}
	return s.loadFn(stateID)
}

func (s *stubStateStore) Save(state contracts.State) error {
	if s.saveFn == nil {
		return errors.New("save not implemented")
	}
	return s.saveFn(state)
}

func (s *stubStateStore) SetActive(stateID string) error {
	if s.setActiveFn == nil {
		return errors.New("set active not implemented")
	}
	return s.setActiveFn(stateID)
}

func (s *stubStateStore) GetActive() (string, error) {
	if s.getActiveFn == nil {
		return "", errors.New("get active not implemented")
	}
	return s.getActiveFn()
}
