package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/composition"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/errors"
)

type CompositionHandler struct {
	cfg     *config.Config
	service *composition.CompositionService
}

func NewCompositionHandler(docService *service.DefaultDocumentService, cfg *config.Config) *CompositionHandler {
	return &CompositionHandler{
		cfg:     cfg,
		service: composition.NewCompositionService(docService),
	}
}

// Create
// @Summary      Create COMPOSITION
// @Description  Work in progress...
// @Description  Creates the first version of a new COMPOSITION in the EHR identified by ehr_id.
// @Description
// @Tags     COMPOSITION
// @Accept   json
// @Produce  json
// @Param    ehr_id      path      string                 true  "EHR identifier. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
// @Param    AuthUserId  header    string                 true  "UserId UUID"
// @Param    Prefer      header    string                 true  "The new EHR resource is returned in the body when the request’s `Prefer` header value is `return=representation`, otherwise only headers are returned."
// @Param    Request     body      model.SwagComposition  true  "COMPOSITION"
// @Success  201         {object}  model.SwagComposition
// @Header   201         {string}  Location  "{baseUrl}/ehr/7d44b88c-4199-4bad-97dc-d78268e01398/composition/8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::1"
// @Header   201         {string}  ETag      "8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::1"
// @Failure  400         "Is returned when the request has invalid ehr_id or invalid content (e.g. content could not be converted to a valid COMPOSITION object)"
// @Failure  404         "Is returned when an EHR with ehr_id does not exist."
// @Failure  422         "Is returned when the content could be converted to a COMPOSITION, but there are semantic validation errors, such as the underlying template is not known or is not validating the supplied COMPOSITION)."
// @Failure  500         "Is returned when an unexpected error occurs while processing a request"
// @Router   /ehr/{ehr_id}/composition [post]
func (h CompositionHandler) Create(c *gin.Context) {
	ehrId := c.Param("ehrid")
	if h.service.Doc.ValidateId(ehrId, types.EHR) == false {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Request body error"})
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Panic("Cant close body request")
		}
	}(c.Request.Body)

	var request model.Composition

	if err = json.Unmarshal(data, &request); err != nil {
		log.Println("Composition Create request unmarshal error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request validation error"})
		return
	}

	if !request.Validate() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request validation error"})
		return
	}

	userId := c.GetString("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	// Checking EHR does not exist
	_, err = h.service.Doc.EhrsIndex.Get(userId)
	if errors.Is(err, errors.IsNotExist) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Composition document creating
	doc, err := h.service.CompositionCreate(userId, ehrId, &request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Composition creating error"})
		return
	}

	h.respondWithDocOrHeaders(ehrId, doc, c)
}

func (h *CompositionHandler) respondWithDocOrHeaders(ehrId string, doc *model.Composition, c *gin.Context) {
	uid := doc.Uid.Value
	c.Header("Location", h.cfg.BaseUrl+"/v1/ehr/"+ehrId+"/composition/"+uid)
	c.Header("ETag", uid)

	prefer := c.Request.Header.Get("Prefer")
	if prefer == "return=representation" {
		c.JSON(http.StatusCreated, doc)
	} else {
		c.AbortWithStatus(http.StatusCreated)
	}
}
