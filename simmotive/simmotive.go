// Package simmotive is a crude port of https://github.com/alexcu/motive-simulator which
// is adapted from Don Hopkins' article The Soul of The Sims which shows an prototype of
// the 'soul' of what became The Sims 1, written January 23, 1997.
package simmotive

import (
	"fmt"
	"log"
	"math/rand"
	"os"
)

func SRand(upper int) float64 {
	return float64(rand.Intn(upper) + 1)
}

// Motive holds the "mental state" of a Sim.
type Motive struct {
	Motive    [mMax]float64 // current state
	oldMotive [mMax]float64 // previous state
	ClockH    int           // hour
	ClockM    int           // minute
	ClockD    int           // day
	Logs      []string      // recent events
}

// NewMotive returns a new motive.
func NewMotive() *Motive {
	return &Motive{
		ClockH: 8,
	}
}

// Clear the console output.
func (m *Motive) Clear() {
	os.Stdout.Write([]byte{0x1B, 0x5B, 0x33, 0x3B, 0x4A, 0x1B, 0x5B, 0x48, 0x1B, 0x5B, 0x32, 0x4A})
}

// Log adds another log message and truncates the entries if there are more than 12.
func (m *Motive) Log(msg string) {
	m.Logs = append(m.Logs, fmt.Sprintf("[Day %d at %02d:%02d]  %s", m.ClockD, m.ClockH, m.ClockM, msg))
	if len(m.Logs) > 12 {
		m.Logs = m.Logs[1:]
	}
}

// The various motives defining the mental state of a Sim.
const (
	mHappyLife   = 0
	mHappyWeek   = 1
	mHappyDay    = 2
	mHappyNow    = 3
	mPhysical    = 4
	mEnergy      = 5
	mComfort     = 6
	mHunger      = 7
	mHygiene     = 8
	mBladder     = 9
	mMental      = 10
	mAlertness   = 11
	mStress      = 12
	mEnvironment = 13
	mSocial      = 14
	mEntertained = 15
	mMax         = 16
)

// 1 tick = 2 minutes game time
const (
	dayTicks  = 720
	weekTicks = 5040
)

// Init initializes the motive.
func (m *Motive) Init() {
	// Clear motive.
	for i := 0; i < mMax; i++ {
		m.Motive[i] = 0
	}
	// Set initial state.
	m.Motive[mEnergy] = 70
	m.Motive[mAlertness] = 20
	m.Motive[mHunger] = -40
}

// Simulates internal m.motive changes
func (m *Motive) SimMotives() {
	var tem float64

	// Increase game clock (Jamie, remove this) :)
	m.ClockM += 2
	if m.ClockM > 58 {
		m.ClockM = 0
		m.ClockH++
		if m.ClockH > 23 {
			m.ClockH = 0
			m.ClockD++
		}
	}

	// Energy.
	if m.Motive[mEnergy] > 0 {
		if m.Motive[mAlertness] > 0 {
			m.Motive[mEnergy] -= (m.Motive[mAlertness] / 100)
		} else {
			m.Motive[mEnergy] -= (m.Motive[mAlertness] / 100) * ((100 - m.Motive[mEnergy]) / 50)
		}
	} else {
		if m.Motive[mAlertness] > 0 {
			m.Motive[mEnergy] -= (m.Motive[mAlertness] / 100) * ((100 - m.Motive[mEnergy]) / 50)
		} else {
			m.Motive[mEnergy] -= (m.Motive[mAlertness] / 100)
		}
	}

	// I had some food.
	if m.Motive[mHunger] > m.oldMotive[mHunger] {
		tem = m.Motive[mHunger] - m.oldMotive[mHunger]
		m.Motive[mEnergy] += tem / 4
	}

	// Comfort.
	if m.Motive[mBladder] < 0 {
		m.Motive[mComfort] += m.Motive[mBladder] / 10 // max -10
	}
	if m.Motive[mHygiene] < 0 {
		m.Motive[mComfort] += m.Motive[mHygiene] / 20 // max -5
	}
	if m.Motive[mHunger] < 0 {
		m.Motive[mComfort] += m.Motive[mHunger] / 20 // max -5
	}
	// dec a max 100/cycle in a cubed curve (seek zero)
	m.Motive[mComfort] -= (m.Motive[mComfort] * m.Motive[mComfort] * m.Motive[mComfort]) / 10000

	// Hunger.
	tem = ((m.Motive[mAlertness] + 100) / 200) * ((m.Motive[mHunger] + 100) / 100) // ^alert * hunger^0

	if m.Motive[mStress] < 0 { // stress -> hunger
		m.Motive[mHunger] += (m.Motive[mStress] / 100) * ((m.Motive[mHunger] + 100) / 100)
	}
	if m.Motive[mHunger] < -99 {
		m.Log("You have starved to death")
		m.Motive[mHunger] = 80
	}

	// Hygiene.
	if m.Motive[mAlertness] > 0 {
		m.Motive[mHygiene] -= 0.3
	} else {
		m.Motive[mHygiene] -= 0.1
	}
	// Hit hygiene limit, take a bath.
	if m.Motive[mHygiene] < -97 {
		m.Log("You smell very bad, mandatory bath")
		m.Motive[mHygiene] = 80
	}

	// Bladder.
	if m.Motive[mAlertness] > 0 {
		m.Motive[mBladder] -= 0.4 // Bladder fills faster while awake.
	} else {
		m.Motive[mBladder] -= 0.2
	}

	// Food eaten goes into bladder (well, not really,
	// but I don't handle number two separately).
	if m.Motive[mHunger] > m.oldMotive[mHunger] {
		tem = m.Motive[mHunger] - m.oldMotive[mHunger]
		m.Motive[mBladder] -= tem / 4
	}

	// If we hit limit, gotta go.
	if m.Motive[mBladder] < -97 {
		if m.Motive[mAlertness] < 0 {
			m.Log("You have wet your bed")
		} else {
			m.Log("You have soiled the carpet")
		}
		m.Motive[mBladder] = 9
	}

	// Alertness.
	if m.Motive[mAlertness] > 0 {
		// max delta at zero
		tem = (100 - m.Motive[mAlertness]) / 50
	} else {
		tem = (m.Motive[mAlertness] + 100) / 50
	}
	if m.Motive[mEnergy] > 0 {
		if m.Motive[mAlertness] > 0 {
			m.Motive[mAlertness] += (m.Motive[mEnergy] / 100) * tem
		} else {
			m.Motive[mAlertness] += (m.Motive[mEnergy] / 100)
		}
	} else {
		if m.Motive[mAlertness] > 0 {
			m.Motive[mAlertness] += (m.Motive[mEnergy] / 100)
		} else {
			m.Motive[mAlertness] += (m.Motive[mEnergy] / 100) * tem
		}
	}
	m.Motive[mAlertness] += (m.Motive[mEntertained] / 300) * tem

	if m.Motive[mBladder] < -50 {
		m.Motive[mAlertness] -= (m.Motive[mBladder] / 100) * tem
	}

	// Stress
	m.Motive[mStress] += m.Motive[mComfort] / 10     // max -10
	m.Motive[mStress] += m.Motive[mEntertained] / 10 // max -10
	m.Motive[mStress] += m.Motive[mEnvironment] / 15 // max -7
	m.Motive[mStress] += m.Motive[mSocial] / 20      // max -5

	// Cut stress while asleep
	if m.Motive[mAlertness] < 0 {
		m.Motive[mStress] = m.Motive[mStress] / 3
	}
	// dec a max 100/cycle in a cubed curve (seek zero)
	m.Motive[mStress] -= (m.Motive[mStress] * m.Motive[mStress] * m.Motive[mStress]) / 10000

	if m.Motive[mStress] < 0 {
		if (SRand(30) - 100) > m.Motive[mStress] {
			if (SRand(30) - 100) > m.Motive[mStress] {
				m.Log("You have lost your temper")
				m.ChangeMotive(mStress, 20)
			}
		}
	}

	// Environment (TODO)

	// Social (TODO)

	// Entertained.
	// cut entertained while asleep
	if m.Motive[mAlertness] < 0 {
		m.Motive[mEntertained] = m.Motive[mEntertained] / 2
	}
	// Calc physical.
	tem = m.Motive[mEnergy]
	tem += m.Motive[mComfort]
	tem += m.Motive[mHunger]
	tem += m.Motive[mHygiene]
	tem += m.Motive[mBladder]
	tem = tem / 5

	// map the linear average into squared curve
	if tem > 0 {
		tem = 100 - tem
		tem = (tem * tem) / 100
		tem = 100 - tem
	} else {
		tem = 100 + tem
		tem = (tem * tem) / 100
		tem = tem - 100
	}
	m.Motive[mPhysical] = tem

	// Calc mental.
	tem += m.Motive[mStress] * 2 // Stress counts *2
	tem += m.Motive[mEnvironment]
	tem += m.Motive[mSocial]
	tem += m.Motive[mEntertained]
	tem = tem / 5

	// Map the linear average into squared curve.
	if tem > 0 {
		tem = 100 - tem
		tem = (tem * tem) / 100
		tem = 100 - tem
	} else {
		tem = 100 + tem
		tem = (tem * tem) / 100
		tem = tem - 100
	}
	m.Motive[mMental] = tem

	// Calc and average happiness.
	// happy = mental + physical
	m.Motive[mHappyNow] = (m.Motive[mPhysical] + m.Motive[mMental]) / 2
	m.Motive[mHappyDay] = ((m.Motive[mHappyDay] * (dayTicks - 1)) + m.Motive[mHappyNow]) / dayTicks
	m.Motive[mHappyWeek] = ((m.Motive[mHappyDay] * (weekTicks - 1)) + m.Motive[mHappyNow]) / weekTicks
	m.Motive[mHappyLife] = ((m.Motive[mHappyLife] * 9) + m.Motive[mHappyWeek]) / 10

	for i := 0; i < mMax; i++ {
		// Check for over/underflow.
		if m.Motive[i] > 100 {
			m.Motive[i] = 100
		} else if m.Motive[i] < -100 {
			m.Motive[i] = -100
		}
		m.oldMotive[i] = m.Motive[i] // Save set in oldMotives (for delta tests).
	}
}

// ChangeMotive changes the given motive by the given value.
// Use this to change m.motives (checks overflow)
func (m *Motive) ChangeMotive(motive int, value float64) {
	m.Motive[motive] += value

	// Check for over/underflow.
	if m.Motive[motive] > 100 {
		m.Motive[motive] = 100
	} else if m.Motive[motive] < -100 {
		m.Motive[motive] = -100
	}
}

// SimJob simulates an 9 hour workday.
// Use this to change m.motives (checks overflow).
func (m *Motive) SimJob() {
	m.ClockH += 9
	if m.ClockH > 24 {
		m.ClockH -= 24
	}

	m.Motive[mEnergy] = ((m.Motive[mEnergy] + 100) * 0.3) - 100
	m.Motive[mHunger] = -60 + SRand(20)
	m.Motive[mHygiene] = -70 + SRand(30)
	m.Motive[mBladder] = -50 + SRand(50)
	m.Motive[mAlertness] = 10 + SRand(10)
	m.Motive[mStress] = -50 + SRand(50)
}

// PrintMotive prints the current state of the given motive.
func (m *Motive) PrintMotive(motive int) {
	var str string
	str += "["
	for i := -25; i < 25; i++ {
		if m.Motive[motive]/4 > float64(i) {
			str += "="
		} else {
			str += "-"
		}
	}
	log.Printf(str + "]\n")
}

// PrintMotives prints the current state of all motives.
func (m *Motive) PrintMotives() {
	log.Printf("Happiness\n")
	log.Printf("=========\n")
	log.Printf("Life happiness   :")
	m.PrintMotive(mHappyLife)
	log.Printf("Week happiness   :")
	m.PrintMotive(mHappyWeek)
	log.Printf("Today's happiness:")
	m.PrintMotive(mHappyDay)
	log.Printf("Happiness now    :")
	m.PrintMotive(mHappyNow)

	log.Printf("\nBasic Needs\n")
	log.Printf("===========\n")
	log.Printf("Physical         :")
	m.PrintMotive(mPhysical)
	log.Printf("Energy           :")
	m.PrintMotive(mEnergy)
	log.Printf("Comfort          :")
	m.PrintMotive(mComfort)
	log.Printf("Hunger           :")
	m.PrintMotive(mHunger)
	log.Printf("Hygiene          :")
	m.PrintMotive(mHygiene)
	log.Printf("Bladder          :")
	m.PrintMotive(mBladder)

	log.Printf("\nHigher Needs\n")
	log.Printf("============\n")
	log.Printf("Mental           :")
	m.PrintMotive(mMental)
	log.Printf("Alertness        :")
	m.PrintMotive(mAlertness)
	log.Printf("Stress           :")
	m.PrintMotive(mStress)
	log.Printf("Environment      :")
	m.PrintMotive(mEnvironment)
	log.Printf("Social           :")
	m.PrintMotive(mSocial)
	log.Printf("Entertained      :")
	m.PrintMotive(mEntertained)
}
