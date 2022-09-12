# gamesheet: Simple agent character sheet

This package defines a struct that can be used to store and manage a character's status and attributes.

## Status effects

A character has several attributes that deremine how effective it is in combat, and its overall health. 

These attributes are:
* Exhaustion
* Hunger
* Thirst
* Stress

### Hunger

Hunger is increasing steadily, however less so while sleeping, more so while moving (or physical exertion).
Hunger can be reduced by eating food.

### Thirst

Thirst is increasing steadily, however less so while sleeping, more so while moving (or physical exertion) (with increased speed in hot environments).
Thirst can be reduced by drinking water (or other liquids).

### Stress

Stress increases sharply when taking damage (or emotional stress) and decreases slowly when not taking damage. Stress decreases faster while resting and sleeping.

### Exhaustion

Exhaustion increases when moving (or physical exertion) and decreases when resting or sleeping.

## Skills (TODO)

Skills are used to determine the probability of success for a given action. A character may choose to improve a specific skill up to a limit which is determined by the level of a character.

## Level (TODO)

A character's level determines how many skill points are available to distribute among the character's skills.

