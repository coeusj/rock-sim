package simulator

// Simulation represents a generic simulation that can update telemetry data
type Simulation[V any] struct {
	Value V
}

// NewSimulation creates a new simulation with initial value
func NewSimulation[V any](initialValue V) *Simulation[V] {
	return &Simulation[V]{Value: initialValue}
}

// Update performs a simulation step and returns the updated value
func (s *Simulation[V]) Update(current V) (V, error) {
	// This is a generic implementation - specific types should embed and override
	return current, nil
}

// NavigationSimulation simulates rocket navigation telemetry
type NavigationSimulation struct {
	*Simulation[Navigation]
}

// NewNavigationSimulation creates a navigation simulation
func NewNavigationSimulation(initial Navigation) *NavigationSimulation {
	return &NavigationSimulation{
		Simulation: NewSimulation(initial),
	}
}

// Update performs navigation simulation step
func (s *NavigationSimulation) Update(current Navigation) (Navigation, error) {
	updated := current
	updated.Altitude += 10.3
	updated.Velocity += 2.5
	s.Value = updated // Update stored value
	return updated, nil
}
