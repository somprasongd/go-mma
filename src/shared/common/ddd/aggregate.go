package ddd

// type Aggregate interface {
// 	GetDomainEvents() []*DomainEvent
// 	ClearDomainEvents()
// }

type Aggregate[TId any] struct {
	Entity[TId]
	DomainEvents []*DomainEvent
}

func (a *Aggregate[TId]) AddDomainEvent(dv *DomainEvent) {
	if a.DomainEvents == nil {
		a.DomainEvents = make([]*DomainEvent, 0)
	}
	a.DomainEvents = append(a.DomainEvents, dv)
}

func (a *Aggregate[TId]) GetDomainEvents() []*DomainEvent {
	return a.DomainEvents
}

func (a *Aggregate[TId]) ClearDomainEvents() {
	a.DomainEvents = nil
}
