# genflora

THIS IS A WORK IN PROGRESS, DO NOT USE.

A climate produces a number of plants that can survive in a given climate. This package provides a simple interface to generate random plant lineages based on a given climate.

Properties like leaf shapes, stem diameter, flower color, etc. influence how much sunlight, precipitation, moisture, and temperature a plant lineage requires.

Base properties are used to seed the initial population for each lineage.

There are presets for plant groups like trees, shrubs, herbs, grasses and others which are used (with slight variations and permutations) to generate random plant lineages.

Mutations are used to introduce new traits into the population in each generation.

A fitnessfunction is used to calculate the fitness of a variant plant lineage. The fitness is used to determine the probability of a plant lineage to be selected for the next generation (or, to survive in the given climate).

## Braindump

Note that this is just rambling without scientific basis. The points below are based on guesswork and Googling some facts and figures.

### Surviving drought

Thick leaves (like from succulents) and thick stems (like from cacti) can store water during drought periods. However, they need to be more fleshy than wooden.

In arid areas, plants that store water usually have thorns to prevent animals from nom nom nomming them.

### Surviving heat

To avoid loss of moisture through evaporation, hot climates produce minimal surface area (small, fleshy leaves, thick stems, no leaves, long but thin(-ish) leaves)

### Surviving cold

In cold climates, low water content in the plant protects the plant from freezing. Hardy, woody stems are one way to help deal with the cold.

If the cold temperature is limited to winters, plants might drop their leaves to reduce water requirements and create a natural insulator to protect the roots from the cold.

Another strategy is to produce a lot of light seeds that can be dispersed by the wind and grow in the spring. The plant itself is not expected to survive the winter.

Plants might grow low and with high density to protect the roots from the cold.

### Surviving on steep slopes

On steep slopes, the roots can't grow very deep, so smaller plants with shallow roots are more likely to survive. If the soil depth is sufficient taller plants might survive.

### Surviving in the shade

Plants that grow in shaded areas or under trees have to be able to survive with less sunlight. Funghi are one example of plants that can survive in the shade or without light. Also broad, thin leaves are a good strategy to absorb as much light as possible.


### Other

Dry climate might result in:
 - small leaves to minimize water loss
   - NOTE: (This article)[https://www.frontiersin.org/articles/10.3389/fpls.2019.00058/full] suggests the opposite. In temperate forests, smaller leaves seem to lose water faster. Weird stuff.
 - small flowers to minimize water loss
 - hardy stems that can withstand drought
 - thick stems that can hold a lot of water

Wet climate might result in:
 - large leaves to maximize photosynthesis
 - large flowers to attract pollinators
 - flexible stems that can withstand flooding

Steep incline might result in:
 - deep roots to anchor the plant
 - short height to minimize wind resistance and lower the center of gravity

Shallow soil might result in:
 - shallow roots to avoid hitting bedrock
 - low height for stability

Low nutrient soil might result in:
 - stunted growth
 - low height

Cold winters might result in:
 - trees and shrubs drop leaves

