package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/service"
)

// @title        IPEHR Gateway API
// @version      0.1
// @description  The IPEHR Gateway is an openEHR compliant EHR server implementation that stores encrypted medical data in a Filecoin distributed file storage.

// @contact.name   API Support
// @contact.url    https://bsn.si/blockchain
// @contact.email  support@bsn.si

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      gateway.ipehr.org
// @BasePath  /v1

type API struct {
	Ehr         *EhrHandler
	EhrStatus   *EhrStatusHandler
	GroupAccess *GroupAccessHandler
	Composition *CompositionHandler
	Query       *QueryHandler
	Cfg         *config.Config

	fs http.FileSystem
}

func New(cfg *config.Config) *API {
	docService := service.NewDefaultDocumentService()
	return &API{
		Ehr:         NewEhrHandler(docService, cfg),
		EhrStatus:   NewEhrStatusHandler(docService, cfg),
		GroupAccess: NewGroupAccessHandler(docService, cfg),
		Composition: NewCompositionHandler(docService, cfg),
		Query:       NewQueryHandler(docService, cfg),
		Cfg:         cfg,
	}
}

func (a *API) Build() *gin.Engine {
	r := gin.New()

	r.Use(a.PrepareContext)

	r.NoRoute(func(c *gin.Context) {
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if p := strings.TrimPrefix(r.URL.Path, "/v1"); len(p) < len(r.URL.Path) {
				if p == "/" || p == "" {
					c.Header("Cache-Control", "no-store, max-age=0")
				}
				c.FileFromFS(p, a.fs)
			} else {
				http.NotFound(w, r)
			}
		}).ServeHTTP(c.Writer, c.Request)
	})

	v1 := r.Group("v1")
	ehr := v1.Group("ehr")
	access := v1.Group("access")
	query := v1.Group("query")

	a.setRedirections(r).
		buildEhrAPI(ehr).
		buildGroupAccessAPI(access).
		buildQueryAPI(query)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func (a *API) buildEhrAPI(r *gin.RouterGroup) *API {
	log.Println("buildEhrAPI")
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	//r.Use(Recovery, app_errors.ErrHandler)

	// Other methods should be authorized
	r.Use(a.Auth)
	r.POST("/", a.Ehr.Create)
	r.GET("/", a.Ehr.GetBySubjectIdAndNamespace)
	r.PUT("/:ehrid", a.Ehr.CreateWithId)
	r.GET("/:ehrid", a.Ehr.GetById)
	r.PUT("/:ehrid/ehr_status", a.EhrStatus.Update)
	r.GET("/:ehrid/ehr_status/:versionid", a.EhrStatus.GetById)
	r.GET("/:ehrid/ehr_status", a.EhrStatus.Get)
	r.POST("/:ehrid/composition", a.Composition.Create)

	return a
}

func (a *API) buildGroupAccessAPI(r *gin.RouterGroup) *API {
	log.Println("buildGroupAccessAPI")
	r.Use(a.Auth)
	r.GET("/group/:group_id", a.GroupAccess.Get)
	r.POST("/group", a.GroupAccess.Create)
	return a
}

func (a *API) buildQueryAPI(r *gin.RouterGroup) *API {
	r.Use(a.Auth)
	r.POST("/", a.Query.ExecPost)
	return a
}

func (a *API) setRedirections(r *gin.Engine) *API {
	redirect := func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "v1/")
	}

	r.GET("/", redirect)
	r.HEAD("/", redirect)
	return a
}

// PrepareContext Fill context with data required for request processing
func (a *API) PrepareContext(c *gin.Context) {
	c.Set("cfg", a.Cfg)
	c.Next()
}
