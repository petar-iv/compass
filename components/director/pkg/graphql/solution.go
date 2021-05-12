package graphql

type Solution struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Version     string  `json:"version"`
}

type SolutionExt struct {
	Solution
	Labels      Labels  `json:"labels"`
}

type SolutionPageExt struct {
	SolutionPage
	Data []*SolutionExt `json:"data"`
}
