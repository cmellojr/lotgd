package bestiary_test

import (
	"math/rand"
	"testing"

	"lotgd/internal/bestiary"
	"lotgd/internal/i18n"
)

func TestMonsterGenerator_TierBounds(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	gen := bestiary.NewMonsterGenerator(rng)

	// Gera 20 monstros do Tier 1 e valida bounds
	for i := 0; i < 20; i++ {
		m := gen.GenerateByTier(1)
		if m.Tier != 1 {
			t.Fatalf("esperado monstro de Tier 1, obtido %d", m.Tier)
		}
		if m.Health <= 0 || m.Attack <= 0 || m.Defense < 0 {
			t.Fatalf("atributos inválidos para o monstro: %+v", m)
		}
	}
}

func TestMonsterGenerator_PlayerScaling(t *testing.T) {
	rng := rand.New(rand.NewSource(123))
	gen := bestiary.NewMonsterGenerator(rng)

	// Jogador de nível alto deve receber monstros de Tier 4
	m := gen.GenerateForPlayer(10)
	if m.Tier != 4 {
		t.Fatalf("jogador de nível 10 deve enfrentar monstros de Tier 4, obtido Tier %d", m.Tier)
	}
}

func TestGenerateDragonOfDay_Deterministic(t *testing.T) {
	day1 := "2026-08-25"
	day2 := "2026-08-26"

	d1a := bestiary.GenerateDragonOfDay(day1)
	d1b := bestiary.GenerateDragonOfDay(day1)
	d2 := bestiary.GenerateDragonOfDay(day2)

	// Mesma data deve gerar exatamente o mesmo dragão
	if d1a.Name != d1b.Name || d1a.Health != d1b.Health || d1a.Attack != d1b.Attack || d1a.Defense != d1b.Defense {
		t.Fatalf("o Dragão do Dia deve ser 100%% determinístico para a mesma data. d1a=%+v, d1b=%+v", d1a, d1b)
	}

	if d1a.ID != i18n.MonsterDragon || !d1a.IsDragon {
		t.Fatalf("o Dragão deve ter ID canônica e IsDragon=true")
	}

	// Datas diferentes devem gerar variações
	if d1a.Health == d2.Health && d1a.Attack == d2.Attack && d1a.Name == d2.Name {
		t.Logf("Aviso: stats coincidentemente iguais entre dias diferentes (raro mas possível)")
	}
}
