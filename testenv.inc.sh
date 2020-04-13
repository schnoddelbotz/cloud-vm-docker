
# consumed by gcloud SDK library ... path to service account JSON
export GOOGLE_APPLICATION_CREDENTIALS=/Users/jan/.config/gcloud/jan-playground.json


# GoogleCloud: Project ID
export CTZZ_PROJECT=hacker-playground-254920
# used by Make for GCP PubSubFn deployment. Why not in env by default?
# https://cloud.google.com/functions/docs/env-var#nodejs_6_nodejs_8_python_37_and_go_111


./cloud-vm-docker completion -s fish | .
