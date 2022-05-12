package pkg

type Tenant struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type RequestPayload struct {
	Customer      string   `json:"customer"`
	Organization  Tenant   `json:"organization"`
	Folders       []Tenant `json:"folders,omitempty"`
	ResourceGroup *Tenant  `json:"resource_group,omitempty"`
}
