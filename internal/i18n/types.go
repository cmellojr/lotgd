package i18n

// LocationID represents distinct game locations in the village and surrounding areas.
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

// NPCID represents canonical NPCs in the village.
type NPCID string

const (
	NPCRosalinda NPCID = "rosalinda"
	NPCAnselmo   NPCID = "anselmo"
	NPCTorin     NPCID = "torin"
	NPCYolanda   NPCID = "yolanda"
	NPCCassandra NPCID = "cassandra"
	NPCRedKnight NPCID = "red_knight"
	NPCTobias    NPCID = "tobias"
)

// MonsterID represents base monster identifiers.
type MonsterID string

const (
	// Tier 1 - Beginner
	MonsterSewerRat     MonsterID = "sewer_rat"
	MonsterMossSpider   MonsterID = "moss_spider"
	MonsterClumsyKobold MonsterID = "clumsy_kobold"
	MonsterGreenSlime   MonsterID = "green_slime"
	MonsterSinisterCrow MonsterID = "sinister_crow"

	// Tier 2 - Apprentice
	MonsterGoblinScout   MonsterID = "goblin_scout"
	MonsterOrcRecruit    MonsterID = "orc_recruit"
	MonsterShadowWolf    MonsterID = "shadow_wolf"
	MonsterRoadBandit    MonsterID = "road_bandit"
	MonsterRustySkeleton MonsterID = "rusty_skeleton"

	// Tier 3 - Veteran
	MonsterBridgeTroll       MonsterID = "bridge_troll"
	MonsterRudeOgre          MonsterID = "rude_ogre"
	MonsterNoviceNecromancer MonsterID = "novice_necromancer"
	MonsterSingingHarpy      MonsterID = "singing_harpy"
	MonsterCrackedStoneGolem MonsterID = "cracked_stone_golem"

	// Tier 4 - Legendary
	MonsterFallenKnight     MonsterID = "fallen_knight"
	MonsterMountainWyvern   MonsterID = "mountain_wyvern"
	MonsterLesserLich       MonsterID = "lesser_lich"
	MonsterColdGazeBasilisk MonsterID = "cold_gaze_basilisk"
	MonsterSwampSpecter     MonsterID = "swamp_specter"

	// Special Boss
	MonsterDragon MonsterID = "the_dragon"
)

// ItemID represents weapons, armors, and consumables.
type ItemID string

const (
	// Weapons
	WeaponStick        ItemID = "stick"
	WeaponDagger       ItemID = "dagger"
	WeaponShortSword   ItemID = "short_sword"
	WeaponBroadsword   ItemID = "broadsword"
	WeaponDragonSlayer ItemID = "dragon_slayer"

	// Armors
	ArmorClothes     ItemID = "clothes"
	ArmorLeather     ItemID = "leather"
	ArmorChainmail   ItemID = "chainmail"
	ArmorPlate       ItemID = "plate"
	ArmorDragonScale ItemID = "dragon_scale"

	// Consumables & Easter Eggs
	PotionHealth           ItemID = "potion_health"
	PotionGarbageCollector ItemID = "potion_gc"
	WeaponNullPointer      ItemID = "null_pointer"
)
