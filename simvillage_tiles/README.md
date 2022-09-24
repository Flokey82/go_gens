# SimVillage_Tiles

This package contains a very simple tile based renderer and will hopefully be a nifty village simulation soon.

## Features

* Simple tile based renderer
* Player controlled character
* NPC characters (rudimentary)
* Drawable prefab objects (rudimentary)
* Collision detection (rudimentary)
* Chunk loading / chunk generation
* Chunk caching (sorta)

## TODO

* Rendering
  * Improve tile render order
  * Find a better dungeon tile set
* Map / world
  * Decouple chunk size from viewport size (?)
  * Generation or loading of larger maps
  * Persistent world (since we use procgen, this'll be interesting)
  * Per-Tile actions / events (doors, triggers, ...)
  * Objects / resources / etc.
* [WIP] Better layer system / named layers
  * [DONE] Create new structs for handling map chunks and layers
  * [DONE] Migrate renderer to MapChunk and Layer types
  * Allow arbitrary layer names
  * Allow enabling per-layer collision detection via layer property
* [WIP] Indoor maps
  * [DONE] Map (world) switching
  * Doors with destination maps / worlds
  * Set player position when transitioning to new map
* Creatures
  * Stats
    * Use gamesheet for stats
  * AI
    * Perception system
    * Pathfinding
    * Decision making (state machines)
  * Actions
    * Open doors / enter buildings
    * Attack, pick up, ...

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/simvillage_tiles/images/rgb.png "Screenshot!")

## Attribution

Dungeon Tile Set by http://blog-buch.rhcloud.com (link is dead)