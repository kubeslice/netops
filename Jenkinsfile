@Library('jenkins-library@main') _
dockerImagePipelineos(
  script: this,
  service: 'netops',
  dockerfile: 'Dockerfile',
  buildContext: '.',
  buildArguments: [PLATFORM:"amd64"]
)
