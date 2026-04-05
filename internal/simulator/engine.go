package simulator

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// Updater defines the interface for updating telemetry values
type Updater[V any] interface {
	Update(current V) (V, error)
}

// Engine manages the simulation lifecycle
type Engine[K comparable, V any] struct {
	updater Updater[V]
	writer  *Writer[K, V]
	logger  *slog.Logger
}

// NewEngine creates a new simulation engine
func NewEngine[K comparable, V any](u Updater[V], w *Writer[K, V], logger *slog.Logger) *Engine[K, V] {
	if logger == nil {
		logger = slog.Default()
	}

	return &Engine[K, V]{
		updater: u,
		writer:  w,
		logger:  logger,
	}
}

// StartSimulation runs the simulation loop asynchronously with proper error handling and context cancellation
func (e *Engine[K, V]) StartSimulation(ctx context.Context, wg *sync.WaitGroup, key K, initialValue V) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		currentValue := initialValue
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for i := 0; i < 30; i++ {
			select {
			case <-ctx.Done():
				e.logger.Info("simulation cancelled", "reason", ctx.Err(), "iterations_completed", i)
				return
			case <-ticker.C:
				// Update simulation state
				updatedValue, err := e.updater.Update(currentValue)
				if err != nil {
					e.logger.Error("failed to update simulation", "error", err, "iteration", i)
					continue
				}
				currentValue = updatedValue

				// Send telemetry
				if err := e.writer.Write(ctx, key, currentValue); err != nil {
					e.logger.Error("failed to write telemetry", "error", err, "iteration", i)
					continue
				}

				e.logger.Debug("simulation step completed", "iteration", i, "value", currentValue)
			}
		}

		e.logger.Info("simulation completed successfully", "total_iterations", 30)
	}()
}
