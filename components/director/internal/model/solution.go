package model

import (
	"github.com/kyma-incubator/compass/components/director/pkg/pagination"
)

type Solution struct {
	ID          string
	Name        string
	Tenant      string
	Version     string
	Description *string
}

type SolutionInput struct {
	Name         string
	Version      string
	Description  *string
	Dependencies []SolutionDependencyInput
	Labels       map[string]interface{}
}

type SolutionDependencyInput struct {
	Application ApplicationFromTemplateInput
	Bundles     []*BundleCreateInput
}

func (i *SolutionInput) ToSolution(id string, tenant string) *Solution {
	if i == nil {
		return nil
	}

	return &Solution{
		ID:          id,
		Name:        i.Name,
		Tenant:      tenant,
		Version:     i.Version,
		Description: i.Description,
	}
}

type SolutionPage struct {
	Data       []*Solution
	PageInfo   *pagination.Page
	TotalCount int
}
