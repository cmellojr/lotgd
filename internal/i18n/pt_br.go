package i18n

import "fmt"

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

// AffixNamesPTBR mapeia os IDs de afixos procedurais para seus nomes em PT-BR.
var AffixNamesPTBR = map[AffixID]string{
	AffixFerocious: "Feroz",
	AffixCowardly:  "Covarde",
	AffixEnraged:   "Enfurecido",
	AffixLucky:     "Sortudo",
	AffixHungry:    "Faminto",
	AffixLazy:      "Preguiçoso",
	AffixShrewd:    "Astuto",
	AffixGigantic:  "Gigantesco",
}

// UITextsPTBR mapeia as chaves de interface genéricas e estáticas para suas traduções em PT-BR.
var UITextsPTBR = map[UIKey]string{
	// Barra de Status e Cabeçalhos
	UIHero:        "Herói:",
	UILevel:       "Nível:",
	UIHealth:      "HP:",
	UIGold:        "Ouro:",
	UIBank:        "Banco:",
	UIDailyFights: "Lutas Diárias:",
	UIWeapon:      "Arma:",
	UIArmor:       "Armadura:",
	UIPotions:     "Poções:",
	UIAttack:      "ATK:",
	UIDefense:     "DEF:",
	UIExperience:  "XP:",
	UIMaxLevel:    "(máx)",

	// Rodapés e Dicas de Atalho
	UIFooterNav:           "[↑/↓] Navegar  [Enter] Selecionar",
	UIFooterNavBack:       "[↑/↓] Navegar  [Enter] Selecionar  [ESC] Voltar",
	UIFooterCombat:        "[A] Atacar  [P] Usar Poção  [F] Fugir",
	UIFooterLogin:         "[Tab] Alternar Campo  [Enter] Confirmar  [Ctrl+C] Sair",
	UIFooterGameOver:      "[Enter/Espaço] Renascer na Capela",
	UIFooterSmith:         "[1/2/3] Abas  [↑/↓] Navegar  [Enter] Comprar/Vender  [ESC] Voltar",

	// Títulos e Subtítulos
	UITownTitle:        "Praça do Vilarejo",
	UITownSubtitle:     "O coração do reino, cercado por tavernas, lojas e o perigo da floresta.",
	UIForestTitle:      "Floresta Sombria",
	UIForestSubtitle:   "Árvores retorcidas e monstros à espreita.",
	UISmithTitle:       "Ferraria do Mestre Torin",
	UISmithSubtitle:    "O cheiro de carvão e o som do martelo contra a bigorna preenchem a oficina.",
	UIDragonTitle:      "Covil do Dragão Ancestral",
	UIDragonSubtitle:   "Uma caverna fumegante recheada de ossos de antigos heróis...",
	UITavernTitle:      "Taverna da Dona Rosalinda",
	UITavernSubtitle:   "Risadas, hidromel e histórias de aventureiros preenchem o ambiente.",
	UIChapelTitle:      "Capela do Frei Anselmo",
	UIChapelSubtitle:   "Um refúgio de paz e cura espiritual no centro do vilarejo.",
	UIGuildTitle:       "Guilda dos Aventureiros",
	UIGuildSubtitle:    "Onde guerreiros treinam sob a tutela de mestres lendários.",
	UIGameOverTitle:    "VOCÊ MORREU!",
	UIGameOverSubtitle: "Sua jornada chegou a um fim trágico...",
	UILoginTitle:       "THE LEGEND OF THE GO DRAGON",
	UILoginSubtitle:    "Entre com seu nome e senha de aventureiro",

	// Menus da Praça
	UITownMenuForest:       "Explorar a Floresta Sombria",
	UITownMenuSmith:        "Visitar a Ferraria do Mestre Torin",
	UITownMenuTavern:       "Entrar na Taverna da Dona Rosalinda",
	UITownMenuChapel:       "Visitar a Capela do Frei Anselmo",
	UITownMenuGuild:        "Ir à Guilda dos Aventureiros",
	UITownMenuDragon:       "Desafiar o Dragão Ancestral",
	UITownMenuBankDeposit:  "Depositar todo o ouro no Banco",
	UITownMenuBankWithdraw: "Sacar todo o ouro do Banco",
	UITownMenuQuit:         "Sair do Jogo",
	UITownDragonReq:        "(Requer Nível %d)",

	// Floresta
	UIForestMenuLookForFight: "Procurar por Combate",
	UIForestMenuReturnTown:   "Retornar à Praça do Vilarejo",
	UIForestNoFightsLeft:     "Você já gastou todas as suas lutas diárias! Volte amanhã para novos desafios.",

	// Ferraria
	UISmithTabWeapons:      "1. Armas",
	UISmithTabArmors:       "2. Armaduras",
	UISmithTabPotions:      "3. Poções",
	UISmithAlreadyEquipped: "Você já possui este equipamento equipado!",
	UISmithNotEnoughGold:   "Você não possui ouro suficiente para esta compra!",
	UISmithNoItemToSell:    "Você não possui nenhum equipamento equipado nesta categoria para vender.",
	UISmithFooterConfirm:   "[S] Aceitar Oferta  [N] Recusar Oferta",

	// Covil do Dragão
	UIDragonAlreadySlain: "O Dragão do Dia já foi derrotado por outro herói! Volte amanhã para o novo surgimento da fera.",

	// Taverna
	UITavernFlirt:         "Flerte com Cassandra (+1 Luta Diária - 10 Ouro)",
	UITavernTalkRedKnight: "Conversar com o Cavaleiro Vermelho",
	UITavernReturnTown:    "Voltar para a Praça",
	UITavernFlirtLimit:    "Cassandra já te inspirou hoje. Volte amanhã para nova dose de coragem!",
	UITavernFlirtSuccess:  "Você oferece uma bebida a Cassandra. Ela sorri graciosamente e sua determinação é renovada! (+1 Luta na Floresta!)",
	UITavernRedKnightMsg:  "O Cavaleiro Vermelho ergue a viseira: 'Você ainda não possui a tempera necessária para cruzar lâminas comigo. Volte quando estiver mais experiente!'",

	// Capela
	UIChapelFullHP:        "Você já está com a saúde máxima! Frei Anselmo abençoa sua caminhada.",
	UIChapelNotEnoughGold: "Você não tem moedas suficientes para a doação da cura.",
	UIChapelHealSuccess:   "Frei Anselmo unge seus ferimentos com óleos sagrados. Vida totalmente restaurada!",

	// Guilda
	UIGuildMaxLevelReached: "Você alcançou o topo da maestria e aprendeu tudo com os mestres da guilda!",

	// Game Over
	UIGameOverRespawnInfo: "Frei Anselmo encontrou seu corpo desacordado e te ressuscitou na Capela com 1 HP.",

	// Login
	UILoginUsernameLabel:  "Aventureiro:",
	UILoginPasswordLabel:  "Senha:",
	UILoginErrInvalidPass: "Senha incorreta para o aventureiro.",
	UILoginErrUserExists:  "Este nome de aventureiro já está registrado.",
}

// MessageTemplatesPTBR mapeia as chaves de mensagens parametrizadas para seus templates em PT-BR.
var MessageTemplatesPTBR = map[MessageKey]string{
	// Banco
	MsgBankDepositSuccess:  "Você depositou %d moedas de ouro no cofre com segurança!",
	MsgBankDepositInvalid:  "A quantia de depósito deve ser positiva.",
	MsgBankDepositNoGold:   "Ouro insuficiente na bolsa para depósito (disponível: %d).",
	MsgBankWithdrawSuccess: "Você retirou %d moedas de ouro do seu cofre.",
	MsgBankWithdrawInvalid: "A quantia de saque deve ser positiva.",
	MsgBankWithdrawNoGold:  "Saldo bancário insuficiente para saque (disponível: %d).",

	// Combate
	MsgCombatPlayerWinDragon: "VITÓRIA LENDÁRIA! Você desferiu o golpe fatal e derrotou %s!",
	MsgCombatPlayerWin:       "Você derrotou %s e ganhou %d XP e %d moedas de ouro!",
	MsgCombatPlayerDefeat:    "%s desferiu um golpe mortal! Você sucumbiu na escuridão...",
	MsgCombatTurnRound:       "Você causou %d de dano%s. %s contra-atacou causando %d de dano.",
	MsgCombatFleeSuccess:     "Você conseguiu recuar estrategicamente para as sombras e escapar de %s!",
	MsgCombatFleeFailDefeat:  "Você tropeçou ao tentar fugir! %s aproveitou e desferiu um golpe letal!",
	MsgCombatFleeFail:        "Falha ao fugir! %s bloqueou seu caminho e te atingiu causando %d de dano.",
	MsgCombatNoPotions:       "Você não possui poções na bolsa.",
	MsgCombatFullHealth:      "Sua vida já está cheia.",

	// Ferraria
	MsgSmithBuyPrompt:          "Comprar %s (%d Ouro)",
	MsgSmithSellPrompt:         "Vender %s (%d Ouro - Recompra)",
	MsgSmithBuySuccess:         "Você comprou %s por %d moedas de ouro!",
	MsgSmithSellSuccess:        "Você vendeu %s por %d moedas de ouro!",
	MsgSmithSellOfferPrompt:   "Mestre Torin examina sua %s: 'Te dou %d moedas de ouro por ela. Aceita?' [S]im / [N]ão",
	MsgSmithSellOfferRejected: "Você recusou a oferta de Mestre Torin. A %s permanece equipada.",

	// Dragão
	MsgDragonReqNotMet: "Sua coragem é grande, mas você precisa ter pelo menos Nível %d para entrar no Covil do Dragão!",

	// Capela
	MsgChapelHealPrompt: "Receber Bênção de Cura (%d Ouro)",

	// Guilda
	MsgGuildChallengePrompt:    "Pronto para desafiar %s pelo Nível %d!",
	MsgGuildNotEnoughXP:        "Experiência insuficiente. Necessário: %d XP (Você tem: %d XP)",
	MsgGuildAlreadyFoughtToday: "Você já desafiou %s hoje! Apenas 1 tentativa por dia é permitida. Retorne no próximo alvorecer.",
	MsgGuildVictory:            "PARABÉNS! Você derrotou %s e ascendeu ao Nível %d!",

	// Game Over
	MsgGameOverLostGold: "Ouro perdido na bolsa: %d moedas",
	MsgGameOverLostXP:   "Experiência perdida: %d XP",

	// Login
	MsgLoginErrCreateAccount: "Erro ao criar conta: %v",
	MsgLoginErrAuth:          "Erro ao autenticar: %v",
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

// ItemDescriptionsPTBR mapeia os IDs de itens para suas descrições traduzidas em PT-BR.
var ItemDescriptionsPTBR = map[ItemID]string{
	WeaponStick:            "Um galho seco encontrado no chão. Melhor que lutar desarmado.",
	WeaponDagger:           "Pequena e afiada, ótima para perfurar pontos fracos de monstros menores.",
	WeaponShortSword:       "Espada confiável de ferro forjada na oficina do Mestre Torin.",
	WeaponBroadsword:       "Lâmina pesada de aço temperado capaz de cortar escamas grossas.",
	WeaponDragonSlayer:     "Arma lendária banhada em sangue de monstros antigos, feita para abater dragões.",
	ArmorClothes:           "Roupas rotas de plebeu que não oferecem proteção real contra presas e garras.",
	ArmorLeather:           "Colete de couro curtido, leve e flexível para iniciantes na floresta.",
	ArmorChainmail:         "Anéis entrelaçados de ferro que absorvem impactos de cortes e flechadas.",
	ArmorPlate:             "Placas maciças de aço polido que protegem o tórax contra golpes devastadores.",
	ArmorDragonScale:       "Armadura forjada com as escamas impenetráveis do próprio Dragão Ancestral.",
	PotionHealth:           "Frasco com líquido vermelho brilhante que cura ferimentos imediatamente.",
	PotionGarbageCollector: "Poção suprema que limpa todas as imperfeições e restaura a saúde completamente.",
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

// GetAffixName retorna o nome traduzido do afixo procedural.
func GetAffixName(id AffixID) string {
	if name, ok := AffixNamesPTBR[id]; ok {
		return name
	}
	return string(id)
}

// GetItemDescription retorna a descrição traduzida do item.
func GetItemDescription(id ItemID) string {
	if desc, ok := ItemDescriptionsPTBR[id]; ok {
		return desc
	}
	return string(id)
}

// GetMessage retorna a mensagem formatada a partir da chave e dos argumentos opcionais.
func GetMessage(key MessageKey, args ...any) string {
	tmpl, ok := MessageTemplatesPTBR[key]
	if !ok {
		return string(key)
	}
	if len(args) == 0 {
		return tmpl
	}
	return fmt.Sprintf(tmpl, args...)
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
