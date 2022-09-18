# gamesheet: Simple agent character sheet

This package defines a struct that can be used to store and manage a character's status and attributes.

I try to be as conservative as possible with the size of the struct to allow for a lot of instances without filling up the memory.

## Current scope
- Leveling / XP requirement calculations
- Skill point generation
- AP / HP leveling
- AP / HP regeneration
- Status (hunger, thirst, exhaustion, stress)
  - Rudimentary support 
  - Increase over time (on tick)
- Handle death (through exhaustion, injury)

## Planned
- Provide means to reduce status values (hunger, thirst, etc.)
- Handle injury (causes stress and damage)
- Adjust limits of Status values based on resilience, etc.

## TODO

### States

Add a "state", which will influence the growth factor for each status.
This could be a struct that can be defined externally, which would allow for a lot of flexibility wherever this is being used. New states like "swimming" or "flying" could be added easily without the need to modify the gamesheet package.

Instead of making the "rate" part of the Status struct, the rate could be part of the state struct, and instead of individually changing each status's rate, we simply refer to the respective property of the state.

Alternative, state could be a number of states, which are summed up to the rates of each status.

For example: 
```
StatusAwake:  (exhaustion:  0.5, hunger: 0.05,  thirst: 0.1)
StatusAsleep: (exhaustion: -1.0, hunger: 0.025, thirst: 0.05, stress: -0.5)
StatusAfraid: (exhaustion:  0.1, stress: 0.2)
StatusRunning:(exhaustion:  0.5, hunger: 0.01,  thirst: 0.1,  stress: 0.1)
StatusCombat: (exhaustion:  0.5, hunger: 0.1,   thirst: 0.2,  stress: 0.1)
```

So, if the state is []{StatusAwake, StatusAfraid}, the resulting rates would be:
```
exhaustion:  0.5 + 0.1   = 0.6
hunger:      0.05        = 0.05
thirst:      0.1         = 0.1
stress:      0.2         = 0.2
```

If the state is []{StatusAsleep}, the resulting rates would be:
```
exhaustion: -1.0
hunger:      0.025
thirst:      0.05
stress:     -0.5
```

#### State duration.

Some states should only be active for a set duration or have an exit condition. Sleeping, for example, should end when we are fully rested or after a set number of time.

Should this be handled by the AI/Player?

