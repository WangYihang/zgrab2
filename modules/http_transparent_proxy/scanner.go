package http_transparent_proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/WangYihang/zgrab2"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Flags struct {
	zgrab2.BaseFlags
	TargetHost     string `long:"target-host" description:"Target host to connect to" default:"ifconfig.me"`
	TargetPort     uint16 `long:"target-port" description:"Target port to connect to" default:"80"`
	RequestTimeout int    `long:"request-timeout" description:"Timeout in seconds" default:"16"`
}

// Validate performs any needed validation on the arguments
func (flags *Flags) Validate(args []string) error {
	return nil
}

// Help returns module-specific help
func (flags *Flags) Help() string {
	return ""
}

type Module struct{}

// NewFlags returns an empty Flags object.
func (module *Module) NewFlags() interface{} {
	return new(Flags)
}

// NewScanner returns a new instance Scanner instance.
func (module *Module) NewScanner() zgrab2.Scanner {
	return new(Scanner)
}

// Description returns an overview of this module.
func (module *Module) Description() string {
	return "Send an HTTP request and read the response, optionally following redirects."
}

// Scanner is the implementation of the zgrab2.Scanner interface.
type Scanner struct {
	Config *Flags
	Domain string
}

func (s *Scanner) Init(flags zgrab2.ScanFlags) error {
	flag, _ := flags.(*Flags)
	s.Config = flag
	return nil
}

func (s *Scanner) InitPerSender(senderID int) error {
	return nil
}

func (s *Scanner) GetName() string {
	return s.Config.Name
}

func (s *Scanner) GetTrigger() string {
	return s.Config.Trigger
}

func (s *Scanner) Protocol() string {
	return "http_transparent_proxy"
}

type Result struct {
	Response *http.Response `json:"response,omitempty"`
}

func CheckTransparentHTTPProxy(index int, host string, port uint16, targetHost string, targetPort uint16, timeout int) error {
	if index%64 == 0 {
		log.Infof("%d, %s:%d, http://%s:%d/", index, host, port, targetHost, targetPort)
	}

	// create a new HTTP client
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// create a new HTTP request
	request, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/", host, port), nil)
	if err != nil {
		// fmt.Println(err)
		return err
	}

	// add headers to the request, if needed
	request.Host = fmt.Sprintf("%s:%d", targetHost, targetPort)

	request.Header.Add("User-Agent", "curl/7.81.0")
	request.Header.Add("NISL-Challenge", uuid.New().String())
	request.Header.Add("NISL-Abuse-Report", "https://pastebin.com/raw/r4g8nddN")
	query := request.URL.Query()
	query.Add("bypass_cache", uuid.New().String())
	request.URL.RawQuery = query.Encode()

	// send the request using the client
	resp, err := client.Do(request)
	if err != nil {
		// fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	// read the response body
	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// fmt.Println(err)
		return err
	}

	strBody := string(rawBody)
	// print the response body
	if strBody == host {
		return nil
	} else {
		return fmt.Errorf(strBody)
	}
}

func (s *Scanner) Scan(t zgrab2.ScanTarget) (zgrab2.ScanStatus, interface{}, error) {
	var port uint16
	if t.Port == -1 {
		port = uint16(s.Config.Port)
	} else {
		port = uint16(t.Port)
	}
	err := CheckTransparentHTTPProxy(t.Index, t.IP.String(), port, s.Config.TargetHost, s.Config.TargetPort, s.Config.RequestTimeout)
	if err != nil {
		return zgrab2.SCAN_PROTOCOL_ERROR, err.Error(), nil
	}
	result := map[string]string{
		"data": "success",
	}
	return zgrab2.SCAN_SUCCESS, result, nil
}

func RegisterModule() {
	var module Module

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	_, err := zgrab2.AddCommand("http_transparent_proxy", "HTTP Proxy Verifier", module.Description(), 80, &module)
	if err != nil {
		log.Fatal(err)
	}
}
