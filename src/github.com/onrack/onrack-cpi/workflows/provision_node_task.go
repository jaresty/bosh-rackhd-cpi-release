package workflows

import "github.com/onrack/onrack-cpi/onrackhttp"

type provisionNodeOptions struct {
	AgentSettingsFile      *string  `json:"agentSettingsFile"`
	AgentSettingsMd5Uri    string   `json:"agentSettingsMd5Uri"`
	AgentSettingsPath      *string  `json:"agentSettingsPath"`
	AgentSettingsURI       string   `json:"agentSettingsUri"`
	Commands               []string `json:"commands"`
	Device                 string   `json:"device"`
	DownloadDir            string   `json:"downloadDir"`
	PublicKeyFile          *string  `json:"publicKeyFile"`
	PublicKeyMd5Uri        string   `json:"publicKeyMd5Uri"`
	PublicKeyURI           string   `json:"publicKeyUri"`
	RegistrySettingsFile   *string  `json:"registrySettingsFile"`
	RegistrySettingsMd5Uri string   `json:"registrySettingsMd5Uri"`
	RegistrySettingsPath   *string  `json:"registrySettingsPath"`
	StemcellFileMd5Uri     string   `json:"stemcellFileMd5Uri"`
	RegistrySettingsURI    string   `json:"registrySettingsUri"`
	StemcellFile           *string  `json:"stemcellFile"`
	StemcellURI            string   `json:"stemcellUri"`
}

type provisionNodeTask struct {
	*onrackhttp.TaskStub
	*onrackhttp.PropertyContainer
	*provisionNodeOptionsContainer
}

type provisionNodeOptionsContainer struct {
	Options provisionNodeOptions `json:"options"`
}

var provisionNodeTemplate = []byte(`{
  "friendlyName": "Provision Node",
  "implementsTask": "Task.Base.Linux.Commands",
  "injectableName": "Task.BOSH.Provision.Node",
  "options": {
    "agentSettingsFile": null,
    "agentSettingsMd5Uri": "{{ api.files }}/md5/{{ options.agentSettingsFile }}/latest",
    "agentSettingsPath": null,
    "agentSettingsUri": "{{ api.files }}/{{ options.agentSettingsFile }}/latest",
    "commands": [
      "curl --retry 3 {{ options.stemcellUri }} -o {{ options.downloadDir }}/{{ options.stemcellFile }}",
      "curl --retry 3 {{ options.agentSettingsUri }} -o {{ options.downloadDir }}/{{ options.agentSettingsFile }}",
      "curl --retry 3 {{ options.registrySettingsUri }} -o {{ options.downloadDir }}/{{ options.registrySettingsFile }}",
      "curl --retry 3 {{ options.publicKeyUri }} -o {{ options.downloadDir }}/{{ options.publicKeyFile }}",
      "curl {{ options.stemcellFileMd5Uri }} | tr -d '\"' > /opt/downloads/stemcellFileExpectedMd5",
      "curl {{ options.agentSettingsMd5Uri }} | tr -d '\"' > /opt/downloads/agentSettingsExpectedMd5",
      "curl {{ options.registrySettingsMd5Uri }} | tr -d '\"' > /opt/downloads/registrySettingsExpectedMd5",
      "curl {{ options.publicKeyMd5Uri }} | tr -d '\"' > /opt/downloads/publicKeyExpectedMd5",
      "md5sum {{ options.downloadDir }}/{{ options.stemcellFile }} | cut -d' ' -f1 > /opt/downloads/stemcellFileCalculatedMd5",
      "md5sum {{ options.downloadDir }}/{{ options.agentSettingsFile }} | cut -d' ' -f1 > /opt/downloads/agentSettingsCalculatedMd5",
      "md5sum {{ options.downloadDir }}/{{ options.registrySettingsFile }} | cut -d' ' -f1 > /opt/downloads/registrySettingsCalculatedMd5",
      "md5sum {{ options.downloadDir }}/{{ options.publicKeyFile }} | cut -d' ' -f1 > /opt/downloads/publicKeyCalculatedMd5",
      "test $(cat /opt/downloads/stemcellFileCalculatedMd5) = $(cat /opt/downloads/stemcellFileExpectedMd5)",
      "test $(cat /opt/downloads/agentSettingsCalculatedMd5) = $(cat /opt/downloads/agentSettingsExpectedMd5)",
      "test $(cat /opt/downloads/registrySettingsCalculatedMd5) = $(cat /opt/downloads/registrySettingsExpectedMd5)",
      "test $(cat /opt/downloads/publicKeyCalculatedMd5) = $(cat /opt/downloads/publicKeyExpectedMd5)",
      "sudo umount {{ options.device }} || true",
      "sudo tar --to-stdout -xvf {{ options.downloadDir }}/{{ options.stemcellFile }} | sudo dd of={{ options.device }}",
      "sudo sfdisk -R {{ options.device }}",
      "sudo mount {{ options.device }}1 /mnt",
      "sudo mkdir -p /mnt/home/vcap/.ssh",
      "sudo cat {{ options.downloadDir }}/{{ options.publicKeyFile }} >> /mnt/home/vcap/.ssh/authorized_keys",
      "sudo chown 1000:1000 /mnt/home/vcap/.ssh/authorized_keys",
      "sudo chown 1000:1000 /mnt/home/vcap/.ssh",
      "sudo chmod 600 /mnt/home/vcap/.ssh/authorized_keys",
      "sudo chmod 700 /mnt/home/vcap/.ssh/",
      "sudo cp {{ options.downloadDir }}/{{ options.agentSettingsFile }} /mnt/{{ options.agentSettingsPath }}",
      "sudo cp {{ options.downloadDir }}/{{ options.registrySettingsFile }} /mnt/{{ options.registrySettingsPath }}",
      "sudo sync"
    ],
    "device": "/dev/sda",
    "downloadDir": "/opt/downloads",
    "publicKeyFile": null,
    "publicKeyMd5Uri": "{{ api.files }}/md5/{{ options.publicKeyFile }}/latest",
    "publicKeyUri": "{{ api.files }}/{{ options.publicKeyFile }}/latest",
    "registrySettingsFile": null,
    "registrySettingsMd5Uri": "{{ api.files }}/md5/{{ options.registrySettingsFile }}/latest",
    "registrySettingsPath": null,
    "registrySettingsUri": "{{ api.files }}/{{ options.registrySettingsFile }}/latest",
    "stemcellFile": null,
    "stemcellFileMd5Uri": "{{ api.files }}/md5/{{ options.stemcellFile }}/latest",
    "stemcellUri": "{{ api.files }}/{{ options.stemcellFile }}/latest"
  },
  "properties": {}
}`)