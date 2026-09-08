package i18n

// LocationNamesPTBR mapeia os identificadores de locais para seus nomes formatados em Português (PT-BR).
var LocationNamesPTBR = map[LocationID]string{
	LocationTown:   "Praça do Vilarejo",
	LocationForest: "Floresta Sombria",
	LocationTavern: "Taverna da Dona Rosalinda",
	LocationChapel: "Capela do Frei Anselmo",
	LocationSmith:  "Ferraria do Mestre Torin",
	LocationGuild:  "Guilda dos Aventureiros",
	LocationDragon: "Covil do Dragão",
}

// NPCNamesPTBR mapeia os identificadores de NPCs para seus nomes canônicos e títulos em PT-BR.
var NPCNamesPTBR = map[NPCID]string{
	NPCRosalinda: "Dona Rosalinda, a Taverneira",
	NPCAnselmo:   "Frei Anselmo, o Curandeiro",
	NPCTorin:     "Mestre Torin, o Ferreiro",
	NPCCassandra: "Cassandra",
	NPCRedKnight: "Cavaleiro Vermelho",
	NPCTobias:    "Mestre Tobias, o Bibliotecário",
}

// MonsterNamesPTBR mapeia os IDs de espécies de monstros para seus nomes traduzidos em PT-BR.
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

	// Chefe Especial Final
	MonsterDragon: "O Dragão",
}

// DragonTitlesPTBR lista os subtítulos temáticos diários do Dragão.
//
// Didática Go: A seleção é determinística: o índice é derivado do hash SHA-256 da data corrente,
// garantindo que todos os jogadores vejam o mesmo título no mesmo dia.
var DragonTitlesPTBR = []string{
	"O Devorador de Goroutines",
	"O Terror dos Ponteiros",
	"A Fúria Ancestral de Gopher",
	"O Destruidor de Compiladores",
	"A Chama Vermelha do Abismo",
}

// UITextsPTBR mapeia as chaves de interface genéricas para suas traduções em PT-BR.
var UITextsPTBR = map[UIKey]string{
	UIMaxLevel: "(máx)",
}

// ItemNamesPTBR mapeia os IDs de itens e equipamentos para seus nomes formatados em PT-BR.
var ItemNamesPTBR = map[ItemID]string{
	WeaponStick:            "Pedaço de Pau",
	WeaponDagger:           "Adaga Afiada",
	WeaponShortSword:       "Espada Curta",
	WeaponBroadsword:       "Montante de Aço",
	WeaponDragonSlayer:     "Matadora de Dragões",
	ArmorClothes:           "Roupas Simples",
	ArmorLeather:           "Armadura de Couro",
	ArmorChainmail:         "Cota de Malha",
	ArmorPlate:             "Armadura de Placas",
	ArmorDragonScale:       "Escamas de Dragão",
	PotionHealth:           "Poção de Vida",
	PotionGarbageCollector: "Elixir do Garbage Collector",
}

// GetLocationName retorna o nome traduzido do local, provendo fallback para a string do ID caso não encontrado.
//
// Didática Go: Usamos a sintaxe `comma-ok` (`value, ok := map[key]`) para verificar com segurança a existência da chave no mapa.
func GetLocationName(id LocationID) string {
	if name, ok := LocationNamesPTBR[id]; ok {
		return name
	}
	return string(id)
}

// GetUIText retorna o texto traduzido de um elemento genérico de interface.
func GetUIText(key UIKey) string {
	if text, ok := UITextsPTBR[key]; ok {
		return text
	}
	return string(key)
}

// GetNPCName retorna o nome traduzido do NPC com seu título.
func GetNPCName(id NPCID) string {
	if name, ok := NPCNamesPTBR[id]; ok {
		return name
	}
	return string(id)
}

// GetMonsterName retorna o nome traduzido do monstro.
func GetMonsterName(id MonsterID) string {
	if name, ok := MonsterNamesPTBR[id]; ok {
		return name
	}
	return string(id)
}

// GetItemName retorna o nome traduzido do item.
func GetItemName(id ItemID) string {
	if name, ok := ItemNamesPTBR[id]; ok {
		return name
	}
	return string(id)
}
