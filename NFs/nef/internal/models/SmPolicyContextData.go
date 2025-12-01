package models

// SmPolicyContextData é o corpo da requisição POST para /sm-policies (CreateSMPolicy).
// Ela contém todas as informações de contexto da PDU Session para o PCF.
type SmPolicyContextData struct {
	// 5.2.2.2.1 Mandatory
	// SUPI (Subscription Permanent Identifier) - Mandatory na criação.
	Supi string `json:"supi"`
	// PDU Session ID - Mandatory.
	PduSessionId int32 `json:"pduSessionId"`
	// DNN (Data Network Name) - Mandatory.
	Dnn string `json:"dnn"`
	// S-NSSAI (Single Network Slice Selection Assistance Information) - Mandatory.
	Snssai Snssai `json:"snssai"`
	// Notification URI para updates - Mandatory.
	NotificationUri string `json:"notificationUri"`
	// Tipo de Sessão (e.g., "IPV4", "IPV6", "IPV4V6", "ETHERNET") - Mandatory.
	PduSessionType string `json:"pduSessionType"`
	// Tipo de Acesso (e.g., "3GPP_ACCESS", "NON_3GPP_ACCESS") - Mandatory.
	AccessType AccessType `json:"accessType"`

	// 5.2.2.2.2 Optional
	// PEI (Permanent Equipment Identifier).
	Pei string `json:"pei,omitempty"`
	// Endereço IPv4 alocado.
	Ipv4Address string `json:"ipv4Address,omitempty"`
	// Prefixo IPv6 alocado.
	Ipv6AddressPrefix string `json:"ipv6AddressPrefix,omitempty"`
	// Prefixo IPv6 liberado (apenas para updates, mas incluído na struct).
	RelIpv6AddressPrefix string `json:"relIpv6AddressPrefix,omitempty"`
	// Lista de prefixos IPv6 do UE.
	UeIpv6AddressPrefixes []string `json:"ueIpv6AddressPrefixes,omitempty"`
	// AMBR da Sessão por assinatura.
	SubscribedSessionAmbr *Ambr `json:"subscribedSessionAmbr,omitempty"`
	// Informação de QoS Default por assinatura.
	SubscribedDefaultQos *SubscribedDefaultQos `json:"subscribedDefaultQos,omitempty"`
	// Localização do UE (e.g., TAI, ECGI, NGRAN-ID).
	UeLocation *UserLocation `json:"ueLocation,omitempty"`
	// Tipo de Localização.
	UeLocationTimestamp string `json:"ueLocationTimestamp,omitempty"`
	// V-PLMN ID.
	VplmnId *PlmnId `json:"vplmnId,omitempty"`
	// Identificador da NF Servidora (SMF).
	ServingNfId *ServingNfId `json:"servingNfId,omitempty"`
	// Características de Charging (Off-line, On-line).
	ChargingCharacteristics string `json:"chargingCharacteristics,omitempty"`
	// Status do 3GPP PS Data Off.
	PsDataOffStatus bool `json:"psDataOffStatus,omitempty"`
	// Tipo de RAT (e.g., "NR", "EUTRA").
	RatType string `json:"ratType,omitempty"`
	// Indicador de PDU Session de emergência.
	EmergencyInd bool `json:"emergencyInd,omitempty"`
	// Indicação de PDU Session de plano de usuário.
	UserPlaneOnly bool `json:"userPlaneOnly,omitempty"`
	// Endereços MAC do UE na PDU Session (para Ethernet).
	UeMac string `json:"ueMac,omitempty"`
	// Modo de seleção do DNN (e.g., "ALLOWED", "NOT_ALLOWED").
	DnnSelectionMode string `json:"dnnSelectionMode,omitempty"`
}

// Sub-estruturas Comuns (Common Data Types)

// Snssai representa o Single Network Slice Selection Assistance Information.
type Snssai struct {
	// Slice/Service Type - Mandatory.
	Sst int32 `json:"sst"`
	// Slice Differentiator (Opcional).
	Sd string `json:"sd,omitempty"`
}

// Ambr representa o Aggregate Maximum Bit Rate para Uplink e Downlink.
type Ambr struct {
	// Uplink AMBR (Ex: "2000 Mbps", "500 Kbps") - Mandatory.
	Uplink string `json:"uplink"`
	// Downlink AMBR (Ex: "4000 Mbps", "1 Gbps") - Mandatory.
	Downlink string `json:"downlink"`
}

// SubscribedDefaultQos representa a informação de QoS Default por assinatura.
type SubscribedDefaultQos struct {
	// 5G QoS Identifier (5QI) - Mandatory.
	FiveQi int32 `json:"5qi"`
	// Allocation and Retention Priority (ARP) - Mandatory.
	Arp Arp `json:"arp"`
	// ... outros campos de QoS (e.g., Ppre, Rqe)
}

// Arp representa a Priority and Pre-emption related information.
type Arp struct {
	// Priority Level (1-15) - Mandatory.
	PriorityLevel int32 `json:"priorityLevel"`
	// Pre-emption Capability ("NOT_PREEMPT", "MAY_PREEMPT") - Mandatory.
	PreemptionCap string `json:"preemptionCap"`
	// Pre-emption Vulnerability ("NOT_PREEMPTABLE", "PREEMPTABLE") - Mandatory.
	PreemptionVuner string `json:"preemptionVuner"`
}

// PlmnId representa a Public Land Mobile Network Identity.
type PlmnId struct {
	Mcc string `json:"mcc"` // Mobile Country Code
	Mnc string `json:"mnc"` // Mobile Network Code
}

// UserLocation (Simplificado) - Contém informações de localização da UE.
type UserLocation struct {
	// Tipo de acesso 3GPP
	Tai *Tai `json:"tai,omitempty"`
	// Identificador da Célula E-UTRAN
	Ecgi *Ecgi `json:"ecgi,omitempty"`
	// Identificador da Célula NR (5G)
	Ncgi *Ncgi `json:"ncgi,omitempty"`
	// Coordenadas Geográficas
	GeographicalCoordinates *Coordinates `json:"geographicalCoordinates,omitempty"`
}

// Tai - Tracking Area Identity
type Tai struct {
	PlmnId PlmnId `json:"plmnId"`
	Tac    string `json:"tac"` // Tracking Area Code
}

// Ecgi - E-UTRAN Cell Global Identifier
type Ecgi struct {
	PlmnId      PlmnId `json:"plmnId"`
	EutraCellId string `json:"eutraCellId"`
}

// Ncgi - NR Cell Global Identifier
type Ncgi struct {
	PlmnId   PlmnId `json:"plmnId"`
	NrCellId string `json:"nrCellId"` // NR Cell Identity
}

// Coordinates (Exemplo de sub-struct para UserLocation)
type Coordinates struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

// ServingNfId (Simplificado) - Representa o identificador da NF Servidora (SMF/AMF)
type ServingNfId struct {
	// NF ID (e.g. FQDN, URI, ou ID específico). O tipo exato é complexo
	// na especificação 3GPP, mas um string para o ID é comum.
	NfId string `json:"nfId"`
}

// Tipos de Enumeração (String Constants)

type AccessType string

const (
	AccessType3GPP    AccessType = "3GPP_ACCESS"
	AccessTypeNon3GPP AccessType = "NON_3GPP_ACCESS"
)
