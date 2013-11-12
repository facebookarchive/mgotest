// Package mgotest provides standalone test instances of mongo sutable for use
// in tests.
package mgotest

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"text/template"

	"labix.org/v2/mgo"

	"github.com/ParsePlatform/go.freeport"
)

var configTemplate, configTemplateErr = template.New("config").Parse(`
port = {{.Port}}
dbpath = {{.DBPath}}
nounixsocket = true
nojournal = true
nohttpinterface = true
noprealloc = true
nssize = 2
smallfiles = true
quiet = true
`)

func init() {
	if configTemplateErr != nil {
		panic(configTemplateErr)
	}
}

// Fatalf is satisfied by testing.T or testing.B.
type Fatalf interface {
	Fatalf(format string, args ...interface{})
}

// Server is a unique instance of a mongod.
type Server struct {
	Port   int
	DBPath string
	T      Fatalf
	cmd    *exec.Cmd
}

// Start the server, this will return once the server has been started.
func (s *Server) Start() {
	port, err := freeport.Get()
	if err != nil {
		s.T.Fatalf(err.Error())
	}
	s.Port = port

	dir, err := ioutil.TempDir("", "mgotest-dbpath-")
	if err != nil {
		s.T.Fatalf(err.Error())
	}
	s.DBPath = dir

	cf, err := ioutil.TempFile(s.DBPath, "config-")
	if err != nil {
		s.T.Fatalf(err.Error())
	}

	if err := configTemplate.Execute(cf, s); err != nil {
		s.T.Fatalf(err.Error())
	}

	s.cmd = exec.Command("mongod", "--config", cf.Name())
	s.cmd.Env = envPlusLcAll()
	s.cmd.Stdout = os.Stdout
	s.cmd.Stderr = os.Stderr
	if err := s.cmd.Start(); err != nil {
		s.T.Fatalf(err.Error())
	}
}

// Stop the server, this will also remove all data.
func (s *Server) Stop() {
	s.cmd.Process.Kill()
	os.RemoveAll(s.DBPath)
}

// URL for the mongo server, suitable for use with mgo.Dial.
func (s *Server) URL() string {
	return fmt.Sprintf("localhost:%d", s.Port)
}

// Session for the mongo server.
func (s *Server) Session() *mgo.Session {
	session, err := mgo.Dial(s.URL())
	if err != nil {
		s.T.Fatalf(err.Error())
	}
	return session
}

// NewStartedServer creates a new server starts it.
func NewStartedServer(t Fatalf) *Server {
	s := &Server{T: t}
	s.Start()
	return s
}

func envPlusLcAll() []string {
	env := os.Environ()
	return append(env, "LC_ALL=C")
}
