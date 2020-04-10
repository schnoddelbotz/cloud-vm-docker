package cloud

import "fmt"

func buildCloudInit(project, cfnRegion, image, command string) string {
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
    Restart=always
    RestartSec=5
    Environment="HOME=/home/cloudservice"
    ExecStartPre=/usr/bin/docker-credential-gcr configure-docker
    ExecStartPre=/usr/bin/docker pull %s
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
`, image, project, cfnRegion, image, command)
}
