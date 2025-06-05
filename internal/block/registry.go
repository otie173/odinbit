package block

import (
	"log"

	"github.com/otie173/odinbit/resource"
)

var (
	defaultBehavior func() = func() {
		log.Println("Block was triggered")
	}

	Barrier                        Block
	Wall                           Block
	Floor1                         Block
	Door1, Door2                   Block
	Grass1, Grass2, Grass3, Grass4 Block
	Grass5, Grass6, Grass7, Grass8 Block
	Tree1, Tree2, Tree3, Tree4     Block
	Stone1, Stone2, Stone3         Block
	Mushroom1, Mushroom2           Block
)

func Load() {
	Barrier = Block{
		ID:       1,
		Texture:  resource.BarrierTexture,
		Passable: false,
		Behavior: defaultBehavior,
	}

	Wall = Block{
		ID:       2,
		Texture:  resource.WallTexture,
		Passable: false,
		Behavior: defaultBehavior,
	}

	Floor1 = Block{
		ID:       3,
		Texture:  resource.Floor1Texture,
		Passable: true,
		Behavior: defaultBehavior,
	}

	Door1 = Block{
		ID:       4,
		Texture:  resource.Door1Texture,
		Passable: false,
		Behavior: defaultBehavior,
	}

	Door2 = Block{
		ID:       5,
		Texture:  resource.Door2Texture,
		Passable: true,
		Behavior: defaultBehavior,
	}

	Grass1 = Block{
		ID:       6,
		Texture:  resource.GrassTextures[0],
		Passable: true,
		Behavior: defaultBehavior,
	}

	Grass2 = Block{
		ID:       7,
		Texture:  resource.GrassTextures[1],
		Passable: true,
		Behavior: defaultBehavior,
	}

	Grass7 = Block{
		ID:       8,
		Texture:  resource.GrassTextures[6],
		Passable: true,
		Behavior: defaultBehavior,
	}

	Grass8 = Block{
		ID:       9,
		Texture:  resource.GrassTextures[7],
		Passable: true,
		Behavior: defaultBehavior,
	}

	Grass3 = Block{
		ID:       10,
		Texture:  resource.GrassTextures[2],
		Passable: false,
		Behavior: defaultBehavior,
	}

	Grass4 = Block{
		ID:       11,
		Texture:  resource.GrassTextures[3],
		Passable: false,
		Behavior: defaultBehavior,
	}

	Grass5 = Block{
		ID:       12,
		Texture:  resource.GrassTextures[4],
		Passable: false,
		Behavior: defaultBehavior,
	}

	Grass6 = Block{
		ID:       13,
		Texture:  resource.GrassTextures[5],
		Passable: false,
		Behavior: defaultBehavior,
	}

	Tree1 = Block{
		ID:       14,
		Texture:  resource.TreeTextures[0],
		Passable: false,
		Behavior: defaultBehavior,
	}

	Tree2 = Block{
		ID:       15,
		Texture:  resource.TreeTextures[1],
		Passable: false,
		Behavior: defaultBehavior,
	}

	Tree3 = Block{
		ID:       16,
		Texture:  resource.TreeTextures[2],
		Passable: false,
		Behavior: defaultBehavior,
	}

	Tree4 = Block{
		ID:       17,
		Texture:  resource.TreeTextures[3],
		Passable: false,
		Behavior: defaultBehavior,
	}

	Stone1 = Block{
		ID:       18,
		Texture:  resource.StoneTextures[0],
		Passable: false,
		Behavior: defaultBehavior,
	}

	Stone2 = Block{
		ID:       19,
		Texture:  resource.StoneTextures[1],
		Passable: false,
		Behavior: defaultBehavior,
	}

	Stone3 = Block{
		ID:       20,
		Texture:  resource.StoneTextures[2],
		Passable: false,
		Behavior: defaultBehavior,
	}

	Mushroom1 = Block{
		ID:       21,
		Texture:  resource.MushroomTextures[0],
		Passable: false,
		Behavior: defaultBehavior,
	}

	Mushroom2 = Block{
		ID:       22,
		Texture:  resource.MushroomTextures[1],
		Passable: false,
		Behavior: defaultBehavior,
	}
}
