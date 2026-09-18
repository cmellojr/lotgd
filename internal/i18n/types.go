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

// AffixID representa os identificadores dos afixos procedurais de monstros.
type AffixID string

const (
	AffixFerocious AffixID = "ferocious"
	AffixCowardly  AffixID = "cowardly"
	AffixEnraged   AffixID = "enraged"
	AffixLucky     AffixID = "lucky"
	AffixHungry    AffixID = "hungry"
	AffixLazy      AffixID = "lazy"
	AffixShrewd    AffixID = "shrewd"
	AffixGigantic  AffixID = "gigantic"
)

// UIKey representa chaves de tradução de elementos genéricos e estáticos da interface.
type UIKey string

const (
	// Barra de Status e Cabeçalhos
	UIHero        UIKey = "hero"
	UILevel       UIKey = "level"
	UIHealth      UIKey = "health"
	UIGold        UIKey = "gold"
	UIBank        UIKey = "bank"
	UIDailyFights UIKey = "daily_fights"
	UIWeapon      UIKey = "weapon"
	UIArmor       UIKey = "armor"
	UIPotions     UIKey = "potions"
	UIAttack      UIKey = "attack"
	UIDefense     UIKey = "defense"
	UIExperience  UIKey = "experience"
	UIMaxLevel    UIKey = "max_level"

	// Rodapés e Dicas de Atalho
	UIFooterNav      UIKey = "footer_nav"
	UIFooterNavBack  UIKey = "footer_nav_back"
	UIFooterCombat   UIKey = "footer_combat"
	UIFooterLogin    UIKey = "footer_login"
	UIFooterGameOver UIKey = "footer_game_over"
	UIFooterSmith    UIKey = "footer_smith"

	// Títulos e Subtítulos de Telas
	UITownTitle        UIKey = "town_title"
	UITownSubtitle     UIKey = "town_subtitle"
	UIForestTitle      UIKey = "forest_title"
	UIForestSubtitle   UIKey = "forest_subtitle"
	UISmithTitle       UIKey = "smith_title"
	UISmithSubtitle    UIKey = "smith_subtitle"
	UIDragonTitle      UIKey = "dragon_title"
	UIDragonSubtitle   UIKey = "dragon_subtitle"
	UITavernTitle      UIKey = "tavern_title"
	UITavernSubtitle   UIKey = "tavern_subtitle"
	UIChapelTitle      UIKey = "chapel_title"
	UIChapelSubtitle   UIKey = "chapel_subtitle"
	UIGuildTitle       UIKey = "guild_title"
	UIGuildSubtitle    UIKey = "guild_subtitle"
	UIGameOverTitle    UIKey = "game_over_title"
	UIGameOverSubtitle UIKey = "game_over_subtitle"
	UILoginTitle       UIKey = "login_title"
	UILoginSubtitle    UIKey = "login_subtitle"

	// Menus da Praça
	UITownMenuForest       UIKey = "town_menu_forest"
	UITownMenuSmith        UIKey = "town_menu_smith"
	UITownMenuTavern       UIKey = "town_menu_tavern"
	UITownMenuChapel       UIKey = "town_menu_chapel"
	UITownMenuGuild        UIKey = "town_menu_guild"
	UITownMenuDragon       UIKey = "town_menu_dragon"
	UITownMenuBankDeposit  UIKey = "town_menu_bank_deposit"
	UITownMenuBankWithdraw UIKey = "town_menu_bank_withdraw"
	UITownMenuQuit         UIKey = "town_menu_quit"
	UITownDragonReq        UIKey = "town_dragon_req"

	// Floresta
	UIForestMenuLookForFight UIKey = "forest_menu_look_for_fight"
	UIForestMenuReturnTown   UIKey = "forest_menu_return_town"
	UIForestNoFightsLeft     UIKey = "forest_no_fights_left"

	// Ferraria
	UISmithTabWeapons      UIKey = "smith_tab_weapons"
	UISmithTabArmors       UIKey = "smith_tab_armors"
	UISmithTabPotions      UIKey = "smith_tab_potions"
	UISmithAlreadyEquipped UIKey = "smith_already_equipped"
	UISmithNotEnoughGold   UIKey = "smith_not_enough_gold"
	UISmithNoItemToSell    UIKey = "smith_no_item_to_sell"
	UISmithFooterConfirm   UIKey = "smith_footer_confirm"

	// Covil do Dragão
	UIDragonAlreadySlain UIKey = "dragon_already_slain"

	// Taverna
	UITavernFlirt         UIKey = "tavern_flirt"
	UITavernTalkRedKnight UIKey = "tavern_talk_red_knight"
	UITavernReturnTown    UIKey = "tavern_return_town"
	UITavernFlirtLimit    UIKey = "tavern_flirt_limit"
	UITavernFlirtSuccess  UIKey = "tavern_flirt_success"
	UITavernRedKnightMsg  UIKey = "tavern_red_knight_msg"

	// Capela
	UIChapelFullHP        UIKey = "chapel_full_hp"
	UIChapelNotEnoughGold UIKey = "chapel_not_enough_gold"
	UIChapelHealSuccess   UIKey = "chapel_heal_success"

	// Guilda
	UIGuildMaxLevelReached UIKey = "guild_max_level_reached"

	// Game Over
	UIGameOverRespawnInfo UIKey = "game_over_respawn_info"

	// Login
	UILoginUsernameLabel  UIKey = "login_username_label"
	UILoginPasswordLabel  UIKey = "login_password_label"
	UILoginErrInvalidPass UIKey = "login_err_invalid_pass"
	UILoginErrUserExists  UIKey = "login_err_user_exists"
)

// MessageKey representa chaves de mensagens parametrizadas e dinâmicas.
type MessageKey string

const (
	// Banco
	MsgBankDepositSuccess  MessageKey = "bank_deposit_success"
	MsgBankDepositInvalid  MessageKey = "bank_deposit_invalid"
	MsgBankDepositNoGold   MessageKey = "bank_deposit_no_gold"
	MsgBankWithdrawSuccess MessageKey = "bank_withdraw_success"
	MsgBankWithdrawInvalid MessageKey = "bank_withdraw_invalid"
	MsgBankWithdrawNoGold  MessageKey = "bank_withdraw_no_gold"

	// Combate
	MsgCombatPlayerWinDragon MessageKey = "combat_player_win_dragon"
	MsgCombatPlayerWin       MessageKey = "combat_player_win"
	MsgCombatPlayerDefeat    MessageKey = "combat_player_defeat"
	MsgCombatTurnRound       MessageKey = "combat_turn_round"
	MsgCombatFleeSuccess     MessageKey = "combat_flee_success"
	MsgCombatFleeFailDefeat  MessageKey = "combat_flee_fail_defeat"
	MsgCombatFleeFail        MessageKey = "combat_flee_fail"
	MsgCombatNoPotions       MessageKey = "combat_no_potions"
	MsgCombatFullHealth      MessageKey = "combat_full_health"

	// Ferraria
	MsgSmithBuyPrompt         MessageKey = "smith_buy_prompt"
	MsgSmithSellPrompt        MessageKey = "smith_sell_prompt"
	MsgSmithBuySuccess        MessageKey = "smith_buy_success"
	MsgSmithSellSuccess       MessageKey = "smith_sell_success"
	MsgSmithSellOfferPrompt   MessageKey = "smith_sell_offer_prompt"
	MsgSmithSellOfferRejected MessageKey = "smith_sell_offer_rejected"

	// Dragão
	MsgDragonReqNotMet MessageKey = "dragon_req_not_met"

	// Capela
	MsgChapelHealPrompt MessageKey = "chapel_heal_prompt"
	MsgHealSuccess      MessageKey = "heal_success"

	// Guilda
	MsgGuildChallengePrompt    MessageKey = "guild_challenge_prompt"
	MsgGuildNotEnoughXP        MessageKey = "guild_not_enough_xp"
	MsgGuildAlreadyFoughtToday MessageKey = "guild_already_fought_today"
	MsgGuildVictory            MessageKey = "guild_victory"

	// Game Over
	MsgGameOverLostGold MessageKey = "game_over_lost_gold"
	MsgGameOverLostXP   MessageKey = "game_over_lost_xp"

	// Login
	MsgLoginErrCreateAccount MessageKey = "login_err_create_account"
	MsgLoginErrAuth          MessageKey = "login_err_auth"
)
