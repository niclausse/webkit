package mode

type Mode int32

const (
	DevelopMode Mode = iota
	TestMode
	ProductionMode
)

func FromString(modeStr string) Mode {
	if modeStr == "develop" || modeStr == "dev" {
		return DevelopMode
	}

	if modeStr == "test" || modeStr == "qa" {
		return TestMode
	}

	if modeStr == "production" || modeStr == "release" {
		return ProductionMode
	}

	return DevelopMode
}
