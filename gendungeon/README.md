# Go Dungeon Generator

**NOTE: This code is a modified fork of https://github.com/brad811/go-dungeon/, which is an implementation of http://journal.stuffwithstuff.com/2014/12/21/rooms-and-mazes/.**

Please check out the original code as it has great features like a local web interface for rendering the dungeon. I stripped out a lot of the additional features, so please check out the original :)

## TODO

- [X] Add 2D dungeon generation
    - [ ] Add new room shapes
- [ ] Add multi-level dungeon generation
    - [X] Levels with identical dimensions
    - [ ] Levels with different dimensions
    - [ ] Add stairs
        - [X] Connect levels using stairs
        - [ ] Add custom number of stairs
        - [ ] Ensure we don't override existing stairs
    - [ ] Add constraint solver to prevent non-overlapping levels


![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/gendungeon/images/lvl0.png "Multilevel Dungeon!")

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/gendungeon/images/lvl1.png "Multilevel Dungeon!")
