package main

func CountIslands(grid [][]int) int {

	var count int = 0
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			if grid[i][j] == 1 {
				count++
				checkCount(grid, i, j)
			}
		}
	}
	return count
}

func checkCount(grid [][]int, i int, j int) {
	if i < 0 || j < 0 || i > len(grid)-1 || j > len(grid[0])-1 {
		return
	}
	if grid[i][j] != 1 {
		return
	}
	grid[i][j] = 0
	checkCount(grid, i-1, j)
	checkCount(grid, i+1, j)
	checkCount(grid, i, j-1)
	checkCount(grid, i, j+1)
}
