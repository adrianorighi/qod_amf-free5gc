package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/free5gc/nef/internal/logger"
	"github.com/free5gc/nef/internal/models"
	"github.com/gin-gonic/gin"
)

// TokenResponse captura a resposta do /oauth2/token do NRF
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func getAccessToken(nrfBaseUri, consumerId, targetService string) (string, error) {
	tokenUrl := fmt.Sprintf("%s/oauth2/token", nrfBaseUri)

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("nfInstanceId", consumerId)
	data.Set("targetNfType", "PCF")
	data.Set("nfType", "SMF")
	data.Set("scope", targetService)

	encodedData := data.Encode()

	bodyReader := strings.NewReader(encodedData)

	req, err := http.NewRequest("POST", tokenUrl, bodyReader)
	if err != nil {
		return "", fmt.Errorf("erro ao criar request para NRF: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	req.Header.Set("Accept", "application/json")

	req.ContentLength = int64(len(encodedData))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao enviar request para NRF: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("NRF respondeu com status %s: %s", resp.Status, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta JSON do NRF: %v", err)
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("NRF não retornou um access_token no JSON")
	}

	return tokenResp.AccessToken, nil
}

func (p *Processor) PostQoDSessionToPcf(
	c *gin.Context,
	qodSession models.QoDManagement,
) {
	logger.TrafInfluLog.Infof("PostQoDSession: %+v", qodSession)
	nrfUrL := "http://nrf.free5gc.org:8000"
	consumerNfInstanceID := "c49c146d-65d8-4c00-b627-73ccf62ecd47"
	targetNfService := "npcf-policyauthorization"
	pcfURL := os.Getenv("PCF_QOD_URL") // "http://pcf.free5gc.org:8000/npcf-policyauthorization/v1/app-sessions"
	gbr5QI := int32(1)

	// Taxa desejada: 250 Mbps (usei string no formato da especificação 3GPP)
	// Fixado para fins de teste (deve-se usar perfis de QoS)
	taxa250Mbps := "250 Mbps"

	if pcfURL == "" {
		logger.TrafInfluLog.Warnf("PCF_QOD_URL not set; skipping forward to PCF")
	} else {
		auth, err := getAccessToken(nrfUrL, consumerNfInstanceID, targetNfService)
		if err != nil {
			log.Fatalf("Falha ao obter token do NRF: %v", err)
		}
		log.Printf("Token recebido com sucesso: %s...", auth[:20])

		pcfRequestBody := models.SmPolicyContextData{
			Supi:           qodSession.Device.NetworkAccessIdentifier,
			PduSessionId:   10,
			Dnn:            "internet",
			PduSessionType: "IPV4V6",
			AccessType:     models.AccessType3GPP,
			Snssai: models.Snssai{
				Sst: 1,
				Sd:  "010203",
			},

			SubscribedSessionAmbr: &models.Ambr{
				Uplink:   taxa250Mbps,
				Downlink: taxa250Mbps,
			},

			SubscribedDefaultQos: &models.SubscribedDefaultQos{
				FiveQi: gbr5QI,
				Arp: models.Arp{
					PriorityLevel:   1,
					PreemptionCap:   "MAY_PREEMPT",
					PreemptionVuner: "NOT_PREEMPTABLE",
				},
			},
			Ipv4Address: qodSession.ApplicationServer.IPv4Address,
			UeLocation: &models.UserLocation{
				Tai: &models.Tai{
					PlmnId: models.PlmnId{Mcc: "208", Mnc: "93"},
				},
			},
		}

		payload, err := json.Marshal(pcfRequestBody)
		if err != nil {
			logger.TrafInfluLog.Errorf("failed to marshal QoD session: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize QoD session"})
			return
		}

		req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodPost, pcfURL, bytes.NewReader(payload))
		if err != nil {
			logger.TrafInfluLog.Errorf("failed to build PCF request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build PCF request"})
			return
		}
		req.Header.Set("Content-Type", "application/json")
		if auth != "" {
			logger.TrafInfluLog.Infof("Using Authorization: Bearer %s", auth[:20])
			req.Header.Set("Authorization", "Bearer "+auth)
		}

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			logger.TrafInfluLog.Errorf("PCF request failed: %v", err)
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to reach PCF"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)
			logger.TrafInfluLog.Errorf("PCF responded %d: %s", resp.StatusCode, string(body))
			c.JSON(resp.StatusCode, gin.H{"error": "pcf error", "status": resp.Status})
			return
		}
	}
	c.JSON(http.StatusCreated, qodSession)
}
