# -*- mode: Python -*-

# For more on Extensions, see: https://docs.tilt.dev/extensions.html
load('ext://restart_process', 'docker_build_with_restart')

local_resource(
  'compile',
  'make build GOOS=linux CGO_ENABLED=0',
  deps=['./cmd', './pkg', './internal'],
)

docker_build_with_restart(
  'gcr.io/outreach-docker/ci',
  '.',
  entrypoint=['/app/bin/ci'],
  dockerfile='deployments/ci/Dockerfile.dev',
  only=[
    './bin',
    './deployments/ci',
  ],
  ssh='default',
  live_update=[
    sync('./bin', '/app/bin'),
  ],
)

templated_yaml = local('./scripts/shell-wrapper.sh deploy-to-dev.sh show')
k8s_yaml(templated_yaml)
k8s_resource('ci', port_forwards=8080,
             resource_deps=['compile'])
