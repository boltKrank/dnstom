package dnswire

// Since DNS was built back in the days when memory and bandwidth was at a premium, everything was represented in the smallest amount of bit use possible.
// I can be a nice flex (and also quite useful to remember what every number is) - but in the intrests of good code we can replace said numbers with const variable names.
// Since a lot of DNS stuff is fixed, we're pretty safe in how many we use provided we don't make the mistake I did many years back
// (use the same number for two diferent consts used in the same scope which just happened to pass QA and caused minor issues in a banking system 18 months later).

const (
	// DNS record types
	TypeA uint16 = 1

	// DNS class
	ClassIN uint16 = 1
)
