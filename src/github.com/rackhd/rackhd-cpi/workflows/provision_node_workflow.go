package workflows

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/rackhd/rackhd-cpi/config"
	"github.com/rackhd/rackhd-cpi/rackhdapi"
)

var provisionNodeWorkflowTemplate = []byte(`{
  "friendlyName": "BOSH Provision Node",
  "injectableName": "Graph.BOSH.ProvisionNode",
  "options": {
    "defaults": {
      "agentSettingsFile": null,
      "agentSettingsPath": null,
      "cid": null,
      "downloadDir": "/opt/downloads",
      "registrySettingsFile": null,
      "registrySettingsPath": null,
      "stemcellFile": null
    }
  },
  "tasks": [
    {
      "label": "bootstrap-ubuntu",
      "taskName": "Task.Linux.Bootstrap.Ubuntu",
      "ignoreFailure": true
    },
    {
      "label": "provision-node",
      "taskName": "Task.BOSH.Provision.Node",
	  "waitOn": {
	    "bootstrap-ubuntu": "finished"
	  }
    },
    {
      "label": "set-id",
      "taskName": "Task.BOSH.SetNodeId",
      "waitOn": {
        "provision-node": "succeeded"
      }
    },
	{
	  "label": "reboot",
	  "taskName": "Task.ProcShellReboot",
	  "waitOn": {
	    "set-id": "succeeded"
	  }
	}
  ]
}`)

type ProvisionNodeWorkflowOptions struct {
	AgentSettingsFile    *string `json:"agentSettingsFile"`
	AgentSettingsPath    *string `json:"agentSettingsPath"`
	CID                  *string `json:"cid"`
	DownloadDir          string  `json:"downloadDir,omitempty"`
	RegistrySettingsFile *string `json:"registrySettingsFile"`
	RegistrySettingsPath *string `json:"registrySettingsPath"`
	StemcellFile         *string `json:"stemcellFile"`
}

type provisionNodeWorkflowOptionsContainer struct {
	Options provisionNodeWorkflowDefaultOptionsContainer `json:"options"`
}

type provisionNodeWorkflowDefaultOptionsContainer struct {
	Defaults ProvisionNodeWorkflowOptions `json:"defaults"`
}

type provisionNodeWorkflow struct {
	*rackhdapi.WorkflowStub
	*provisionNodeWorkflowOptionsContainer
	Tasks []rackhdapi.WorkflowTask `json:"tasks"`
}

func RunProvisionNodeWorkflow(c config.Cpi, nodeID string, workflowName string, options ProvisionNodeWorkflowOptions) error {
	req := rackhdapi.RunWorkflowRequestBody{
		Name:    workflowName,
		Options: map[string]interface{}{"defaults": options},
	}

	return rackhdapi.RunWorkflow(rackhdapi.WorkflowPoster, rackhdapi.WorkflowFetcher, c, nodeID, req)
}

func PublishProvisionNodeWorkflow(cpiConfig config.Cpi, uuid string) (string, error) {
	tasks, workflow, err := generateProvisionNodeWorkflow(uuid)
	if err != nil {
		return "", err
	}

	for i := range tasks {
		err = rackhdapi.PublishTask(cpiConfig, tasks[i])
		if err != nil {
			return "", err
		}
	}

	w := provisionNodeWorkflow{}
	err = json.Unmarshal(workflow, &w)
	if err != nil {
		log.Error(fmt.Sprintf("error umarshalling workflow: %s", err))
		return "", err
	}

	err = rackhdapi.PublishWorkflow(cpiConfig, workflow)
	if err != nil {
		return "", err
	}

	return w.Name, nil
}

func generateProvisionNodeWorkflow(uuid string) ([][]byte, []byte, error) {
	p := provisionNodeTask{}
	err := json.Unmarshal(provisionNodeTemplate, &p)
	if err != nil {
		log.Error(fmt.Sprintf("error unmarshalling provision node task template: %s\n", err))
		return nil, nil, fmt.Errorf("error unmarshalling provision node task template: %s\n", err)
	}

	p.Name = fmt.Sprintf("%s.%s", p.Name, uuid)
	p.UnusedName = fmt.Sprintf("%s.%s", p.UnusedName, "UPLOADED_BY_RACKHD_CPI")

	pBytes, err := json.Marshal(p)
	if err != nil {
		log.Error(fmt.Sprintf("error marshalling provision node task template: %s\n", err))
		return nil, nil, fmt.Errorf("error marshalling provision node task template: %s\n", err)
	}

	s := setNodeIDTask{}
	err = json.Unmarshal(setNodeIDTemplate, &s)
	if err != nil {
		log.Error(fmt.Sprintf("error unmarshalling set node id task template: %s\n", err))
		return nil, nil, fmt.Errorf("error unmarshalling set node id task template: %s\n", err)
	}

	s.Name = fmt.Sprintf("%s.%s", s.Name, uuid)
	s.UnusedName = fmt.Sprintf("%s.%s", s.UnusedName, "UPLOADED_BY_RACKHD_CPI")

	sBytes, err := json.Marshal(s)
	if err != nil {
		log.Error(fmt.Sprintf("error marshalling set node id task template: %s\n", err))
		return nil, nil, fmt.Errorf("error marshalling set node id task template: %s\n", err)
	}

	w := provisionNodeWorkflow{}
	err = json.Unmarshal(provisionNodeWorkflowTemplate, &w)
	if err != nil {
		log.Error(fmt.Sprintf("error unmarshalling provision node workflow template: %s\n", err))
		return nil, nil, fmt.Errorf("error unmarshalling provision node workflow template: %s\n", err)
	}

	w.Name = fmt.Sprintf("%s.%s", w.Name, uuid)
	w.UnusedName = fmt.Sprintf("%s.%s", w.UnusedName, "UPLOADED_BY_RACKHD_CPI")
	w.Tasks[1].TaskName = fmt.Sprintf("%s.%s", w.Tasks[1].TaskName, uuid)
	w.Tasks[2].TaskName = fmt.Sprintf("%s.%s", w.Tasks[2].TaskName, uuid)

	wBytes, err := json.Marshal(w)
	if err != nil {
		log.Error(fmt.Sprintf("error marshalling provision node workflow template: %s\n", err))
		return nil, nil, fmt.Errorf("error marshalling provision node workflow template: %s\n", err)
	}

	return [][]byte{pBytes, sBytes}, wBytes, nil
}