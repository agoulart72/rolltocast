package app

import (
	"math/rand"
)

type Spell struct {
	Level int
	DC    int
	Burn  int
}

var SpellCost = [...]Spell{
	{Level: 0, DC: 0, Burn: 0},   // 0 0
	{Level: 1, DC: 8, Burn: 8},   // 12 2  (10)
	{Level: 2, DC: 10, Burn: 10}, // 13 3  (10)
	{Level: 3, DC: 12, Burn: 12}, // 15 5  (10)
	{Level: 4, DC: 15, Burn: 15}, // 16 6  (10)
	{Level: 5, DC: 17, Burn: 17}, // 17 7  (10)
	{Level: 6, DC: 19, Burn: 19}, // 19 9  (11)
	{Level: 7, DC: 22, Burn: 22}, // 20 10 (12)
	{Level: 8, DC: 24, Burn: 24}, // 21 11 (13)
	{Level: 9, DC: 26, Burn: 26}, // 23 13 (13)
}

type Sorcerer struct {
	Level            int
	Bonus            int
	MaxSpellRank     int
	CurrentSpellRank int
	Burn             int
	CurrentHitDie    int
}

var SorcererByLevel = [...]Sorcerer{
	{Level: 1, Bonus: 5, MaxSpellRank: 1},
	{Level: 2, Bonus: 5, MaxSpellRank: 1},
	{Level: 3, Bonus: 5, MaxSpellRank: 2},
	{Level: 4, Bonus: 6, MaxSpellRank: 2},
	{Level: 5, Bonus: 7, MaxSpellRank: 3},
	{Level: 6, Bonus: 7, MaxSpellRank: 3},
	{Level: 7, Bonus: 7, MaxSpellRank: 4},
	{Level: 8, Bonus: 8, MaxSpellRank: 4},
	{Level: 9, Bonus: 9, MaxSpellRank: 5},
	{Level: 10, Bonus: 9, MaxSpellRank: 5},

	{Level: 11, Bonus: 9, MaxSpellRank: 6},
	{Level: 12, Bonus: 9, MaxSpellRank: 6},
	{Level: 13, Bonus: 10, MaxSpellRank: 7},
	{Level: 14, Bonus: 10, MaxSpellRank: 7},
	{Level: 15, Bonus: 10, MaxSpellRank: 8},
	{Level: 16, Bonus: 10, MaxSpellRank: 8},
	{Level: 17, Bonus: 11, MaxSpellRank: 9},
	{Level: 18, Bonus: 11, MaxSpellRank: 9},
	{Level: 19, Bonus: 11, MaxSpellRank: 9},
	{Level: 20, Bonus: 11, MaxSpellRank: 9},
}

func (s *Sorcerer) LongRest() {
	s.Bonus = SorcererByLevel[s.Level-1].Bonus
	s.MaxSpellRank = SorcererByLevel[s.Level-1].MaxSpellRank
	s.CurrentSpellRank = SorcererByLevel[s.Level-1].MaxSpellRank
	s.Burn = 0
	s.CurrentHitDie = s.Level
}

func (s *Sorcerer) ShortRest() {

	burnedSlots := SorcererByLevel[s.Level-1].MaxSpellRank - s.CurrentSpellRank

	if burnedSlots > 0 && s.CurrentHitDie > s.Level/2 {
		s.CurrentSpellRank += 1
		s.CurrentHitDie -= s.CurrentSpellRank
	}

}

func (s *Sorcerer) Cast(rank int) (success bool, backlash bool) {

	success = false
	backlash = false

	if s.CurrentSpellRank <= 0 {
		return false, false
	}

	spellCheckRoll := rand.Intn(20) + 1

	success = (spellCheckRoll + s.Bonus) >= SpellCost[rank].DC

	if spellCheckRoll < 20 {
		backlash = spellCheckRoll <= s.Burn
	}

	if success {
		if SpellCost[rank].Burn-s.Level > 0 {
			s.Burn += SpellCost[rank].Burn - s.Level
		} else {
			s.Burn += 1
		}

	}
	if backlash {
		s.Burn = 0
		s.CurrentSpellRank = s.CurrentSpellRank - 1
	}

	return success, backlash
}
