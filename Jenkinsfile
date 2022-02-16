@Library('jenkins-library@kubeslice-mesh') _
dockerImagePipeline(
  script: this,
  service: 'kubeslice-netops',
  dockerfile: 'Dockerfile',
  buildContext: '.',
  buildArguments: [PLATFORM:"amd64"]
)
