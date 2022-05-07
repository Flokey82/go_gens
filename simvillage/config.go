package simvillage

const (
	VERSION = 1.1
	VERBOSE = true

	// Tier 0 -- World events
	// Tier 1 -- Character information
	// Tier 2 -- Daily happenings
	// Tier 3 -- Details
	// Tier 4 -- Math details
	// Tier 5 -- Everything
	LOGGING_VERBOSITY = 3

	// How many villagers to start with,
	// cannot select less than 8
	// Default: 10
	STARTING_POP = 4

	// FRIENDLY_CHANCE
	// Controls how likely neutral villagers will respond
	// positively to interactions
	// Default: 0.6
	// Range: 0.0 - 1.0
	FRIENDLY_CHANCE = 0.6

	// How often settlers will arrive
	// at your village per year
	// Default = 12
	SETTLER_CHANCE = 22

	// WORLD_SADNESS
	// This controls how likely it is for people to have a good day
	// Default: 0.7
	// Accepted Range: 0.0 - 1.0
	AVG_HAPPY = 0.7

	// SOCIAL_CHANCE
	// This controls the % of the population that has a social
	// interaction every day. Setting this to 0.0 will disable
	// all social events.
	// Warning: Setting this above 1 might cause performance
	// issues.
	// Default: 0.6
	// Accepted Range: 0.0 <=
	SOCIAL_CHANCE = 0.6

	// DISEASE_ENABLED
	// Enables and disables the disease module
	DISEASE_ENABLED = false

	// DISEASE_CHANCE
	// This controls how frequenty disease will break out
	// on average per year
	// Default: 1
	DISEASE_CHANCE = 1

	// DISEASE_SEVERITY
	// This controls how much of a population will be effected
	// by the sickness
	// Default: 0.2
	// Accepted Range: 0.0 - 1.0
	DISEASE_SEVERITY = 0.1
)
