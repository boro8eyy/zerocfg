package zerocfg

const noSource = "default"

// node represents a single configuration option, including its name, description, aliases, value, and metadata.
type node struct {
	Name        string
	Description string
	Aliases     []string
	Value       Value
	setSource   string
	isSecret    bool
	isRequired  bool
	caller      string
}

func (n *node) pathName() string {
	if n.caller == "" {
		return n.Name
	}

	return n.caller + ":" + n.Name
}

func (n *node) source() string {
	if n.setSource == "" {
		return noSource
	}

	return n.setSource
}

// Value is the interface implemented by all configuration option types in zerocfg.
//
// Requirements:
//   - Must support setting its value from a string:
//     Set(string) error
//     The string is produced by zerocfg's ToString conversion.
//   - Must report its type name for identification and documentation:
//     Type() string
type Value interface {
	Set(string) error
	Type() string
}

// OptNode is a function that modifies a node during option registration.
// It is used to apply additional behaviors such as aliases, secret marking, grouping, or required flags.
//
// Example:
//
//	Int("db.port", 5432, "database port", Alias("p"), Required())
type OptNode func(*node)

// Alias returns an OptNode that adds an alias to a configuration option.
// Aliases allow options to be referenced by alternative names (e.g., for CLI flags).
//
// Example:
//
//	port := Int("db.port", 5432, "database port", Alias("p"))
func Alias(alias string) OptNode {
	return func(n *node) {
		n.Aliases = append(n.Aliases, alias)
	}
}

// Secret returns an OptNode that marks a configuration option as secret.
// Secret options are masked in rendered output (e.g., Show) to avoid leaking sensitive values.
//
// Example:
//
//	password := Str("db.password", "", "database password", Secret())
func Secret() OptNode {
	return func(n *node) {
		n.isSecret = true
	}
}

// Group returns an OptNode that applies a Grp to a configuration option.
// This sets the option's name prefix and applies all group modifiers.
//
// Example:
//
//	g := NewGroup("db")
//	host := Str("host", "localhost", "db host", Group(g)) // becomes "db.host"
func Group(g *Grp) OptNode {
	return func(n *node) {
		n.Name = g.key(n.Name)
		g.applyOpts(n)
	}
}

// Required returns an OptNode that marks a configuration option as required.
// Required options must be set by a configuration source or an error will be returned by Parse.
//
// Example:
//
//	user := Str("db.user", "", "database user", Required())
func Required() OptNode {
	return func(n *node) {
		n.isRequired = true
	}
}
