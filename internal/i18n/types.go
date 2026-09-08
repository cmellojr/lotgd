package i18n

// LocationID representa os identificadores únicos de locais no vilarejo e arredores.
//
// Didática Go: Criamos um alias de tipo `type LocationID string` sobre o tipo primitivo string
// para garantir segurança de tipos em tempo de compilação (type safety), impedindo a substituição
// acidental de IDs de locais por outras strings.
type LocationID string

const (
	LocationTown   LocationID = "town"
	LocationForest LocationID = "forest"
	LocationTavern LocationID = "tavern"
	LocationChapel LocationID = "chapel"
	LocationSmith  LocationID = "smith"
	LocationGuild  LocationID = "guild"
	LocationDragon LocationID = "dragon"
)

// NPCID representa os identificadores dos NPCs canônicos do vilarejo.
type NPCID string

const (
	NPCRosalinda NPCID = "rosalinda"
	NPCAnselmo   NPCID = "anselmo"
	NPCTorin     NPCID = "torin"
	NPCCassandra NPCID = "cassandra"
	NPCRedKnight NPCID = "red_knight"
	NPCTobias    NPCID = "tobias"
)

// MonsterID representa os identificadores das espécies de monstros.
type MonsterID string

const (
	// Tier 1 - Iniciantes (Níveis 1 a 2)
	MonsterSewerRat     MonsterID = "sewer_rat"
	MonsterMossSpider   MonsterID = "moss_spider"
	MonsterClumsyKobold MonsterID = "clumsy_kobold"
	MonsterGreenSlime   MonsterID = "green_slime"
	MonsterSinisterCrow MonsterID = "sinister_crow"

	// Tier 2 - Aventureiros (Níveis 3 a 4)
	MonsterGoblinScout   MonsterID = "goblin_scout"
	MonsterOrcRecruit    MonsterID = "orc_recruit"
	MonsterShadowWolf    MonsterID = "shadow_wolf"
	MonsterRoadBandit    MonsterID = "road_bandit"
	MonsterRustySkeleton MonsterID = "rusty_skeleton"

	// Tier 3 - Veteranos (Níveis 5 a 7)
	MonsterBridgeTroll       MonsterID = "bridge_troll"
	MonsterRudeOgre          MonsterID = "rude_ogre"
	MonsterNoviceNecromancer MonsterID = "novice_necromancer"
	MonsterSingingHarpy      MonsterID = "singing_harpy"
	MonsterCrackedStoneGolem MonsterID = "cracked_stone_golem"

	// Tier 4 - Lendários (Níveis 8 a 10)
	MonsterFallenKnight     MonsterID = "fallen_knight"
	MonsterMountainWyvern   MonsterID = "mountain_wyvern"
	MonsterLesserLich       MonsterID = "lesser_lich"
	MonsterColdGazeBasilisk MonsterID = "cold_gaze_basilisk"
	MonsterSwampSpecter     MonsterID = "swamp_specter"

	// Chefe Especial Final
	MonsterDragon MonsterID = "the_dragon"
)

// ItemID representa os identificadores de armas, armaduras e consumíveis do jogo.
type ItemID string

const (
	// Armas
	WeaponStick        ItemID = "stick"
	WeaponDagger       ItemID = "dagger"
	WeaponShortSword   ItemID = "short_sword"
	WeaponBroadsword   ItemID = "broadsword"
	WeaponDragonSlayer ItemID = "dragon_slayer"

	// Armaduras
	ArmorClothes     ItemID = "clothes"
	ArmorLeather     ItemID = "leather"
	ArmorChainmail   ItemID = "chainmail"
	ArmorPlate       ItemID = "plate"
	ArmorDragonScale ItemID = "dragon_scale"

	// Consumíveis e Poções
	PotionHealth           ItemID = "potion_health"
	PotionGarbageCollector ItemID = "potion_gc"
)
