package gamecs

import (
	"fmt"
	"log"

	"github.com/Flokey82/aitree"
	"github.com/Flokey82/go_gens/vectors"
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
	log.Println(fmt.Sprintf("%d: ActionWander", l.ai.id))
	if l.EndCondition() {
		return aitree.StateSuccess
	}

	// Check if we are already following a path.
	// If not, we select a random point within 128 meters.
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
	log.Println(fmt.Sprintf("%d: ActionMoveTo", l.ai.id))
	if l.FailFunc() {
		return aitree.StateFailure
	}
	// TODO: Only set target if none is set, or set the target in another action.
	if l.ai.Target != l.PosFunc() {
		l.ai.SetTarget(l.PosFunc())
	}

	// Check if there is still active pathfinding going on.
	// If not, we have arrived and return success.
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
	log.Println(fmt.Sprintf("%d: ActionPickUpItem", l.ai.id))

	// Select an item.
	it := l.ItemFunc()

	// Check if we can see the item.
	// If we can't, it is either too far away, someone snatched it up.
	if !l.ai.CanSee(it) {
		return aitree.StateFailure // We must be too far away
	}

	// Get our agent from the entity manager...
	a := l.ai.w.mgr.GetEntityFromID(l.ai.id)

	// ... and try to add the object to our inventory.
	// If we can't, return the failure state.
	if !a.CInventory.Add(it) {
		return aitree.StateFailure
	}

	log.Println(fmt.Sprintf("%d: picked up %.2f, %.2f", l.ai.id, l.ai.Target.X, l.ai.Target.Y))
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
	log.Println(fmt.Sprintf("%d: ActionConsumeItem", l.ai.id))

	// Select an item.
	it := l.ItemFunc()
	if it == nil {
		return aitree.StateFailure
	}

	// Remove the item from the world.
	l.ai.w.mgr.RemoveItem(it)

	// Signal to our agent that we ate.
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
	log.Println(fmt.Sprintf("%d: ActionTransferItems", l.ai.id))

	// Get the destination inventory.
	dstInv := l.TargetFunc()
	if dstInv == nil {
		return aitree.StateFailure
	}

	// Get our agent from the entity manager...
	a := l.ai.w.mgr.GetEntityFromID(l.ai.id)

	// ... and attempt to transfer all items to the dst inventory.
	// If we can't, return the failure state.
	if !a.CInventory.TransferAll(dstInv) {
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
	log.Println(fmt.Sprintf("%d: ActionIsTrue", l.ai.id))

	// Run the eval function and return the state accordingly.
	if l.Eval() {
		return aitree.StateSuccess
	}
	return aitree.StateFailure
}

type ActionAttack struct {
	ai         *CAi
	TargetFunc func() *Agent
}

func newActionAttack(ai *CAi, f func() *Agent) *ActionAttack {
	return &ActionAttack{
		ai:         ai,
		TargetFunc: f,
	}
}

func (l *ActionAttack) Tick() aitree.State {
	log.Println(fmt.Sprintf("%d: ActionAttack", l.ai.id))

	// Attempt to get the target.
	tgt := l.TargetFunc()
	if tgt == nil {
		return aitree.StateFailure
	}

	// The target is already dead, so we succeed.
	if tgt.Dead() {
		return aitree.StateSuccess
	}

	// Check if the target is in range, and if so, attempt to injure it.
	// TODO: Check if we have enough action points.
	if vectors.Dist2(tgt.Pos, l.ai.w.mgr.GetEntityFromID(l.ai.id).Pos) < 0.2 {
		tgt.Injure(10, l.ai.id)
		log.Println(fmt.Sprintf("%d: Hit %d for 10 damage (%f health remaining)", l.ai.id, tgt.id, tgt.Health()))
		return aitree.StateRunning
	}

	// The target was not in range, so we fail this action.
	return aitree.StateFailure
}
