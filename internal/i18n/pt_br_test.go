package i18n

import "testing"

func TestLocalizationMaps(t *testing.T) {
	if name := GetLocationName(LocationTown); name != "Praça do Vilarejo" {
		t.Errorf("expected 'Praça do Vilarejo', got %q", name)
	}

	if npc := GetNPCName(NPCRosalinda); npc != "Dona Rosalinda, a Taverneira" {
		t.Errorf("expected Rosalinda title, got %q", npc)
	}

	if mon := GetMonsterName(MonsterSewerRat); mon != "Rato-do-Esgoto" {
		t.Errorf("expected 'Rato-do-Esgoto', got %q", mon)
	}

	if item := GetItemName(WeaponDragonSlayer); item != "Matadora de Dragões" {
		t.Errorf("expected 'Matadora de Dragões', got %q", item)
	}

	if unknown := GetLocationName(LocationID("unknown_loc")); unknown != "unknown_loc" {
		t.Errorf("expected fallback to ID for unknown location, got %q", unknown)
	}
}

// TestUITextsCoverAllKeys garante que toda chave de interface e mensagem declarada
// tem tradução registrada, conforme a verificação prevista no ADR-0005.
func TestUITextsCoverAllKeys(t *testing.T) {
	uiKeys := []UIKey{
		UIHero, UILevel, UIHealth, UIGold, UIBank, UIDailyFights,
		UIWeapon, UIArmor, UIPotions, UIAttack, UIDefense, UIExperience, UIMaxLevel,
		UIFooterNav, UIFooterNavBack, UIFooterCombat, UIFooterLogin, UIFooterGameOver, UIFooterSmith,
		UITownTitle, UITownSubtitle, UIForestTitle, UIForestSubtitle, UISmithTitle, UISmithSubtitle,
		UIDragonTitle, UIDragonSubtitle, UITavernTitle, UITavernSubtitle, UIChapelTitle, UIChapelSubtitle,
		UIGuildTitle, UIGuildSubtitle, UIGameOverTitle, UIGameOverSubtitle, UILoginTitle, UILoginSubtitle,
		UITownMenuForest, UITownMenuSmith, UITownMenuTavern, UITownMenuChapel, UITownMenuGuild, UITownMenuDragon,
		UITownMenuBankDeposit, UITownMenuBankWithdraw, UITownMenuQuit, UITownDragonReq,
		UIForestMenuLookForFight, UIForestMenuReturnTown, UIForestNoFightsLeft,
		UISmithTabWeapons, UISmithTabArmors, UISmithTabPotions, UISmithAlreadyEquipped, UISmithNotEnoughGold, UISmithNoItemToSell,
		UIDragonAlreadySlain, UITavernFlirt, UITavernTalkRedKnight, UITavernReturnTown, UITavernFlirtLimit, UITavernFlirtSuccess,
		UITavernRedKnightMsg, UIChapelFullHP, UIChapelNotEnoughGold, UIChapelHealSuccess, UIGuildMaxLevelReached, UIGameOverRespawnInfo,
		UILoginUsernameLabel, UILoginPasswordLabel, UILoginErrInvalidPass, UILoginErrUserExists,
	}

	for _, key := range uiKeys {
		text, ok := UITextsPTBR[key]
		if !ok {
			t.Errorf("chave de interface %q sem tradução em UITextsPTBR", key)
			continue
		}
		if text == "" {
			t.Errorf("chave de interface %q com tradução vazia", key)
		}
	}

	if text := GetUIText(UIMaxLevel); text != "(máx)" {
		t.Errorf("expected '(máx)', got %q", text)
	}

	if unknown := GetUIText(UIKey("unknown_ui")); unknown != "unknown_ui" {
		t.Errorf("expected fallback to key for unknown UI text, got %q", unknown)
	}
}

func TestGetMessageAndDescriptions(t *testing.T) {
	if affix := GetAffixName(AffixFerocious); affix != "Feroz" {
		t.Errorf("expected 'Feroz', got %q", affix)
	}

	if desc := GetItemDescription(WeaponStick); desc == "" {
		t.Errorf("expected non-empty description for WeaponStick")
	}

	msg := GetMessage(MsgBankDepositSuccess, 100)
	expectedMsg := "Você depositou 100 moedas de ouro no cofre com segurança!"
	if msg != expectedMsg {
		t.Errorf("expected %q, got %q", expectedMsg, msg)
	}

	if fallback := GetMessage(MessageKey("unknown_key")); fallback != "unknown_key" {
		t.Errorf("expected fallback 'unknown_key', got %q", fallback)
	}
}
