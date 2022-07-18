package ansible

import (
  "fmt"
  "os"
  "os/user"
)

type Runner interface {
  Command() error  
  Run()
}

type Params struct {
  Cmd  Cmd
  User	   string
  Log	   string
}

func New (p Params) Runner {
  a := &Ansible{Cmd: p.Cmd}
  u, err := user.Lookup(p.User)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  a.User = u
  a.Logger = fileLog(p.Log)

  err = a.Command()
  if err != nil {
    fmt.Println(err)
  }
  return a
}
