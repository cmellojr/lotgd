package i18n

// LocationNamesPTBR maps locations to Portuguese display names.
var LocationNamesPTBR = map[LocationID]string{
	LocationTown:   "Praça do Vilarejo",
	LocationForest: "Floresta Sombria",
	LocationTavern: "Taverna da Dona Rosalinda",
	LocationChapel: "Capela do Frei Anselmo",
	LocationSmith:  "Ferraria do Mestre Torin",
	LocationGuild:  "Guilda dos Aventureiros",
	LocationDragon: "Covil do Dragão",
}

// NPCNamesPTBR maps canonical NPCs to display names with titles.
var NPCNamesPTBR = map[NPCID]string{
	NPCRosalinda: "Dona Rosalinda, a Taverneira",
	NPCAnselmo:   "Frei Anselmo, o Curandeiro",
	NPCTorin:     "Mestre Torin, o Ferreiro",
	NPCYolanda:   "Yolanda, a Cigana Mercadora",
	NPCCassandra: "Cassandra",
	NPCRedKnight: "Cavaleiro Vermelho",
	NPCTobias:    "Mestre Tobias, o Bibliotecário",
}

// MonsterNamesPTBR maps monster IDs to canonical PT-BR names.
var MonsterNamesPTBR = map[MonsterID]string{
	// Tier 1
	MonsterSewerRat:     "Rato-do-Esgoto",
	MonsterMossSpider:   "Aranha-de-Musgo",
	MonsterClumsyKobold: "Kobold Trapalhão",
	MonsterGreenSlime:   "Slime Verde",
	MonsterSinisterCrow: "Corvo Sinistro",

	// Tier 2
	MonsterGoblinScout:   "Goblin Batedor",
	MonsterOrcRecruit:    "Orc Recruta",
	MonsterShadowWolf:    "Lobo das Sombras",
	MonsterRoadBandit:    "Bandido de Estrada",
	MonsterRustySkeleton: "Esqueleto Enferrujado",

	// Tier 3
	MonsterBridgeTroll:       "Troll da Ponte",
	MonsterRudeOgre:          "Ogro Malcriado",
	MonsterNoviceNecromancer: "Necromante Iniciante",
	MonsterSingingHarpy:      "Harpia Cantora",
	MonsterCrackedStoneGolem: "Golem de Pedra Rachada",

	// Tier 4
	MonsterFallenKnight:     "Cavaleiro Caído",
	MonsterMountainWyvern:   "Wyvern da Montanha",
	MonsterLesserLich:       "Lich Menor",
	MonsterColdGazeBasilisk: "Basilisco de Olhar Frio",
	MonsterSwampSpecter:     "Espectro do Pântano",

	// Special Boss
	MonsterDragon: "O Dragão",
}

// MonsterPrefixesPTBR contains adjective prefixes for procedural monsters.
var MonsterPrefixesPTBR = []string{
	"Feroz",
	"Covarde",
	"Enfurecido",
	"Sortudo",
	"Faminto",
	"Preguiçoso",
	"Astuto",
	"Gigantesco",
}

// ItemNamesPTBR maps item IDs to Portuguese names.
var ItemNamesPTBR = map[ItemID]string{
	WeaponStick:            "Pedaço de Pau",
	WeaponDagger:           "Adaga Afiada",
	WeaponShortSword:       "Espada Curta",
	WeaponBroadsword:       "Montante de Aço",
	WeaponDragonSlayer:     "Matadora de Dragões",
	WeaponNullPointer:      "Ponteiro Nulo Amaldiçoado",
	ArmorClothes:           "Roupas Simples",
	ArmorLeather:           "Armadura de Couro",
	ArmorChainmail:         "Cota de Malha",
	ArmorPlate:             "Armadura de Placas",
	ArmorDragonScale:       "Escamas de Dragão",
	PotionHealth:           "Poção de Vida",
	PotionGarbageCollector: "Elixir do Garbage Collector",
}

// GetLocationName returns the localized location name.
func GetLocationName(id LocationID) string {
	if name, ok := LocationNamesPTBR[id]; ok {
		return name
	}
	return string(id)
}

// GetNPCName returns the localized NPC name.
func GetNPCName(id NPCID) string {
	if name, ok := NPCNamesPTBR[id]; ok {
		return name
	}
	return string(id)
}

// GetMonsterName returns the localized monster name.
func GetMonsterName(id MonsterID) string {
	if name, ok := MonsterNamesPTBR[id]; ok {
		return name
	}
	return string(id)
}

// GetItemName returns the localized item name.
func GetItemName(id ItemID) string {
	if name, ok := ItemNamesPTBR[id]; ok {
		return name
	}
	return string(id)
}
