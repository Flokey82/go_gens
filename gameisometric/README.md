# gameisometric

This is a simple isometric 'game' experiment using the ebitengine and is written in Go.
The code is heavily inspired by the [ebiten isometric example](https://ebitengine.org/en/examples/isometric.html). I modified the code to use [gendungeon](https://github.com/Flokey82/go_gens/gendungeon) for level generation instead of just randomly place tiles.

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/gameisometric/images/screenshot.png "Screenshot of isometric stuff!")

Some interesting info on the math behind isometric maths can be found here:
https://clintbellanger.net/articles/isometric_math/

## NOTE

Part of (or All) the graphic tiles used in this program is the public domain roguelike tileset 'RLTiles'.
You can find the original tileset at: http://rltiles.sf.net

See: https://opengameart.org/content/dungeon-crawl-32x32-tiles

Further, I am using Kenney's fantastic isometric tile set.

See: https://www.kenney.nl/assets/isometric-miniature-library

## TODO

- [X] Add basic isometric rendering
- [X] Use gendungeon for level generation
    - [X] Add multilevel dungeon generation
    - [X] Add stairs
    - [X] Add doors
    - [X] Add furniture
    - [ ] Tidy up level handling
    - [ ] Add clutter generation
- [ ] Add actual game
    - [ ] Add player and player movement
    - [ ] Add enemies
    - [ ] Add combat
    - [ ] Add items
    - [ ] Make stairs and doors functional
