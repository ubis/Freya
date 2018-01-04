package inventory

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
	12: "left_ear",
	13: "right_ear",
	14: "left_bracelet",
	15: "right_bracelet",
	16: "finger3",
	17: "finger4",
	18: "belt",
	19: "extended_pet",
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
