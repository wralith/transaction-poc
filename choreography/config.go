package choreography

const ProtocolAddr = "http://localhost"
const FirstStepAppPort = ":3000"
const SecondStepAppPort = ":3001"
const ThirdStepAppPort = ":3002"

type FirstStepAppURL struct {
	Port string
}

func GetFirstStepAppConfig() FirstStepAppURL {
	return FirstStepAppURL{
		Port: FirstStepAppPort,
	}
}

type SecondStepAppConfig struct {
	Port            string
	FirstStepAppURL string
}

func GetSecondStepAppConfig() SecondStepAppConfig {
	return SecondStepAppConfig{
		Port:            SecondStepAppPort,
		FirstStepAppURL: ProtocolAddr + FirstStepAppPort,
	}
}

type ThirdStepAppConfig struct {
	Port             string
	SecondStepAppURL string
}

func GetThirdStepAppConfig() ThirdStepAppConfig {
	return ThirdStepAppConfig{
		Port:             ThirdStepAppPort,
		SecondStepAppURL: ProtocolAddr + SecondStepAppPort,
	}
}
