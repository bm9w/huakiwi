// [metadata]
// creation_date = "2020/05/04"
// maturity = "production"
// updated_date = "2021/03/03"

// [rule]
// author = ["Elastic"]
// description = """
// Adversaries may attempt to clear or disable the Bash command-line history in an attempt to evade detection or forensic
// investigations.
// """
// from = "now-9m"
// index = ["auditbeat-*", "logs-endpoint.events.*"]
// language = "eql"
// license = "Elastic License v2"
// name = "Tampering of Bash Command-Line History"
// risk_score = 47
// rule_id = "7bcbb3ac-e533-41ad-a612-d6c3bf666aba"
// severity = "medium"
// tags = ["Elastic", "Host", "Linux", "Threat Detection", "Defense Evasion"]
// timestamp_override = "event.ingested"
// type = "eql"

// query = '''
// process where event.type in ("start", "process_started") and
//  (
//   (process.args : ("rm", "echo") and process.args : (".bash_history", "/root/.bash_history", "/home/*/.bash_history")) or
//   (process.name : "history" and process.args : "-c") or
//   (process.args : "export" and process.args : ("HISTFILE=/dev/null", "HISTFILESIZE=0")) or
//   (process.args : "unset" and process.args : "HISTFILE") or
//   (process.args : "set" and process.args : "history" and process.args : "+o")
//  )
// '''

// [[rule.threat]]
// framework = "MITRE ATT&CK"
// [[rule.threat.technique]]
// id = "T1070"
// name = "Indicator Removal on Host"
// reference = "https://attack.mitre.org/techniques/T1070/"
// [[rule.threat.technique.subtechnique]]
// id = "T1070.003"
// name = "Clear Command History"
// reference = "https://attack.mitre.org/techniques/T1070/003/"

// [rule.threat.tactic]
// id = "TA0005"
// name = "Defense Evasion"
// reference = "https://attack.mitre.org/tactics/TA0005/"

package rules

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/mosajjal/ebpf-edr/types"
)

func defense_evasion_deletion_of_bash_command_line_history() error {

	var credSub types.EventSubscriber

	go func(s types.EventSubscriber) {
		log.Println("Running defense_evasion_deletion_of_bash_command_line_history rule")
		s.Source = make(chan types.EventStream, 100)
		s.Subscribe()
		for {
			select {
			case event := <-s.Source:
				args_concat := event.Cmd + " " + strings.Join(event.Args, " ")
				if strings.Contains(args_concat, ".bash_history") || (strings.Contains(args_concat, "history") && strings.Contains(args_concat, "-c")) || (strings.Contains(args_concat, "export") && strings.Contains(args_concat, "HISTFILE")) || (strings.Contains(args_concat, "set") && strings.Contains(args_concat, "history") && strings.Contains(args_concat, "+o")) {

					event_json, _ := json.Marshal(event)
					log.Printf("Deletion of bash history. Severity: Medium. Details: %s\n", string(event_json))

				}
			case <-types.GlobalQuit:
				return
				//todo:write quit
			}
		}
	}(credSub)
	return nil
}

var _ = defense_evasion_deletion_of_bash_command_line_history()