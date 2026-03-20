package ports

type Spring interface {
	Update(pos float64, vel float64, target float64) (nextPos float64, nextVel float64)
}
