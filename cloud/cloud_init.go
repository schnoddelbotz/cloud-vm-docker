package cloud

import (
	"fmt"
	"strings"
)

func buildCloudInit(project, cfnRegion, image string, command []string) string {
	// should set vm shutdown token
	// should use task as first arg?
	// should quote all command parts
	my_command := strings.Join(command, " ")
	return fmt.Sprintf(`
#cloud-config
users:
- name: cloudservice
  uid: 2000

write_files:
- path: /etc/systemd/system/cloudservice.service
  permissions: 0644
  owner: root
  content: |
    [Unit]
    Description=Start inventory-optimisation docker container
    Wants=gcr-online.target
    After=gcr-online.target

    [Service]
    User=cloudservice
    Restart=no
    Environment="HOME=/home/cloudservice"
    ExecStartPre=/usr/bin/docker-credential-gcr configure-docker
    ExecStart=/usr/bin/docker run --rm \
        -v/var/run/docker.sock:/var/run/docker.sock \
        -v/home/cloudservice/.docker/config.json:/home/cloudservice/.docker/config.json \
        -eGCP_PROJECT=%s -eGCP_REGION=%s --name=ctzz %s %s
    ExecStop=/usr/bin/docker stop ctzz
    ExecStopPost=/usr/bin/docker rm ctzz

runcmd:
- usermod -aG docker cloudservice
- docker-credential-gcr configure-docker
- systemctl daemon-reload
- systemctl start cloudservice.service
`, project, cfnRegion, image, my_command)
}
