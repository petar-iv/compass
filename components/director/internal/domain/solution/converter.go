package solution

import (
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/pkg/errors"
)

type AppTemplateConverter interface {
	ApplicationFromTemplateInputFromGraphQL(in graphql.ApplicationFromTemplateInput) model.ApplicationFromTemplateInput
}

type BundleConverter interface {
	MultipleCreateInputFromGraphQL(in []*graphql.BundleCreateInput) ([]*model.BundleCreateInput, error)
}

type converter struct {
	appTemplateConv AppTemplateConverter
	bundleConv      BundleConverter
}

func NewConverter(appTemplateConv AppTemplateConverter, bundleConv BundleConverter) *converter {
	return &converter{
		appTemplateConv: appTemplateConv,
		bundleConv:      bundleConv,
	}
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

func (c *converter) InputFromGraphQL(in graphql.SolutionInput) (model.SolutionInput, error) {
	var labels map[string]interface{}
	if in.Labels != nil {
		labels = in.Labels
	}

	deps := make([]model.SolutionDependencyInput, 0)
	for _, v := range in.Dependencies {
		// TODO check if impossible
		if v.Application == nil {
			continue
		}

		appTplIn := c.appTemplateConv.ApplicationFromTemplateInputFromGraphQL(*v.Application)
		bunddlesIn, err := c.bundleConv.MultipleCreateInputFromGraphQL(v.Bundles)
		if err != nil {
			return model.SolutionInput{}, errors.Wrapf(err, "while converting bundles for solution with name %s", in.Name)
		}
		deps = append(deps, model.SolutionDependencyInput{
			Application: appTplIn,
			Bundles:     bunddlesIn,
		})
	}
	return model.SolutionInput{
		Name:         in.Name,
		Description:  in.Description,
		Version:      in.Version,
		Labels:       labels,
		Dependencies: deps,
	}, nil
}

func (c *converter) MultipleFromEntities(entities SolutionCollection) []*model.Solution {
	var items []*model.Solution
	for _, ent := range entities {
		model := ent.ToModel()

		items = append(items, model)
	}
	return items
}
