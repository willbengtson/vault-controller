package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/vault/api"
	log "github.com/Sirupsen/logrus"
)

const VERSION = "0.0.1"

var (
	vaultTokenFile string
	logFormat string
	logLevel string
	addr string
)

func main() {
	flag.StringVar(&addr, "addr", "", "Vault Server Address")
	flag.StringVar(&vaultTokenFile, "token_file", "", "File path to vault auth response (JSON)")
	flag.StringVar(&logFormat, "log_fmt", "", "Log Format - json or tty")
	flag.StringVar(&logLevel, "log_level", "", "Log Level - info, debug, warn")
	versionPtr := flag.Bool("version", false, "Prints version")
	flag.Parse()

	if *versionPtr {
		fmt.Println(VERSION)
		os.Exit(0)
	}
	
	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)

	// Set the log format.  Default to Text
	if logFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{})
	}

	// Set the log level and default to INFO
	switch logLevel {
		case "warn":
			log.SetLevel(log.WarnLevel)
		case "debug":
			log.SetLevel(log.DebugLevel)
		default:
			log.SetLevel(log.InfoLevel)
	}

	log.Info("Starting vault-controller...")

	// Check to see if we are passing in a Vault address.  If so, use it.  
	// If not, attempt to get it from VAULT_ADDR environment variable
	// If bot are not set, panic
	if addr == "" {
		log.Debug("--addr not set.  Attemping to get addr from VAULT_ADDR environment variable")
		addr = os.Getenv("VAULT_ADDR")
		if addr == "" {
			log.Fatal("VAULT_ADDR not set...exiting")
		}
	}

	var secret *api.Secret

	// Setup the Vault Client with default settings and then assign the address
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}
	client.SetAddress(addr)

	if vaultTokenFile != "" {
		log.Info("Reading vault secret file from ", vaultTokenFile)
		if _, err := os.Stat(vaultTokenFile); err != nil {
		    if os.IsNotExist(err) {
		        log.Fatal("Vault secret file does not exist: ", vaultTokenFile)
		    } else {
		        // other error
		    }
		}
		file, err := ioutil.ReadFile(vaultTokenFile)
		if err != nil {
			log.Fatal("Could not read Vault auth response file: ", err)
		}
		secret, err = api.ParseSecret(bytes.NewReader(file))
		if err != nil {
			log.Fatal("Could not parse Vault secret: ", err)
		}

		// Log the Vault Auth response if we are in Debug logging
		logVaultSecret(secret)

		if !secret.Auth.Renewable {
			log.Fatal("Vault token is not renewable.")
		}

		client.SetToken(secret.Auth.ClientToken)
	} else { 
		// Vault token response file must be set
		log.Fatal("Vault auth response not found.  Exiting...")
	}

	// Set the retry to 5 seconds in case Vault is unavailable or there is an error in renewal
	retryDelay := 5 * time.Second
	go func() {
		for {
			tokeRenewal, err := client.Auth().Token().RenewSelf(secret.Auth.LeaseDuration)
			if err != nil {
				log.Info("token-renew: Renew client token error: ", err, "; retrying in ", retryDelay)
				time.Sleep(retryDelay)
				continue
			}

			// To play it safe, let's renew at the half way point.
			nextRenew := tokeRenewal.Auth.LeaseDuration / 2
			log.Info("token-renew: Successfully renewed the client token; next renewal in ", nextRenew, " seconds")
			// Sleep for the renewal wait period
			time.Sleep(time.Duration(nextRenew) * time.Second)
		}
	}()

	// Check for interrupt signal and exit cleanly
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown signal received, exiting vault-controller...")
}

// Function to log the Vault Secret Auth Response.  This only works if log_level is debug
func logVaultSecret(authResponse *api.Secret) {
	jsonObject, err := json.MarshalIndent(&authResponse, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	log.Debug(string(jsonObject))
}
