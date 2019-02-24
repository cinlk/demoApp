package httpModel


type HttpSearchTopWords struct {

	Name string `json:"name"`
}

type HttpSearchWord struct {
	Type string `json:"type"`
	Words []string `json:"words"`
}