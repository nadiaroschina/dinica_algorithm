package main

import "fmt"

func min(a, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}

// calculating level graph using breadth-first search
// e = (v, u) is in level graph <=> level[v] == level[u]+1
func bfs(n, s, t int, graph [][]int, flow [][]int, capacity [][]int, level *[]int) bool {
	// initializing level as in standard bfs; setting s as a root
	*level = make([]int, n)
	for i := 0; i < n; i++ {
		(*level)[i] = -1
	}
	q := make([]int, 0)
	q = append(q, s)
	(*level)[s] = 0
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		for _, v := range graph[u] {
			// checking if the edge is in residual graph
			if capacity[u][v]-flow[u][v] > 0 && (*level)[v] == -1 {
				(*level)[v] = (*level)[u] + 1
				q = append(q, v)
			}
		}
	}
	// we are finished only if t is unreachable from s in residual graph
	return (*level)[t] != -1
}

// creating blocking flow in level graph by finding augmenting paths using depth-first search
// we consider current path to be s -> u, starting with u = s and finishing with u = t after
func dfs(n, u, s, t int, currMin int, graph [][]int, capacity [][]int, flow *[][]int, level []int, it []int) int {
	// currMin is the lowest value of the residual edge in current path
	if u == t {
		return currMin
	}
	for it[u] < len(graph[u]) {
		// no need to consider paths with edges u-w, where w index is less than it[u] if graph[u],
		// as we already blocked all the paths from such w in previous dfs
		v := graph[u][it[u]]
		// checking if (u, v) is in both level graph and residual graph
		if level[v] == level[u]+1 && capacity[u][v]-(*flow)[u][v] > 0 {
			// calculating the lowest value of the residual edge in the rest of current path
			restMin := dfs(n, v, s, t, min(currMin, capacity[u][v]-(*flow)[u][v]), graph, capacity, flow, level, it)
			// updating the flow
			if restMin > 0 {
				(*flow)[u][v] += restMin
				(*flow)[v][u] -= restMin
				return restMin
			}
		}
		// we don't need to consider v for blocking flow anymore
		it[u]++
	}
	return 0
}

func main() {

	const INF = 1e9 + 7

	// vertex and edges count
	var n, m int
	// starting and terminal flow vertexes
	var s, t int
	fmt.Scan(&n, &m, &s, &t)

	// initializing graph as adjacency lists
	graph := make([][]int, n)
	// edges capacity and current flow as tables
	capacity := make([][]int, n)
	flow := make([][]int, n)
	for i := 0; i < n; i++ {
		capacity[i] = make([]int, n)
		flow[i] = make([]int, n)
	}
	// vertex levels in level graph
	level := make([]int, n)

	// building residual graph and initializing edges capacities
	for i := 0; i < m; i++ {
		var u, v int
		var c int
		fmt.Scan(&u, &v, &c)
		graph[u] = append(graph[u], v)
		graph[v] = append(graph[v], u)
		capacity[u][v] = c
	}

	// while t is reachable from s in residual graph
	for bfs(n, s, t, graph, flow, capacity, &level) {
		// vertex iterator for faster finding of blocking flow
		it := make([]int, n)
		for dfs(n, s, s, t, INF, graph, capacity, &flow, level, it) > 0 {
		}
	}

	// calculating resulting flow
	res := 0
	for i := 0; i < n; i++ {
		res += flow[i][t]
	}
	fmt.Println(res)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if flow[i][j] > 0 {
				fmt.Print(flow[i][j], " ")
			} else {
				fmt.Print(0, " ")
			}
		}
		fmt.Println()
	}
}
