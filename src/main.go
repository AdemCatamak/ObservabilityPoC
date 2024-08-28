package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

func main() {

	log.Println("main function was triggered")

	configManager := configManager{}
	configManager.Setup()

	// Create a context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure the context is cancelled on program exit

	// Start the background task
	go recordTestConfig(ctx, TestConfigGaugeMetric)

	// Build Http Server
	port := viper.GetString("PORT")
	server := buildHttpServer(port)

	// Run the HTTP server
	go func() {
		log.Println("Starting server on :" + port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// Wait for an interrupt signal to shut down
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)

	<-stopSignal

	// Cancel the context to stop the background task
	cancel()

	// Wait for the server to shut down gracefully
	log.Println("Shutting down server...")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Println("Server gracefully stopped")

}

func buildHttpServer(port string) *http.Server {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("GET /index",
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprintf(w, "Time: %s", time.DateTime)
		},
	)

	mux.HandleFunc("GET /echo/{statusCode}",
		func(w http.ResponseWriter, r *http.Request) {
			urlParts := strings.Split(r.URL.Path, "/")
			statusCodeStr := urlParts[len(urlParts)-1]
			statusCode, err := strconv.Atoi(statusCodeStr)
			if err != nil {
				http.Error(w, "status-code must be an integer.", http.StatusBadRequest)
				return
			}

			w.WriteHeader(statusCode)
			hostname, err := os.Hostname()
			if err != nil {
				hostname = "unknown"
			}
			_, _ = w.Write([]byte(fmt.Sprintf("Status Code: %d from %s", statusCode, hostname)))
		},
	)

	muxWithMiddleware := requestCounterMetricMiddleware(mux)
	muxWithMiddleware = responseStatusCodeCounterMetricMiddleware(muxWithMiddleware)

	server := &http.Server{Addr: ":" + port, Handler: muxWithMiddleware}
	return server
}

func requestCounterMetricMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.HasSuffix(r.URL.Path, "/metrics") == false {

			hostname, err := os.Hostname()
			if err != nil {
				hostname = "unknown-host"
			}
			RequestCounterMetric.WithLabelValues(hostname).Inc()
		}

		next.ServeHTTP(w, r)
	})
}

func responseStatusCodeCounterMetricMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path
		if strings.Contains(path, "echo") {
			path = "echo"
		}

		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown-host"
		}
		cRW := &customResponseWriter{w, http.StatusOK}

		next.ServeHTTP(cRW, r)

		ResponseStatusCodeCounterMetric.With(prometheus.Labels{"status_code": fmt.Sprint(cRW.statusCode), "path": path, "node": hostname}).Inc()
	})
}

func recordTestConfig(ctx context.Context, testConfigMetric *prometheus.GaugeVec) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("recordTestConfig was stopped...")
			return

		default:
			feedMetric := func(key string) {
				var err error
				var f float64

				if key == "Test1" && viper.GetBool("Test1Error") {
					err = errors.New("Unexpected Error")
				} else {
					f = viper.GetFloat64(key)
				}

				if err != nil {
					testConfigMetric.DeleteLabelValues("Test1")
				} else {
					testConfigMetric.With(prometheus.Labels{"test_config_name": key}).Set(f)
				}
			}

			enabled := viper.GetBool("TestMetricEnabled")
			if enabled {
				feedMetric("Test1")
				feedMetric("Test2")
				feedMetric("Test3")
			} else {
				testConfigMetric.Reset()
			}

			time.Sleep(5 * time.Second)
		}
	}
}
