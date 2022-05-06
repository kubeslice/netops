@Library('jenkins-library@kubeslice-mesh') _
dockerImagePipeline(
  script: this,
  service: 'aveshadev/netops',
  dockerfile: 'Dockerfile',
  buildContext: '.',
  buildArguments: [PLATFORM:"amd64"]
)
