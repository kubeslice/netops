@Library('jenkins-library@opensource') _
dockerImagePipeline(
  script: this,
  service: 'aveshadev/netops',
  dockerfile: 'Dockerfile',
  buildContext: '.',
  buildArguments: [PLATFORM:"amd64"]
)
