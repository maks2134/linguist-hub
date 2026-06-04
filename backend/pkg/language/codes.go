package language

type UILocale string

type LangCode string

type Role string

const (
	RoleLearning Role = "learning"
	RoleNative   Role = "native"
)

const DefaultNative LangCode = "english"
