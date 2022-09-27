@Library('jenkins-library@main') _
dockerImagePipeline(
  script: this,
  service: 'netops',
  dockerfile: 'Dockerfile',
  buildContext: '.',
  buildArguments: [PLATFORM:"amd64"]
)
