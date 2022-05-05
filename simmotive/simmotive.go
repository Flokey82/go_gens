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

func (m *Motive) Clr() {
	os.Stdout.Write([]byte{0x1B, 0x5B, 0x33, 0x3B, 0x4A, 0x1B, 0x5B, 0x48, 0x1B, 0x5B, 0x32, 0x4A})
}

func SRand(upper int) float64 {
	return float64(rand.Intn(upper) + 1)
}

type Motive struct {
	Motive    [16]float64
	oldMotive [16]float64
	ClockH    int
	ClockM    int
	ClockD    int
	Logs      []string
}

func NewMotive() *Motive {
	return &Motive{
		ClockH: 8,
	}
}

func (m *Motive) Log(msg string) {
	m.Logs = append(m.Logs, fmt.Sprintf("[Day %d at %02d:%02d]  %s", m.ClockD, m.ClockH, m.ClockM, msg))
	if len(m.Logs) > 12 {
		m.Logs = m.Logs[1:]
	}
}

const (
	mHappyLife = 0
	mHappyWeek = 1
	mHappyDay  = 2
	mHappyNow  = 3

	mPhysical = 4
	mEnergy   = 5
	mComfort  = 6
	mHunger   = 7
	mHygiene  = 8
	mBladder  = 9

	mMental      = 10
	mAlertness   = 11
	mStress      = 12
	mEnvironment = 13
	mSocial      = 14
	mEntertained = 15
)

// 1 tick = 2 minutes game time
const DAYTICKS = 720
const WEEKTICKS = 5040

func (m *Motive) InitMotives() {
	var count int

	for count = 0; count < 16; count++ {
		m.Motive[count] = 0
	}

	m.Motive[mEnergy] = 70
	m.Motive[mAlertness] = 20
	m.Motive[mHunger] = -40
}

// Simulates internal m.motive changes
func (m *Motive) SimMotives(count int) {
	var tem float64
	//var z int
	// Rect r = { 100, 100, 140, 140 };

	// inc game clock (Jamie, remove this)
	m.ClockM += 2

	if m.ClockM > 58 {
		m.ClockM = 0
		m.ClockH++
		if m.ClockH > 23 {
			m.ClockH = 0
			m.ClockD++
		}
	}

	// energy
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

	// I had some food
	if m.Motive[mHunger] > m.oldMotive[mHunger] {
		tem = m.Motive[mHunger] - m.oldMotive[mHunger]
		m.Motive[mEnergy] += tem / 4
	}

	// comfort
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

	// hunger
	tem = ((m.Motive[mAlertness] + 100) / 200) * ((m.Motive[mHunger] + 100) / 100) // ^alert * hunger^0

	if m.Motive[mStress] < 0 { // stress -> hunger
		m.Motive[mHunger] += (m.Motive[mStress] / 100) * ((m.Motive[mHunger] + 100) / 100)
	}
	if m.Motive[mHunger] < -99 {
		m.Log("You have starved to death")
		m.Motive[mHunger] = 80
	}

	// hygiene
	if m.Motive[mAlertness] > 0 {
		m.Motive[mHygiene] -= 0.3
	} else {
		m.Motive[mHygiene] -= 0.1
	}
	// hit limit, bath
	if m.Motive[mHygiene] < -97 {
		m.Log("You smell very bad, mandatory bath")
		m.Motive[mHygiene] = 80
	}

	// bladder
	if m.Motive[mAlertness] > 0 {
		// bladder fills faster while awake
		m.Motive[mBladder] -= 0.4
	} else {
		m.Motive[mBladder] -= 0.2
	}
	// food eaten goes into bladder
	if m.Motive[mHunger] > m.oldMotive[mHunger] {
		tem = m.Motive[mHunger] - m.oldMotive[mHunger]
		m.Motive[mBladder] -= tem / 4
	}
	// hit limit, gotta go
	if m.Motive[mBladder] < -97 {
		if m.Motive[mAlertness] < 0 {
			m.Log("You have wet your bed")
		} else {
			m.Log("You have soiled the carpet")
		}
		m.Motive[mBladder] = 9
	}

	// alertness
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
	// stress
	m.Motive[mStress] += m.Motive[mComfort] / 10     // max -10
	m.Motive[mStress] += m.Motive[mEntertained] / 10 // max -10
	m.Motive[mStress] += m.Motive[mEnvironment] / 15 // max -7
	m.Motive[mStress] += m.Motive[mSocial] / 20      // max -5

	// cut stress while asleep
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
	// environment

	// social

	// enterntained
	// cut enterntained while asleep
	if m.Motive[mAlertness] < 0 {
		m.Motive[mEntertained] = m.Motive[mEntertained] / 2
	}
	// calc physical
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

	// calc mental
	tem += m.Motive[mStress] // stress counts *2
	tem += m.Motive[mStress]
	tem += m.Motive[mEnvironment]
	tem += m.Motive[mSocial]
	tem += m.Motive[mEntertained]
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
	m.Motive[mMental] = tem

	// calc and average happiness
	// happy = mental + physical
	m.Motive[mHappyNow] = (m.Motive[mPhysical] + m.Motive[mMental]) / 2
	m.Motive[mHappyDay] = ((m.Motive[mHappyDay] * (DAYTICKS - 1)) + m.Motive[mHappyNow]) / DAYTICKS
	m.Motive[mHappyWeek] = ((m.Motive[mHappyDay] * (WEEKTICKS - 1)) + m.Motive[mHappyNow]) / WEEKTICKS
	m.Motive[mHappyLife] = ((m.Motive[mHappyLife] * 9) + m.Motive[mHappyWeek]) / 10

	for z := 0; z < 16; z++ {
		if m.Motive[z] > 100 {
			m.Motive[z] = 100 // check for over/under flow
		}
		if m.Motive[z] < -100 {
			m.Motive[z] = -100
		}
		m.oldMotive[z] = m.Motive[z] // save set in oldMotives (for delta tests)
	}
}

// use this to change m.motives (checks overflow)
func (m *Motive) ChangeMotive(motive int, value float64) {
	m.Motive[motive] += value
	if m.Motive[motive] > 100 {
		m.Motive[motive] = 100
	}
	if m.Motive[motive] < -100 {
		m.Motive[motive] = -100
	}
}

// use this to change m.motives (checks overflow)
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
