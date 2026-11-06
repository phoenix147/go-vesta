package vesta

import (
	"log/slog"
	"net/http"
	"strconv"
)

func RunApi(httpAddr string, v *Vestatherm) {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/", homeHandler(v))
	mux.HandleFunc("GET /api/v1/heatpump", wpHandler(v))
	mux.HandleFunc("GET /api/v1/tiledstove", koHandler(v))
	mux.HandleFunc("GET /api/v1/buffer", pufferHandler(v))
	mux.HandleFunc("POST /api/v1/buffer/ww/temp/{temp}", wwTempHandler(v))

	http.ListenAndServe(httpAddr, mux)
}

func homeHandler(v *Vestatherm) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v.GoToHome() // avoid screensaver
		v.GoToHome()
		w.WriteHeader(http.StatusOK)
	}
}

func wpHandler(v *Vestatherm) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v.GoToHome() // avoid screensaver
		v.GoToHeatpump()
		w.WriteHeader(http.StatusOK)
	}
}

func koHandler(v *Vestatherm) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v.GoToHome() // avoid screensaver
		v.GoToTiledStove()
		w.WriteHeader(http.StatusOK)
	}
}

func pufferHandler(v *Vestatherm) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v.GoToHome() // avoid screensaver
		v.GoToBuffer()
		w.WriteHeader(http.StatusOK)
	}
}

func wwTempHandler(v *Vestatherm) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v.GoToHome() // avoid screensaver

		tempString := r.PathValue("temp")

		temp, err := strconv.Atoi(tempString)
		if err != nil || temp < 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := v.SetWaterTemperature(temp); err != nil {
			slog.Warn("error setting temperature", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
