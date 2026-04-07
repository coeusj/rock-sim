package simulations

import (
	"context"
	"sync"
)

type Simulator interface {
	Start(ctx context.Context, wg *sync.WaitGroup)
	Stop() error
	Update() error
}
