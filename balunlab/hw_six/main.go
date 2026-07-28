package main

import (
	"bytes"
	"fmt"
	"unsafe"
)

type Option func(*GamePerson)

const (
	UnknownGamePersonType = iota
	BuilderGamePersonType
	BlacksmithGamePersonType
	WarriorGamePersonType
)

func WithName(name string) Option {
	return func(person *GamePerson) {
		copy(person.name[:], []byte(name))
	}
}

func WithCoordinates(x, y, z int32) Option {
	return func(person *GamePerson) {
		person.x = x
		person.y = y
		person.z = z
	}
}

func WithGold(gold uint32) Option {
	return func(person *GamePerson) {
		person.gold = gold
	}
}

func WithMana(mana uint16) Option {
	return func(person *GamePerson) {
		person.stats1 &^= 0b1111111111000000
		person.stats1 |= (mana & 0b1111111111) << 6
	}
}

func WithHealth(health uint16) Option {
	return func(person *GamePerson) {
		person.stats2 &^= 0b1111111111000000
		person.stats2 |= (health & 0b1111111111) << 6
	}
}

func WithRespect(respect uint16) Option {
	return func(person *GamePerson) {
		person.stats3 &^= 0b1111
		person.stats3 |= respect & 0b1111
	}
}

func WithStrength(strength uint16) Option {
	return func(person *GamePerson) {
		person.stats1 &^= 0b0000000000111100
		person.stats1 |= (strength & 0b1111) << 2
	}
}

func WithExperience(experience uint16) Option {
	return func(person *GamePerson) {
		person.stats2 &^= 0b1111
		person.stats2 |= experience & 0b1111
	}
}

func WithLevel(level uint16) Option {
	return func(person *GamePerson) {
		person.stats3 &^= 0b11110000
		person.stats3 |= (level & 0b1111) << 4
	}
}

func WithHouse() Option {
	return func(person *GamePerson) {
		person.stats3 |= 1 << 15
	}
}

func WithGun() Option {
	return func(person *GamePerson) {
		person.stats3 |= 1 << 14
	}
}

func WithFamily() Option {
	return func(person *GamePerson) {
		person.stats3 |= 1 << 13
	}
}

func WithType(personType int) Option {
	return func(person *GamePerson) {
		person.stats1 &^= 0b11

		if personType == BuilderGamePersonType {
			person.stats1 |= 1 << 0
		}

		if personType == BlacksmithGamePersonType {
			person.stats1 |= 1 << 1
		}

		if personType == WarriorGamePersonType {
			person.stats1 |= 1 << 0
			person.stats1 |= 1 << 1
		}
	}
}

type GamePerson struct {
	x      int32    // 4 байта
	y      int32    // 4 байта
	z      int32    // 4 байта
	gold   uint32   // 4 байта
	name   [42]byte // 42 байта
	stats1 uint16   // 2 байта
	stats2 uint16   // 2 байта
	stats3 uint16   // 2 байта
}

// stats1
// 15........6 | 5..2     | 1..0
// mana          strength   type

// stats2
// 15........6 | 5..3     | 3..0
// health        unused     experience

// stats3
// 15........13 | 12..8   | 7..4  | 3..0
// bools          unused     level   respect

func NewGamePerson(options ...Option) *GamePerson {
	p := &GamePerson{}

	for _, opt := range options {
		opt(p)
	}

	return p
}

func (p *GamePerson) Name() string {
	return string(bytes.TrimRight(p.name[:], "\x00"))
}

func (p *GamePerson) X() int32 {
	return p.x
}

func (p *GamePerson) Y() int32 {
	return p.y
}

func (p *GamePerson) Z() int32 {
	return p.z
}

func (p *GamePerson) Gold() uint32 {
	return p.gold
}

func (p *GamePerson) Mana() int {
	return int((p.stats1 >> 6) & 0b1111111111)
}

func (p *GamePerson) Health() int {
	return int((p.stats2 >> 6) & 0b1111111111)
}

func (p *GamePerson) Respect() int {
	return int(p.stats3 & 0b1111)
}

func (p *GamePerson) Strength() int {
	return int((p.stats1 >> 2) & 0b1111)
}

func (p *GamePerson) Experience() int {
	return int(p.stats2 & 0b1111)
}

func (p *GamePerson) Level() int {
	return int((p.stats3 >> 4) & 0b1111)
}

func (p *GamePerson) HasHouse() bool {
	return (p.stats3 & (1 << 15)) != 0
}

func (p *GamePerson) HasGun() bool {
	return (p.stats3 & (1 << 14)) != 0
}

func (p *GamePerson) HasFamily() bool {
	return (p.stats3 & (1 << 13)) != 0
}

func (p *GamePerson) Type() int {
	firstBit := p.stats1 & 1
	secondBit := (p.stats1 >> 1) & 1

	if firstBit == 1 && secondBit == 0 {
		return BuilderGamePersonType
	}

	if firstBit == 0 && secondBit == 1 {
		return BlacksmithGamePersonType
	}

	if firstBit == 1 && secondBit == 1 {
		return WarriorGamePersonType
	}

	return UnknownGamePersonType
}

func main() {
	fmt.Println(unsafe.Sizeof(GamePerson{}))
	// person := NewGamePerson(WithCoordinates(1, 2, 3))
}
