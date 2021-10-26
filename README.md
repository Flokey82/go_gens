# Go Generators
This repository contains various small attempts at procedural generation, simulation, and other things that might be useful for some people.

## Village Economy Generator
This package attempts to generate a stable, self sustaining economy given some initial buildings that serve as seed. The user can define building types that consume and produce resources, which the simulation will use to add new buildings until the sum of produced resources outweighs the sum of consumed resources. Please note that the output of this process is only meant to be used as input for a more sophisticated generator, only ensuring that besides required import of unavailable resources, the village/settlement can be in theory self-sustaining.

## Simple State Machine
This package provides a very basic/minimalistic state machine implemented in Go. It is neither efficient nor fast, but a start if you want to experiment with very, very, veeeery simple AI :) Feel free to fork and/or improve!
