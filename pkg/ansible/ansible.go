package ansible

import (
  "fmt"
  "os/exec"
  "os/user"
  "os"
  "path/filepath"
  "strconv"
  "syscall"
)

type Ansible struct {
  Cmd
  User		  *user.User
  Logger	  fileLog
}

type fileLog string

func (fl fileLog) Write(data []byte) (int, error) {
  f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY, 0664)
  if err != nil {
    return 0, err
  }
  defer f.Close()
  return f.Write(data)
}

func (fl fileLog) String() string {
  return string(fl)
}

func (a Ansible) Run() {
  usrHome := a.User.HomeDir
  dir := filepath.Join(usrHome, "/rosetta/maneuvers")
  gid, err := idToU32(a.User.Gid)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  uid, err := idToU32(a.User.Uid)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  ansibleCmd := &exec.Cmd {
    Path: a.Path,
    Args: a.Args,
    Dir: dir,
    Stdout: os.Stdout,
    Stderr: os.Stderr,
    SysProcAttr: &syscall.SysProcAttr{
      Credential: &syscall.Credential{
	Uid: uid,
	Gid: gid,
	NoSetGroups: true,
      },
    },
  }

  if a.Logger != "" {
    ansibleCmd.Stderr = a.Logger
  }
  fmt.Println(ansibleCmd.String())
  if err := ansibleCmd.Run(); err != nil {
    fmt.Println(err)
  } 

  return
}

func idToU32 (id string) (uint32, error) {
  id64, err := strconv.ParseUint(id, 10, 32)
  if err != nil {
    return 0, err
  }
  id32 := uint32(id64)
  return id32, nil
}
