# dragonfly-world-generation

A WIP world generator for dragonfly.

Progress:
  ☑ Biomes
  ☑ Terrain shape
  ☐ Surface and correct materials
  ☑ Caves & pillars
  ☐ Terrain features
  ☐ Structures
  ☐ Nether and End dimensions

## Setup

Include this in your `main()` function:

```golang
[...]

// First load resources for surface definitions, terrain features and structures
// These resources can be found in the official BDS installation,
// but you don't need the whole install just these two folders.
worldgen.LoadResources("/home/<USER>/BDS/definitions", "/home/<USER>/BDS/behavior_packs/vanilla/structures")

[...]

// Then slap your favourite seed in this function and your done!
conf.Generator = worldgen.NewDefaultGenerator(123)

// Generators with more options are coming soon!

[...]
```