package settings

func (c *client) LastSource() string {
	return c.Dir.Source.Last
}
func (c *client) LastTarget() string {
	return c.Dir.Target.Last
}
func (c *client) LastSwapped() string {
	return c.Dir.Output.Last
}
func (c *client) Source() string {
	return c.Dir.Source.Folder
}
func (c *client) Target() string {
	return c.Dir.Target.Folder
}
func (c *client) Swapped() string {
	return c.Dir.Output.Folder
}
