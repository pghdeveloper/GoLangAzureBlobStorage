package lib

type Containers struct {
	//ContainerIds []string `json:"containerIds"`
	ContainerIds []string `validate:"containerIds,required"`
}

type InMemoryFile struct {
	FileName string
	Content  []byte
}