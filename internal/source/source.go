package source

type Source string

const (
	API   Source = "api"   // External Traefik API config
	Local Source = "local" // Local config
	Agent Source = "agent" // Agent config
)

func (s Source) Valid() bool {
	switch s {
	case API, Local, Agent:
		return true
	default:
		return false
	}
}
