# Gamerogueish

This is a sorta-kinda roguelike using the fantastic package https://github.com/BigJk/ramen, which is a simple console emulator written in go that can be used to create various ascii / text (roguelike) games.

Right now, the code is a really basic re-factor of the roguelike example that comes with Ramen, but I'll use it as a basis for various (at least to me) interesting experiments :)


![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/gamerogueish/images/rgb.png "rogue-ish")

## Done

* Custom world generator functions

## TODO

* FOV / 'Fog of war'
  * [DONE] Basic radius based FOV
  * Raycasting based FOV
  * See: https://github.com/ajhager/rog/blob/master/fov.go
* Creatures
  * AI
* Documentation
* Inventory
* Items
* Combat
* Map generation
  * Neighbor rooms not centered (optionally)
  * Connections / doors not centered (optionally)
  * Caves
  * Custom seed
