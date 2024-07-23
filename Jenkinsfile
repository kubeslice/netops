@Library('jenkins-library@opensource-release') _
dockerImagePipeline(
  script: this,
  service: 'netops',
  dockerfile: 'Dockerfile',
  buildContext: '.',
  buildArguments: [PLATFORM:"amd64"]
)
