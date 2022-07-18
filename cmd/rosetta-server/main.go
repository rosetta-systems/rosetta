package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	
	"github.com/rosetta-systems/rosetta/pkg/ansible"
)

func main() {
	http.HandleFunc("/v0/site/update", siteHandler)
	http.HandleFunc("/v0/maneuvers/update", manHandler)
	http.HandleFunc("/v0/ping", pingHandler)
	log.Fatal(http.ListenAndServe(":8008", nil))
}

// Rebuild & Deploy rosetta.systems container
func siteHandler(w http.ResponseWriter, r *http.Request) {
	var ansi ansible.Runner
	params := ansible.Params{
		User:	"rosetta",
		Log:	"goansi.log",
		Cmd:	ansible.Cmd{
			AnsibleCommand: "playbook",
			Args: []string{	"-u", "rosetta", "-i", "inventory/production", "--tags", "deploy", "playbooks/rosetta.yml" },
		},
	}

	ansi = ansible.New(params)
	ansi.Run()
	w.WriteHeader(http.StatusOK)
}

// Update Maneuvers from latest GitHub Merge
func manHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var payload struct {
		Ref string `json:"ref"`
	}
	err := dec.Decode(&payload)
	if err != nil {
		http.Error(w, "Error: Could Not Read JSON Body", http.StatusBadRequest)
		log.Println(err)
		return
	}
	
	w.WriteHeader(http.StatusOK)

	if strings.Contains(payload.Ref, "main") {
    // checkout repo
		var ansi ansible.Runner
		params := ansible.Params{
			User:	"rosetta",
			Log:	"goansi.log",
			Cmd:	ansible.Cmd{
				AnsibleCommand: "playbook",
				Args: []string{	"-u", "rosetta", "-i", "inventory/production", "--tags", "update-git", "playbooks/rosetta.yml" },
			},
		}

		ansi = ansible.New(params)
		ansi.Run()

    // Update Rosetta Maneuvers, rebuild rosetta-server
		params = ansible.Params{
			User:	"rosetta",
			Log:	"goansi.log",
			Cmd:	ansible.Cmd{
				AnsibleCommand: "playbook",
				Args: []string{	"-u", "rosetta", "-i", "inventory/production", "--skip-tags", "deploy,git", "playbooks/rosetta.yml" },
			},
		}

		ansi = ansible.New(params)
		ansi.Run()
	}
	
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
//	var ansi ansible.Runner
//	params := ansible.Params{
//		User:	"rosetta",
//		Log:	"goansi.log",
//		Cmd:	ansible.Cmd{
//			AnsibleCommand: "playbook",
//			Args: []string{	"-u", "rosetta", "-v", "ping.yml" },
//		},
//	}
//
//	ansi = ansible.New(params)
//	ansi.Run()
	w.Write([]byte("\npong\n\n"))
}
