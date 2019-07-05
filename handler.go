package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

var (
	gitRepoDir    = mustGetenv("GIT_REPO_DIR", "/repo")
	webhookSecret = mustGetenv("WEBHOOK_SECRET")
	dcExtraFlags  = mustGetenv("DC_EXTRA_FLAGS", "")
)

func index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		if r.URL.Path != "/" {
			http.Error(w, makeResponse(true, http.StatusText(http.StatusNotFound)), http.StatusNotFound)
			return
		}
		if r.Header.Get("x-webhook-secret") != webhookSecret {
			http.Error(w, makeResponse(true, http.StatusText(http.StatusUnauthorized)), http.StatusUnauthorized)
			return
		}

		err := executeAction()
		if err != nil {
			http.Error(w, makeResponse(true, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, makeResponse(false, "success"))
	})
}

func executeAction() error {
	gitPull := exec.Command("git", "-C", gitRepoDir, "pull")
	stdoutStderr, err := gitPull.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Printf("git pull: %s\n", stdoutStderr)
	dcUp := exec.Command("sh", "-c", fmt.Sprintf(`cd %s && docker-compose up -d %s`, gitRepoDir, dcExtraFlags))
	stdoutStderr, err = dcUp.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Printf("dc up: %s\n", stdoutStderr)
	return nil
}

func makeResponse(error bool, message string) string {
	return fmt.Sprintf(`{"error": %s, "message": "%s"}`, strconv.FormatBool(error), message)
}

// MustGetenv returns env value for key and panics if the value is not set.
func mustGetenv(key string, def ...string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	panic(fmt.Errorf("%s env not set", key))
	return ""
}
