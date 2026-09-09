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

// TestUITextsCoverAllKeys garante que toda chave de interface declarada em
// types.go tem tradução registrada, conforme a verificação prevista no
// ADR-0005.
func TestUITextsCoverAllKeys(t *testing.T) {
	for _, key := range []UIKey{UIMaxLevel} {
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
