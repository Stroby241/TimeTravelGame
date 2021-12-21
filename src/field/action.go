package field

import (
	"fmt"
	. "github.com/Stroby241/TimeTravelGame/src/math"
)

const (
	actionStay    = 1
	actionMove    = 2
	actionSupport = 3
)

type Action struct {
	TimePos
	Kind    int
	Support int // For Move Action
}

func NewAction() *Action {
	return &Action{
		Kind: actionStay,
	}
}

func (t *Timeline) SetAction(unit *Unit, pos TimePos) {

	if unit.Action.Kind == actionMove {
		for i := len(t.moveUnits) - 1; i >= 0; i-- {
			if t.moveUnits[i] == unit {
				t.moveUnits = append(t.moveUnits[:i], t.moveUnits[i+1:]...)
			}
		}
	} else if unit.Action.Kind == actionSupport {

		if _, actionUnit := t.GetUnitAtPos(unit.Action.TimePos); actionUnit != nil && actionUnit.FactionId == unit.FactionId {
			actionUnit.Support--
		}

		for _, actionUnit := range t.moveUnits {
			if actionUnit.Action.TilePos == unit.Action.TilePos && actionUnit.FactionId == unit.FactionId {
				actionUnit.Action.Support--
			}
		}

		for i := len(t.supportUnits) - 1; i >= 0; i-- {
			if t.supportUnits[i] == unit {
				t.supportUnits = append(t.supportUnits[:i], t.supportUnits[i+1:]...)
			}
		}
	}

	// If TimePos is the same -> Stay
	if unit.SamePos(pos) {
		unit.Action.Kind = actionStay
		unit.Action.FieldPos = CardPos{}
		unit.Action.TilePos = AxialPos{}
		return
	}

	// If is to an own Unit -> Support
	if _, actionUnit := t.GetUnitAtPos(pos); actionUnit != nil && actionUnit.FactionId == unit.FactionId {
		unit.Action.Kind = actionSupport
		unit.Action.TimePos = pos
		unit.Action.Support = 0

		actionUnit.Support++
		t.supportUnits = append(t.supportUnits, unit)
		return
	}

	// If is to an own Move -> Support
	for _, actionUnit := range t.moveUnits {
		if actionUnit.Action.SamePos(actionUnit.Action.TimePos); actionUnit.FactionId == unit.FactionId {

			unit.Action.Kind = actionSupport
			unit.Action.TimePos = pos
			unit.Action.Support = 0

			actionUnit.Action.Support++
			t.supportUnits = append(t.supportUnits, unit)
			return
		}
	}

	// Else -> Move
	unit.Action.Kind = actionMove
	unit.Action.TimePos = pos
	t.moveUnits = append(t.moveUnits, unit)
}

type targetPos struct {
	TimePos
	moveUnits   []*Unit
	presentUnit *Unit
	winningUnit *Unit
}

func (t *Timeline) SubmitRoundUnits() {

	var targetPositions []*targetPos
	for _, unit := range t.moveUnits {

		// Find all aktive Units
		isAktive := false
		for _, pos := range t.ActiveFields {
			if unit.FieldPos == pos {
				isAktive = true
			}
		}
		if !isAktive {
			continue
		}

		// Check if there is already a target Position
		found := false
		for i, positon := range targetPositions {
			if positon.SamePos(unit.Action.TimePos) {
				targetPositions[i].moveUnits = append(targetPositions[i].moveUnits, unit)
				found = true
			}
		}

		// If not add target Position
		if !found {
			position := &targetPos{
				TimePos:   unit.Action.TimePos,
				moveUnits: []*Unit{unit},
			}

			_, position.presentUnit = t.GetUnitAtPos(position.TimePos)
			if position.presentUnit != nil {
				position.moveUnits = append(position.moveUnits, position.presentUnit)
			}

			targetPositions = append(targetPositions, position)
		}
	}

	changeHappend := true

	for changeHappend {
		changeHappend = false
		for _, position := range targetPositions {

			oldWinningUnit := position.winningUnit
			winningAmount := 0
			for _, unit := range position.moveUnits {

				amount := unit.Action.Support + 1
				if amount > winningAmount {
					position.winningUnit = unit
					winningAmount = amount
				} else if amount == winningAmount {
					position.winningUnit = nil
				}
			}

			if position.presentUnit != nil && position.winningUnit == nil {
				position.winningUnit = position.presentUnit
			}

			if position.winningUnit != oldWinningUnit {
				changeHappend = true
			}
		}
	}

	for _, position := range targetPositions {
		if position.winningUnit != nil {
			position.winningUnit.copyToField()
		}
	}
	fmt.Println("DEbug")
}

/*
type targetPos struct {
	oldFieldPos CardPos
	fieldPos    CardPos
	pos         AxialPos
	moveUnits   []*Unit

	loopUnit *Unit
}

func (t *Timeline) ApplyUnitsActions() {

	var timeTravelPositions []*targetPos
	moveUnit := func(unit *Unit, position *targetPos) {
		for j := len(t.supportUnits) - 1; j >= 0; j-- {
			if t.supportUnits[j].Action.ToFieldPos == unit.FieldPos && t.supportUnits[j].Action.ToPos == unit.TilePos {

				t.SetAction(t.supportUnits[j], t.supportUnits[j].FieldPos, t.supportUnits[j].TilePos)

			} else if t.supportUnits[j].Action.ToFieldPos == unit.Action.ToFieldPos && t.supportUnits[j].Action.ToPos == unit.Action.ToPos {

				t.SetAction(t.supportUnits[j], t.supportUnits[j].FieldPos, t.supportUnits[j].TilePos)
			}
		}

		unit.FieldPos = position.fieldPos
		unit.TilePos = position.pos

		t.SetAction(unit, unit.FieldPos, unit.TilePos)

		if position.oldFieldPos != position.fieldPos {
			timeTravelPositions = append(timeTravelPositions, position)
		}
	}

	for _, positon := range targetPositons {
		notAktive := true
		for _, pos := range t.ActiveFields {
			if pos == positon.oldFieldPos {
				notAktive = false
			}
		}
		if notAktive {
			positon.fieldPos = positon.oldFieldPos.Add(CardPos{Y: t.FieldBounds.Y})
		}
	}

	for len(targetPositons) > 0 {

		madeMove := false

		for i, positon := range targetPositons {
			positon.loopUnit = nil

			var winningUnit *Unit
			winningSupport := 0

			var loopWinningUnit *Unit
			loopWinningSupport := math.MaxInt32

			_, presentUnit := t.GetUnitAtPos(positon.oldFieldPos, positon.pos)
			if presentUnit != nil {
				winningSupport = presentUnit.Support + 1
				loopWinningSupport = presentUnit.Support
			}

			for _, unit := range positon.moveUnits {

				if (unit.Action.Support + 1) > winningSupport {
					winningUnit = unit
					winningSupport = unit.Action.Support + 1
				} else if (unit.Action.Support + 1) == winningSupport {
					winningUnit = nil
				}

				if (unit.Action.Support + 1) > loopWinningSupport {
					loopWinningUnit = unit
					loopWinningSupport = unit.Action.Support + 1
				} else if (unit.Action.Support + 1) == loopWinningSupport {
					loopWinningUnit = nil
				}
			}

			if winningUnit != nil {
				if presentUnit != nil {
					t.RemoveUnitAtPos(positon.fieldPos, positon.pos)
				}

				moveUnit(winningUnit, positon)

				targetPositons = append(targetPositons[:i], targetPositons[i+1:]...)

				madeMove = true
				break

			} else if loopWinningUnit != nil && presentUnit != nil {
				positon.loopUnit = loopWinningUnit
			}
		}

		for _, positon := range targetPositons {
			if positon.loopUnit != nil {

				loop := []*targetPos{positon}

				findingPos := true
				loopDone := false
				for findingPos {
					findingPos = false

					for _, testPosition := range targetPositons {
						if positon.loopUnit != nil && testPosition != positon &&
							loop[len(loop)-1].loopUnit.TilePos == testPosition.pos {

							loop = append(loop, testPosition)
							findingPos = true

							if testPosition.loopUnit.TilePos == loop[0].pos {
								loopDone = true
							}

							break
						}
					}
				}

				if loopDone {
					for _, pos := range loop {
						moveUnit(pos.loopUnit, pos)

						for i, targetPositon := range targetPositons {
							if targetPositon == pos {
								targetPositons = append(targetPositons[:i], targetPositons[i+1:]...)
								break
							}
						}
					}
					madeMove = true
					break
				}
			}
		}

		if !madeMove {
			break
		}
	}

	for _, pos := range timeTravelPositions {
		if t.Fields[pos.fieldPos] != nil {
			continue
		}

		field := t.Fields[pos.oldFieldPos]
		t.CopyField(pos.fieldPos, field)
		t.ActiveFields = append(t.ActiveFields, pos.fieldPos)
	}
}
*/
