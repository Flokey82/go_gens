# Go Generators
This repository contains various small attempts at procedural generation, simulation, and other things that might be useful for some people.

PLEASE NOTE: These packages are under constant development for maybe a month or two until I've figured out a consistent style etc.

DO NOT USE YET!

## Simulation and Procedural Generation

### gencellular: Simple Cellular Automata in Golang (WIP)
This package currently only implements Conway's Game of Life.

### genvillage: Village Economy Generator
This package attempts to generate a stable, self sustaining economy given some initial buildings that serve as seed. The user can define building types that consume and produce resources, which the simulation will use to add new buildings until the sum of produced resources outweighs the sum of consumed resources. Please note that the output of this process is only meant to be used as input for a more sophisticated generator, only ensuring that besides required import of unavailable resources, the village/settlement can be in theory self-sustaining.

### genfloortxt: Simple 2D Text Floorplan Renderer (WIP)
This package provides a little demo on how to 'render' walls as unicode characters, imitating the cp437 symbol set.

## Game Components

### gameloop: Simple Game Loop
This is a very, very basic game loop... nothing fancy about it.