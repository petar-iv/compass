package solution

import (
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
)

type converter struct{}

func NewConverter() *converter {
	return &converter{}
}

func (c *converter) ToGraphQL(in *model.Solution) *graphql.Solution {
	if in == nil {
		return nil
	}

	return &graphql.Solution{
		ID:          in.ID,
		Name:        in.Name,
		Version:     in.Version,
		Description: in.Description,
	}
}

func (c *converter) MultipleToGraphQL(in []*model.Solution) []*graphql.Solution {
	var solutions []*graphql.Solution
	for _, r := range in {
		if r == nil {
			continue
		}

		solutions = append(solutions, c.ToGraphQL(r))
	}

	return solutions
}

func (c *converter) InputFromGraphQL(in graphql.SolutionInput) model.SolutionInput {
	var labels map[string]interface{}
	if in.Labels != nil {
		labels = in.Labels
	}

	return model.SolutionInput{
		Name:            in.Name,
		Description:     in.Description,
		Version: in.Version,
		Labels:          labels,
	}
}

func (c *converter) MultipleFromEntities(entities SolutionCollection) []*model.Solution {
	var items []*model.Solution
	for _, ent := range entities {
		model := ent.ToModel()

		items = append(items, model)
	}
	return items
}
