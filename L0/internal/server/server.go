package server

import (
	"L0/internal/repository"
	"L0/internal/util"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/stan.go"
)

type Server struct {
	config *util.Config
	logger *util.Logger
	server *http.Server
	mux    *http.ServeMux
	repos  *repository.Repository
	cache  *repository.Cache
}

func NewServer(config *util.Config, repo *repository.Repository, cache *repository.Cache) *Server {
	return &Server{
		config: config,
		logger: &util.Logger{},
		mux:    http.NewServeMux(),
		repos:  repo,
		cache:  cache,
	}
}

func (s *Server) Start() error {
	if err := s.logger.InitLogger(s.config.LevelLoger); err != nil {
		return err
	}

	if err := s.config.InitConfig(); err != nil {
		return err
	}

	if err := s.repos.InitDB(); err != nil {
		return err
	}

	ns, err := s.сonfigureNats()
	if err != nil {
		s.logger.Logger.Err(err).Msg("couldn't connect to the nats channel, check the connection")
		return err
	}

	if err = s.repos.NatsGenerateDate(ns); err != nil {
		s.logger.Logger.Err(err).Msg("failed to add data to nats channel")
	}

	if err = s.repos.NatsSubcribe(ns); err != nil {
		s.logger.Logger.Err(err).Msg("the subscription to the channel was not successful")
	}

	s.mux.HandleFunc("GET /order/{id}/", s.hadlerGetOrder)

	s.server = &http.Server{
		Addr:    ":" + s.config.ServerPort,
		Handler: s.mux,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		s.logger.Logger.Info().Msg("received shutdown signal")
		ns.Close()
		s.repos.Close()
		s.server.Close()
		os.Exit(0)
	}()

	return http.ListenAndServe(":"+s.config.ServerPort, s.mux)
}

func (s *Server) hadlerGetOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	order, ok := s.cache.GetOrderByUID(id)
	if !ok {
		http.Error(w, "data not found", http.StatusInternalServerError)
	}

	json, err := json.MarshalIndent(order, "", " ")
	if err != nil {
		s.logger.Logger.Err(err).Msg("failed to convert json to byte array")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(json); err != nil {
		s.logger.Logger.Err(err).Msg("unable to convert json")
	}

}

func (s *Server) сonfigureNats() (stan.Conn, error) {
	sc, err := stan.Connect(s.config.NatsClussterID, s.config.NatsClientID, stan.NatsURL(s.config.NatsURL))
	if err != nil {
		return nil, err
	}

	return sc, nil
}
