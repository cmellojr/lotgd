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
