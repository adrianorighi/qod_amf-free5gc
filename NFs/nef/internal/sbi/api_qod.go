package sbi

import (
	"net/http"

	"github.com/free5gc/nef/internal/logger"
	"github.com/free5gc/nef/internal/models"
	"github.com/free5gc/openapi"
	"github.com/gin-gonic/gin"
)

func (s *Server) getQodRoutes() []Route {
	return []Route{
		{
			Method:  http.MethodPost,
			Pattern: "/sessions",
			APIFunc: s.apiPostQoDManagementTransactions,
		},
	}
}

func (s *Server) apiPostQoDManagementTransactions(gc *gin.Context) {
	var qosData models.QoDManagement
	reqBody, err := gc.GetRawData()
	logger.SBILog.Infof("Request Body: %s", string(reqBody))
	if err != nil {
		logger.SBILog.Errorf("Get Request Body error: %+v", err)
		gc.JSON(http.StatusInternalServerError,
			openapi.ProblemDetailsSystemFailure(err.Error()))
		return
	}
	err = openapi.Deserialize(&qosData, reqBody, "application/json")
	if err != nil {
		logger.SBILog.Errorf("Deserialize Request Body error: %+v", err)
		gc.JSON(http.StatusBadRequest,
			openapi.ProblemDetailsMalformedReqSyntax(err.Error()))
		return
	}

	// Log the decoded QoDManagement struct
	logger.SBILog.Infof("QoDManagement decoded: %+v", qosData)

	// Chame o método do Processor (ajustado para função getter)
	s.Processor().PostQoDSessionToPcf(gc, qosData)
}
