package ansible

import (
  //"errors"
  "fmt"
  "os/exec"
)

// Cmd defines the variable options for running ansible
type Cmd struct {
  AnsibleCommand  supportedAnsibleCommand
  Path		  string
  Args		  []string
}

type supportedAnsibleCommand string

// Command only to be used when all Cmd fields are passed in together
func (c *Cmd) Command() error {
  // Ensure AnsibleCommand is supported
  err := c.validateAnsibleCommand()
  if err != nil {
    return err
  }
  
  // Get 'AnsibleCommand' executable
  ansibleExecPath, err := exec.LookPath("ansible-" + string(c.AnsibleCommand))
  if err != nil {
    return err
  }
  c.Path = ansibleExecPath

  c.AnsibleArgs(ansibleExecPath, c.Args)

  return nil
}

func (c *Cmd) validateAnsibleCommand() error {
  const (
    playbook supportedAnsibleCommand = "playbook"
  )
  ac := c.AnsibleCommand
  switch ac {
  case playbook:
    return nil
  }
  return InvalidAnsibleCommand(ac)
}

type InvalidAnsibleCommand string
func (e InvalidAnsibleCommand) Error() string {
  return fmt.Sprintf("%s is not a supported ansbile command.\nMust be one of:\n\tplaybook", string(e))
}

func (c *Cmd) AnsibleArgs(execPath string, args []string) {
  //TODO create validateAnsibleArgs
  args = append(args, "")
  copy(args[1:], args)
  args[0] = execPath
  c.Args = args 
}
