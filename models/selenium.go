package models

import (
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/tebeka/selenium"
)

type Selenium struct {
	seleniumPath      string
	port              int
	browserName       string
	debug             bool
	remoteURL         string
	userAgent         string
	winName           string
	width             int
	height            int
	driver            selenium.WebDriver
	service           *selenium.Service
	responseCount     uint32
	responseCallbacks []ResponseCallback
	wg                *sync.WaitGroup
	lock              *sync.RWMutex
}

type ResponseCallback func(selenium.WebDriver)

func NewSelenium(options ...func(*Selenium)) *Selenium {
	s := &Selenium{
		wg:                &sync.WaitGroup{},
		lock:              &sync.RWMutex{},
		responseCallbacks: make([]ResponseCallback, 0, 10),
	}
	s.Init()
	for _, f := range options {
		f(s)
	}
	return s
}

func (s *Selenium) Start() {
	opts := []selenium.ServiceOption{
		selenium.Output(os.Stderr),
	}
	selenium.SetDebug(s.debug)
	service, err := selenium.NewChromeDriverService(s.seleniumPath, s.port, opts...)
	if err != nil {
		panic(err)
	}
	s.service = service
	caps := selenium.Capabilities{"browserName": s.browserName}
	driver, err := selenium.NewRemote(caps, s.remoteURL)
	if err != nil {
		panic(err)
	}
	driver.ResizeWindow(s.winName, s.width, s.height)
	s.driver = driver
}

func (s *Selenium) Init() {
	s.browserName = "chrome"
	s.port = 9515
	s.debug = false
	s.seleniumPath = "/usr/local/bin/chromedriver"
	s.remoteURL = fmt.Sprintf("http://localhost:%d/wd/hub", s.port)
	s.winName = ""
	s.height = 200
	s.width = 200
}

func PickUnusedPort() int {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0
	}
	port := l.Addr().(*net.TCPAddr).Port
	if err := l.Close(); err != nil {
		return 0
	}
	return port
}

func SeleniumPath(p string) func(*Selenium) {
	return func(s *Selenium) {
		s.seleniumPath = p
	}
}

func Port(p int) func(*Selenium) {
	return func(s *Selenium) {
		s.port = p
		s.remoteURL = fmt.Sprintf("http://localhost:%d/wd/hub", s.port)
	}
}

func Debug(d bool) func(*Selenium) {
	return func(s *Selenium) {
		s.debug = d
	}
}

func BrowserName(name string) func(*Selenium) {
	return func(s *Selenium) {
		s.browserName = name
	}
}

func RemoteURL(r string) func(*Selenium) {
	return func(s *Selenium) {
		s.remoteURL = r
	}
}

func WinName(m string) func(*Selenium) {
	return func(s *Selenium) {
		s.winName = m
	}
}

func Size(w, h int) func(*Selenium) {
	return func(s *Selenium) {
		s.width = w
		s.height = h
	}
}

func (s *Selenium) Wait() {
	s.wg.Wait()
}

func (s *Selenium) handleOnResponse(r selenium.WebDriver) {
	for _, f := range s.responseCallbacks {
		f(r)
	}
}

func (s *Selenium) Destroy() {
	s.service.Stop()
	s.driver.Quit()
}

func (s *Selenium) OnResponseCallback(fn func(selenium.WebDriver)) {
	s.responseCallbacks = append(s.responseCallbacks, fn)
}

func (s *Selenium) GET(URL string) error {
	s.wg.Add(1)
	defer s.wg.Done()
	err := s.driver.Get(URL)
	if err != nil {
		return err
	}
	s.handleOnResponse(s.driver)
	return nil
}
