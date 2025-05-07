package structs

// cmd struct
type CMD struct {
	Name    string
	Help    string
	Execute func(t *Target)
}
