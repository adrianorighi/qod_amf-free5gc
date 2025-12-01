package models

// Device contém as informações de identificação do dispositivo.
type Device struct {
	NetworkAccessIdentifier string `json:"networkAccessIdentifier"`
}

// ApplicationServer contém os endereços do servidor de aplicação.
type ApplicationServer struct {
	IPv4Address string `json:"ipv4Address"`
	IPv6Address string `json:"ipv6Address"`
}

// DevicePorts especifica os ranges de portas do dispositivo.
type DevicePorts struct {
	Ranges []PortRange `json:"ranges"`
}

// PortRange define um intervalo de portas, com um início e um fim.
type PortRange struct {
	From int `json:"from"`
	To   int `json:"to"`
}

// Represents the QoD Management Object.
type QoDManagement struct {
	Device            Device            `json:"device"`
	ApplicationServer ApplicationServer `json:"applicationServer"`
	DevicePorts       DevicePorts       `json:"devicePorts"`
	Duration          int               `json:"duration"`
	QoSProfile        string            `json:"qosProfile" yaml:"qosProfile" bson:"qosProfile,omitempty"`
}
