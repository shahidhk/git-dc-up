package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	gitRepoDir    = mustGetenv("GIT_REPO_DIR", "/repo")
	webhookSecret = mustGetenv("WEBHOOK_SECRET")
	dcExtraFlags  = mustGetenv("DC_EXTRA_FLAGS", "")
	isGithub      = mustGetenv("IS_GITHUB", "1")
)

func index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		if r.URL.Path != "/" {
			http.Error(w, makeResponse(true, http.StatusText(http.StatusNotFound)), http.StatusNotFound)
			return
		}
		if (isGithub == "1" && !isValidSignature(r, webhookSecret)) || (isGithub == "0" && r.Header.Get("x-webhook-secret") != webhookSecret) {
			http.Error(w, makeResponse(true, http.StatusText(http.StatusUnauthorized)), http.StatusUnauthorized)
			return
		}

		// if !isValidSignature(r, webhookSecret) {
		// 	fmt.Println("signature invalid")
		// 	http.Error(w, makeResponse(true, http.StatusText(http.StatusUnauthorized)), http.StatusUnauthorized)
		// 	return
		// }

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

func isValidSignature(r *http.Request, key string) bool {
	// Assuming a non-empty header
	gotHash := strings.SplitN(r.Header.Get("X-Hub-Signature"), "=", 2)
	if gotHash[0] != "sha1" {
		return false
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Cannot read the request body: %s\n", err)
		return false
	}

	hash := hmac.New(sha1.New, []byte(key))
	if _, err := hash.Write(b); err != nil {
		log.Printf("Cannot compute the HMAC for request: %s\n", err)
		return false
	}

	expectedHash := hex.EncodeToString(hash.Sum(nil))
	log.Println("EXPECTED HASH:", expectedHash)
	return gotHash[1] == expectedHash
}
