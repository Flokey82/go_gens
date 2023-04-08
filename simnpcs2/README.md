# Simnpcs2

This package is a skeleton for a rebuilt version of the simnpcs and gamecs package. It is a work in progress, so ignore this for now.


![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/simnpcs2/images/sample.webp "Pixel People!")

## TODO

- [X] Add pathfinding
- [X] Add noise for obstacle generation
- [ ] Add needs and remediaion (food, water, sleep, etc.)
- [ ] Add AI
- [ ] Add combat / interactions
- [ ] Add inventory
- [ ] Exports
    - [X] Unify rendering
    - [X] Add GIF export
        - [ ] Make it optional
    - [X] Add WEBP export
        - [ ] Make it optional
        - [ ] Prevent cgo crashes on exit

## Brainstorming

This section just contains random thoughts and ideas.

### AI

The AI should be able to do the following:

- [ ] Recognize any needs (based on Maslow's hierarchy of needs)
- [ ] Perform simple actions like:
    - [X] Walk to a certain point
    - [X] Perceive the environment
    - [ ] Find items (of specific types)
    - [ ] Pick up items
    - [ ] Consume items
    - [ ] Interact with other NPCs
- [ ] Can chain certain actions together to complete tasks
      (e.g. hungry? -> find food -> eat -> resume previous task)
- [ ] Can be interrupted by other NPCs
- [ ] Tasks can be sorted by priority (and feasibility)
- [ ] Tasks can be interrupted and resumed (or completely restarted?)

### Needs

Needs are a way to determine what an NPC should do. They are based on Maslow's hierarchy of needs. The needs are:

- [ ] Physiological (e.g. food, water, sleep)
- [ ] Safety (e.g. shelter, protection)
- [ ] Belongingness (e.g. friends, family)
- [ ] Esteem (e.g. respect, recognition)
- [ ] Self-actualization (e.g. personal growth, fulfillment)

### Perception

Perception is a way to determine what an NPC can see. It should be able to do the following:

- [X] See other NPCs
- [X] See items
- [ ] See obstacles
- [ ] See the environment (e.g. water, trees, etc.)

### Interactions

Interactions are a way for NPCs to interact with each other. They should be able to do the following:

- [ ] Talk to each other
- [ ] Trade items
- [ ] Fight each other

### Inventory

The inventory is a way for NPCs to store items. It should be able to do the following:

- [ ] Store items
- [ ] Remove items
- [ ] Check if an item is in the inventory

### Items

Items represent objects in the environment that can be moved, picked up, and, depending on the type, consumed. They should be able to do the following:

- [ ] Be picked up
- [ ] Be consumed
- [ ] Be dropped
- [ ] Be moved
- [ ] Be stored in an inventory

### Obstacles

Obstacles are objects in the environment that NPCs cannot move through. They should be able to do the following:

- [X] Be generated based on a noise function / heightmap

### Tasks

Tasks are a way to determine what an NPC should do. They should be able to do the following:

- [ ] Be sorted by priority
- [ ] Be sorted by feasibility
- [ ] Be interrupted and resumed (or completely restarted?)