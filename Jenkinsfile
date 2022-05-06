@Library('jenkins-library@opensource') _
dockerImagePipeline(
  script: this,
  service: 'aveshasystems/netops',
  dockerfile: 'Dockerfile',
  buildContext: '.',
  buildArguments: [PLATFORM:"amd64"]
)
