package gamecs

import (
	"fmt"
	"github.com/Flokey82/go_gens/aitree"
	"github.com/Flokey82/go_gens/vectors"
	"log"
)

type ActionWander struct {
	ai           *CAi
	EndCondition func() bool
}

func newActionWander(ai *CAi, f func() bool) *ActionWander {
	return &ActionWander{
		ai:           ai,
		EndCondition: f,
	}
}

func (l *ActionWander) Tick() aitree.State {
	log.Println("ActionWander")
	if l.EndCondition() {
		return aitree.StateSuccess
	}

	// There is nothing to eat.
	// Select a random point within 128 meters.
	if !l.ai.CAiPath.active {
		l.ai.SetTarget(vectors.RandomVec2(128.0))
	}

	return aitree.StateRunning
}

type ActionMoveTo struct {
	ai       *CAi
	FailFunc func() bool
	PosFunc  func() vectors.Vec2
}

func newActionMoveTo(ai *CAi, ff func() bool, f func() vectors.Vec2) *ActionMoveTo {
	return &ActionMoveTo{
		ai:       ai,
		FailFunc: ff,
		PosFunc:  f,
	}
}

func (l *ActionMoveTo) Tick() aitree.State {
	log.Println("ActionMoveTo")
	if l.FailFunc() {
		return aitree.StateFailure
	}
	// TODO: Only set target if none is set, or set the target in another action.
	if l.ai.Target != l.PosFunc() {
		l.ai.SetTarget(l.PosFunc())
	}

	// There is nothing to eat.
	// Select a random point within 128 meters.
	if !l.ai.CAiPath.active {
		return aitree.StateSuccess
	}

	return aitree.StateRunning
}

type ActionPickUpItem struct {
	ai       *CAi
	ItemFunc func() *Item
}

func newActionPickUpItem(ai *CAi, f func() *Item) *ActionPickUpItem {
	return &ActionPickUpItem{
		ai:       ai,
		ItemFunc: f,
	}
}

func (l *ActionPickUpItem) Tick() aitree.State {
	log.Println("ActionPickUpItem")

	it := l.ItemFunc()
	if !l.ai.CanSee(it) {
		return aitree.StateFailure // We must be too far away
	}
	a := l.ai.w.mgr.GetEntityFromID(l.ai.id)
	if !a.CInventory.Add(it) {
		return aitree.StateFailure
	}
	// TODO: Message that we're munching, so we'd need to reset hunger.
	log.Println(fmt.Sprintf("picked up %.2f, %.2f", l.ai.Target.X, l.ai.Target.Y))
	return aitree.StateSuccess
}

type ActionConsumeItem struct {
	ai       *CAi
	ItemFunc func() *Item
}

func newActionConsumeItem(ai *CAi, f func() *Item) *ActionConsumeItem {
	return &ActionConsumeItem{
		ai:       ai,
		ItemFunc: f,
	}
}

func (l *ActionConsumeItem) Tick() aitree.State {
	log.Println("ActionConsumeItem")

	it := l.ItemFunc()
	if it == nil {
		return aitree.StateFailure
	}
	l.ai.w.mgr.RemoveItem(it)
	l.ai.CAiStatus.Eat()
	// TODO: Message that we're munching, so we'd need to reset hunger.
	log.Println(fmt.Sprintf("ate %.2f, %.2f", l.ai.Target.X, l.ai.Target.Y))

	return aitree.StateSuccess
}

type ActionTransferItems struct {
	ai         *CAi
	TargetFunc func() *CInventory
}

func newActionTransferItems(ai *CAi, f func() *CInventory) *ActionTransferItems {
	return &ActionTransferItems{
		ai:         ai,
		TargetFunc: f,
	}
}

func (l *ActionTransferItems) Tick() aitree.State {
	log.Println("ActionTransferItems")

	it := l.TargetFunc()
	if it == nil {
		return aitree.StateFailure
	}
	if !l.ai.w.mgr.GetEntityFromID(l.ai.id).CInventory.TransferAll(it) {
		return aitree.StateFailure
	}

	return aitree.StateSuccess
}

type ActionIsTrue struct {
	ai   *CAi
	Eval func() bool
}

func newActionIsTrue(ai *CAi, ef func() bool) *ActionIsTrue {
	return &ActionIsTrue{
		ai:   ai,
		Eval: ef,
	}
}

func (l *ActionIsTrue) Tick() aitree.State {
	log.Println("ActionIsTrue")
	if l.Eval() {
		return aitree.StateSuccess
	}
	return aitree.StateFailure
}
