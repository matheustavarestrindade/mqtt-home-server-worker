package http

import (
	"fmt"
	"net/http"

	"github.com/matheustavarestrindade/mqtt-home-server-worker/internal/database"
	"github.com/matheustavarestrindade/mqtt-home-server-worker/internal/http/endpoints"
)

type Server struct {
	Port            int
	sensorsEndpoint *endpoints.SensorEndpoints
}

func NewServer(port int, database *database.Database) *Server {
	server := &Server{
		Port:            port,
		sensorsEndpoint: endpoints.NewSensorEndpoints(database),
	}

	http.HandleFunc("/sensors", server.sensorsEndpoint.GetSensorsByID)
	http.HandleFunc("/sensor/data", server.sensorsEndpoint.GetSensorDataByIDAndTimestamp)
	go func() {
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()
	return server
}
