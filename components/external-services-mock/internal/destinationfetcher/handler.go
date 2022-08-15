package destinationfetcher

import (
	"net/http"
	"path"

	"github.com/kyma-incubator/compass/components/director/pkg/log"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetSensitiveData(writer http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	destinationName := path.Base(req.URL.String())
	log.C(ctx).Infof("Sending sensitive data for destination: %s", destinationName)
	data, ok := destinationsData[destinationName]

	if !ok {
		http.Error(writer, "Not Found", http.StatusNotFound)
		return
	}

	if _, err := writer.Write(data); err != nil {
		log.C(ctx).Errorf("Failed to respond with error %s", err.Error())
	}
}

func (h *Handler) GetSubaccountDestinationsPage(writer http.ResponseWriter, req *http.Request) {
	writer.WriteHeader(http.StatusOK)
}
