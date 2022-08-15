package model

type DestinationInput struct {
	Name                    string `json:"Name"`
	Type                    string `json:"Type"`
	URL                     string `json:"URL"`
	Authentication          string `json:"Authentication"`
	XFSystemName            string `json:"XFSystemName"`
	CommunicationScenarioId string `json:"communicationScenarioId"`
	ProductName             string `json:"product.name"`
}
