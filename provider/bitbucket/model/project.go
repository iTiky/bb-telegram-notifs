package model

type Project struct {
	ID          int64  `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
	Type        string `json:"type"`
	Links       Links  `json:"links"`
}

func (p Project) SelfLink() *Link {
	l, found := p.Links["self"]
	if !found {
		return nil
	}

	if len(l) > 0 {
		return &l[0]
	}

	return nil
}

type Projects struct {
	Page
	Values []Project `json:"values"`
}
