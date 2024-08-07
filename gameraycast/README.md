NOTE: This is currently a modified fork of https://github.com/samuel-pratt/raycaster, which will be reworked quite substantially to combine features of various raycasting implementations using ebiten and sdl.

- https://github.com/samuel-pratt/raycaster
- https://github.com/kyriacos/go-raycaster
- https://github.com/TheInvader360/dungeon-crawler
- https://github.com/Myu-Unix/ray_engine
- https://lodev.org/cgtutor/raycasting.html

This is an excellent collection of tutorials on raycasting:

- https://github.com/vinibiavatti1/RayCastingTutorial

Further, excellent resources on raycasting:

- https://lodev.org/cgtutor/raycasting.html

The ultimate goal is to create a basic dungeon crawler for the various dungeon and tile map generators that accumulate in this repository.

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/gameraycast/images/basic.png "Screenshot of basic raycast!")

## Textures

Some textures are from Wolfenstein 3D, which is a public domain game.

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/gameraycast/images/textured.png "Screenshot of textured raycast!")

## TODO

- [X] Basic raycasting, movement, rendering
- [X] Add different wall types and colors
- [X] Custom FOV and ray resolution
- [ ] Add textures
    - [X] Basic, single texture
    - [X] Multiple textures
- [ ] Add minimap
- [ ] Add custom maps
    - [X] Remove hardcoded map
    - [X] Add example dungeon map import
- [ ] Add sprites
- [ ] Add enemies
- [ ] Add map marker for player start
- [ ] Add doors
    - [X] Add door texture
    - [X] Add door interaction
    - [ ] Allow doors to close again instead of removing them
    - [ ] Add door opening / closing animation
- [ ] Add stairs