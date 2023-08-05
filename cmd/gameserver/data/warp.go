package data

type Warp struct {
	Id       int
	World    int
	Location []struct {
		X int
		Y int
	}
}

type WarpData struct {
	Warp []struct {
		World int
		Warps []Warp
	}
}

func (wp *WarpData) FindWarp(world byte, id byte) *Warp {
	for _, v := range wp.Warp {
		if v.World != int(world) {
			continue
		}

		for _, v1 := range v.Warps {
			if v1.Id != int(id) {
				continue
			}

			return &v1
		}

		break
	}

	return nil
}
