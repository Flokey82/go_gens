# Gamerogueish

This is a sorta-kinda roguelike using the fantastic package https://github.com/BigJk/ramen, which is a simple console emulator written in go that can be used to create various ascii / text (roguelike) games.

Right now, the code is a really basic re-factor of the roguelike example that comes with Ramen, but I'll use it as a basis for various (at least to me) interesting experiments :)


![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/gamerogueish/images/rgb.png "rogue-ish")

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/gamerogueish/images/rgb2.png "next to the exit")

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/gamerogueish/images/rgb3.png "story panel")

## Keybindings

* ASDW - Move
* Space - Attack
* TAB - Toggle UI selection
* Arrow Up / Down - Select UI element
* Enter - Consume / Equip / Loot selected UI element
* Backspace - Drop selected UI element
* P - Add potion to inventory (for testing)
* T - Add trap to inventory (for testing)

## TODO

* Creatures
  * [DONE] Basic movement (random)
    * Custom movement speed
  * [DONE] AI (basic)
    * Custom perception radius
    * Should recognize traps and see items
  * [DONE] Pathfinding
    * Optimize
  * [DONE] Traps should affect enemies (hacky)
    * Only enemies that are within player view are affected
  * [DONE] Randomized loot
  * Generate names
* Documentation
* Inventory
  * [DONE] Basic inventory
  * [DONE] Item add / remove
* Items
  * [DONE] Basic items
  * [DONE] Item generation
  * [DONE] Consumable items
  * [DONE] Equippable items
  * [DONE] Enemy inventory
    * [DONE] Display items on dead enemies
    * [DONE] Looting of dead enemies
      * SELECTIVE looting of dead enemies
  * Item pickup / drop
    * [DONE] Item drop
      * Confirmation before dropping
    * [DONE] Item pickup
  * Item effects
  * [DONE] Item triggers
  * Hidden items
    * [DONE] Hidden traps
    * [DONE] Reveal hidden items on touch
    * Reveal hiddem items if we have high enough perception or a potion of perception
    * Hidden status should be per entity (player, enemy) so that enemies can also run into traps and avoid them when they know where they are
* Combat
  * [DONE] Player death
  * Enemies should attack on each turn
* Map generation
  * [DONE] Custom seed
  * [DONE] Custom world generator functions
  * [DONE] Creature placement
  * Item placement
  * [DONE] Dungeons
    * Neighbor rooms not centered (optionally)
    * Connections / doors not centered (optionally)
    * [DONE] Water puddles
      * Prevent water from blocking doors
  * Caves
  * Overworld / outdoor
* FOV / 'Fog of war'
  * [DONE] Basic radius based FOV
  * Raycasting based FOV
  * Merge map and FOV?
* Items / Entities rendering
  * Interface for items / entities in the world
* UI
  * Add scenes / screens
    * Scene interface
      * Input handling
      * Rendering
    * Scenes
      * Character creation
      * Main menu
      * [DONE] Main game
      * [DONE] Game over
      * [DONE] Win (located exit)
  * Deduplicate UI code (items, enemies, etc)
  * Highlight active UI components
  * Move player info out of selectable UI
  * [DONE] De-dupe scene textbox code
  * [DONE] Add scene textbox
    * [DONE] Pagination
    * [DONE] Line break / paragraph support
    * Move to top level (so that it can be used in all scenes)
    * Make closing key configurable (and customize bottom text accordingly)
* Gameplay
  * [DONE] Add win condition
  * Cutscenes / events / triggers
    * Trigger scenes through items (e.g. stairs, level exit)
      * [DONE] Level exit
      * [DONE] Trap
      * Text / message
    * Limit triggers to entity types (e.g. only player can trigger stairs)
  * Collects stats or score to show on win/lose

## Interesting stuff

* FOV
  * https://github.com/ajhager/rog/blob/master/fov.go
  * http://journal.stuffwithstuff.com/2015/09/07/what-the-hero-sees/
  * http://www.roguebasin.com/index.php?title=Field_of_Vision
* Loot
  * http://journal.stuffwithstuff.com/2014/07/05/dropping-loot/
  * https://www.reddit.com/r/roguelikedev/comments/2y3rkg/faq_friday_7_loot/
* Game loop
  * http://journal.stuffwithstuff.com/2014/07/15/a-turn-based-game-loop/
* Map generation
  * http://journal.stuffwithstuff.com/2014/12/21/rooms-and-mazes/