package controllers

import (
	"sync"

	"gaussgo/internal/apperrors"
	"gaussgo/internal/contracts"
	"gaussgo/internal/types"
)

var _ contracts.SceneController = (*SceneController)(nil)

type SceneController struct {
	mu      sync.RWMutex
	history []string
	current string
	ctx     *types.RuntimeContext
}

func NewSceneController(initialScene string) *SceneController {
	return &SceneController{current: initialScene}
}

func NewSceneControllerWithContext(initialScene string, ctx *types.RuntimeContext) *SceneController {
	controller := &SceneController{current: initialScene, ctx: ctx}
	controller.syncContext()
	return controller
}

func (c *SceneController) Current() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.current
}

func (c *SceneController) Navigate(sceneID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.current != "" {
		c.history = append(c.history, c.current)
	}
	c.current = sceneID
	c.syncContext()
}

func (c *SceneController) Back() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.history) == 0 {
		return apperrors.New(apperrors.CodeNotFound, apperrors.ErrorTypeDomain, "scene history is empty")
	}
	idx := len(c.history) - 1
	c.current = c.history[idx]
	c.history = c.history[:idx]
	c.syncContext()
	return nil
}

func (c *SceneController) syncContext() {
	if c.ctx != nil {
		c.ctx.CurrentSceneID = c.current
	}
}
