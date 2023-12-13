package inventory

type EquipmentType int

const (
	Helmet = iota
	Suit
	Gloves
	Boots
	RightHand
	LeftHand
	Epaulet
	Amulet
	Ring1
	Ring2
	Vehicle
	Pet
	Unknown
	LeftEarring
	RightEarring
	LeftBracelet
	RightBracelet
	Ring3
	Ring4
	Belt

	Invalid = -1
)

var eqTypes = map[int]string{
	0:  "helmet",
	1:  "suit",
	2:  "gloves",
	3:  "boots",
	4:  "right_hand",
	5:  "left_hand",
	6:  "epaulet",
	7:  "amulet",
	8:  "finger1",
	9:  "finger2",
	10: "vehicle",
	11: "pet",
	12: "unk",
	13: "left_ear",
	14: "right_ear",
	15: "left_bracelet",
	16: "right_bracelet",
	17: "finger3",
	18: "finger4",
	19: "belt",
}

// Returns equipment slot id by given name
func MapEquipment(name string) uint16 {
	for key, value := range eqTypes {
		if value == name {
			return uint16(key)
		}
	}

	return 65535
}
