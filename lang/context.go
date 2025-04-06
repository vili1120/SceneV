package lang

type Context struct {
  DisplayName string
  Parent *Context
  ParentEntryPos *Position
}
