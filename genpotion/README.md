# genpotion: Potion generator / Alchemy

This package implements a simple alchemy system similar to the potion crafting systems found in the Elder Scrolls games.

Please note that this code is inspired by two fantastic projects, please check them out!

* https://github.com/hogart/alchemy
* https://github.com/TemirkhanN/alchemist

# Design

Current scope:

- A potion is a combination of ingredients and provides one or multiple effects.
- The effects of a potion are determined by the ingredients and the effects that they share with each other.
- A potion can only be created if every ingredient shares at least one effect with at least one other ingredient.
- Effects unique to one ingredient are not considered.

Future scope:

- The strength of each effect is determined by the number of ingredients that have that particular effect.
- The strength of each effect is determined by the quality of the effect of the contributing ingredient.
- The price of a potion is dependent on the price of each ingredient as well as the quality and number of its effects.
- The strength of each effect is also dependent on the skill of the alchemist.
- The success of crafting a potion is also dependent on an ingredient.
- Special effects can be unlocked based on specific combinations of effects, ingredients, and/or skill.
- The name of a potion should be dependent on the ingredients
  - Apple, Sugar Cane: "Sweet Potion of Stamina"
  - Liquorice: "Bitter Draught of Magica"
  - etc.