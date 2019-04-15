package srv

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/elastic/go-elasticsearch/v7"

	"github.com/kainobor/estest/app/lib/logger"
	"github.com/kainobor/estest/app/lib/session"
	"github.com/kainobor/estest/app/pkg/product"
	productRepo "github.com/kainobor/estest/app/pkg/product/repo"
	productSvc "github.com/kainobor/estest/app/pkg/product/svc"
	"github.com/kainobor/estest/app/pkg/user"
	userRepo "github.com/kainobor/estest/app/pkg/user/repo"
	userSvc "github.com/kainobor/estest/app/pkg/user/svc"
	prodCtrlV1 "github.com/kainobor/estest/app/srv/ctrl/v1/product"
	userCtrlV1 "github.com/kainobor/estest/app/srv/ctrl/v1/user"
	"github.com/kainobor/estest/config"
)

type Server struct {
	config  config.Server
	router  *mux.Router
	userSvc user.Service
	prodSvc product.Service
}

func New(conf config.Server) *Server {
	return &Server{
		config: conf,
	}
}

func (srv *Server) Start(ctx context.Context, es *elasticsearch.Client) {
	log := logger.New(ctx)
	if srv == nil {
		log.Fatalw("Server can't be nil")
	}

	srv.buildServices(es)

	srv.router = mux.NewRouter()
	srv.addV1Routes(srv.router.PathPrefix("/v1").Subrouter())

	go func() {
		if err := http.ListenAndServe(":"+srv.config.Port, srv.router); err != nil {
			log.Fatalw("Can't start server", "error", err)
		}
	}()

	log.Info("Start to listen")

	// Wait until end
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	select {
	case s := <-c:
		log.Infof("Listening finished with signal %s", s.String())
	}
}

func (srv *Server) addV1Routes(r *mux.Router) {
	prod := prodCtrlV1.New(srv.prodSvc)
	usr := userCtrlV1.New(srv.userSvc)
	r.HandleFunc("/login", usr.Login).Methods(http.MethodGet)
	r.HandleFunc("/products", srv.withAuth(prod.GetProduct)).Methods(http.MethodGet)
}

func (srv *Server) buildServices(dbClient *elasticsearch.Client) {
	usrRepo := userRepo.NewRepository(dbClient)
	prodRepo := productRepo.New(dbClient)
	srv.prodSvc = productSvc.New(prodRepo)
	srv.userSvc = userSvc.New(usrRepo)
}

func (srv *Server) withAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := session.IsAuth(r)
		if err != nil {
			log := logger.New(r.Context())
			log.Errorf("failed to authorize: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "UserID", id)
		r.WithContext(ctx)

		handlerFunc.ServeHTTP(w, r)
	})
}
