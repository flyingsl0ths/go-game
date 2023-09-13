package main

type Grid[T any] struct {
	objects []T
	columns int
	rows    int
}

func (grid *Grid[T]) get(x int, y int) *T {
	if y > grid.columns || x > grid.rows {
		return nil
	}

	return &grid.objects[y+x*len(grid.objects)]
}
