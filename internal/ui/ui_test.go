package ui

import (
	"fmt"
	"strings"
	"testing"

	"lotgd/internal/engine"
	"lotgd/internal/i18n"
)

func TestRenderStatusBar_NilPlayer(t *testing.T) {
	out := RenderStatusBar(nil, 80)
	if out != "" {
		t.Errorf("expected empty string for nil player, got %q", out)
	}
}

func TestRenderStatusBar_ValidPlayer(t *testing.T) {
	p := &engine.Player{
		Username:     "SirGalahad",
		Level:        3,
		Health:       45,
		MaxHealth:    60,
		Gold:         250,
		BankGold:     1000,
		ForestFights: 12,
		Weapon:       engine.Item{ID: i18n.WeaponShortSword},
		Armor:        engine.Item{ID: i18n.ArmorLeather},
		PotionsCount: 2,
	}

	out := RenderStatusBar(p, 80)
	if out == "" {
		t.Fatalf("expected non-empty status bar output")
	}

	// Verify key information appears in rendered output
	expectedSubstrings := []string{
		"SirGalahad",
		"3",
		"45",
		"60",
		"250",
		"1000",
		"12",
	}

	for _, substr := range expectedSubstrings {
		if !strings.Contains(out, substr) {
			t.Errorf("expected status bar to contain %q, but it did not.\nOutput:\n%s", substr, out)
		}
	}
}

func TestRenderHPBar_Boundaries(t *testing.T) {
	tests := []struct {
		name        string
		ratio       float64
		totalBlocks int
	}{
		{name: "Negative ratio (clamped to 0)", ratio: -0.5, totalBlocks: 10},
		{name: "Zero ratio (dead)", ratio: 0.0, totalBlocks: 10},
		{name: "Half ratio", ratio: 0.5, totalBlocks: 10},
		{name: "Full ratio", ratio: 1.0, totalBlocks: 10},
		{name: "Over 1.0 ratio (clamped to 1)", ratio: 1.5, totalBlocks: 10},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bar := renderHPBar(tc.ratio, tc.totalBlocks)
			if !strings.HasPrefix(bar, "[") || !strings.HasSuffix(bar, "]") {
				t.Errorf("expected bar to be enclosed in brackets, got %q", bar)
			}
		})
	}
}

func TestStyles_NoPanic(t *testing.T) {
	// Verify that basic styles render without panic
	sample := "Test Text"
	if TitleStyle.Render(sample) == "" {
		t.Errorf("TitleStyle rendered empty string")
	}
	if AppStyle.Render(sample) == "" {
		t.Errorf("AppStyle rendered empty string")
	}
	if ErrorNoticeStyle.Render(sample) == "" {
		t.Errorf("ErrorNoticeStyle rendered empty string")
	}
	if SuccessNoticeStyle.Render(sample) == "" {
		t.Errorf("SuccessNoticeStyle rendered empty string")
	}
	if SelectedMenuItemStyle.Render(sample) == "" {
		t.Errorf("SelectedMenuItemStyle rendered empty string")
	}
}

func TestRenderStatusBar_ShowsCombatStatsAndXP(t *testing.T) {
	p := &engine.Player{
		Username:    "Hero",
		Level:       1,
		Experience:  40,
		Health:      20,
		MaxHealth:   20,
		BaseAttack:  5,
		BaseDefense: 2,
		Weapon:      engine.WeaponsCatalog[0],
		Armor:       engine.ArmorsCatalog[0],
	}
	out := RenderStatusBar(p, 120)

	for _, want := range []string{
		"ATK:", fmt.Sprintf("%d", p.TotalAttack()),
		"DEF:", fmt.Sprintf("%d", p.TotalDefense()),
		"XP:", "40/",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("barra de estado não mostra %q:\n%s", want, out)
		}
	}
}

func TestRenderStatusBar_MaxLevelUsesLocalizedTag(t *testing.T) {
	p := &engine.Player{
		Username:    "Hero",
		Level:       engine.MaxLevel,
		Experience:  22000,
		Health:      20,
		MaxHealth:   20,
		BaseAttack:  5,
		BaseDefense: 2,
		Weapon:      engine.WeaponsCatalog[0],
		Armor:       engine.ArmorsCatalog[0],
	}
	out := RenderStatusBar(p, 120)

	want := fmt.Sprintf("%d %s", p.Experience, i18n.GetUIText(i18n.UIMaxLevel))
	if !strings.Contains(out, want) {
		t.Errorf("no nível máximo a barra de estado não mostra %q:\n%s", want, out)
	}
}
