package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func startExporter(binaryPath string, configFile string, terminationChan chan os.Signal) error {
	cmd := exec.Command(binaryPath, "--config.file="+configFile, "--web.listen-address=127.0.0.1:9116")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	go func() {
		<-terminationChan
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}()

	return cmd.Run()
}

func main() {
	snmpConfig := flag.String("snmp-config", "/etc/snmp_exporter/snmp.yml", "absolute path to snmp exporter config file")
	exporterExecutable := flag.String("exporter-executable", "/bin/snmp_exporter", "absolute path to snmp exporter")
	token := flag.String("token", "changeIt!", "token to verify requests")
	flag.Parse()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGTERM)
	go func() {
		err := startExporter(*exporterExecutable, *snmpConfig, termChan)
		log.Fatal(err)
	}()

	http.HandleFunc("/metrics", func(writer http.ResponseWriter, request *http.Request) {
		params := request.URL.Query()
		requestToken := params.Get("token")
		module := params.Get("module")
		target := params.Get("target")

		if requestToken != *token {
			log.Printf("Invalid token: %s", token)
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		if module == "" || target == "" {
			writer.Write([]byte("Missing parameters"))
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		metricsUrl := fmt.Sprintf("http://localhost:9116/snmp?module=%s&target=%s", module, target)
		log.Printf("Get metrics for url %s", metricsUrl)
		res, err := http.Get(metricsUrl)

		if err != nil {
			log.Printf("Error tunneling request: %s", err.Error())
			writer.WriteHeader(http.StatusBadGateway)
			return
		}
		_, _ = io.Copy(writer, res.Body)
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
