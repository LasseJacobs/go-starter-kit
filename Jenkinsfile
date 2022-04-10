pipeline {
    agent any
    environment {
        REGISTRY_HOST = "localhost:32000"

        APP_NAME = "gokit"
        GIT_DIGEST = "null"
    }
    stages {
        stage ('Calc digest') {
            steps {
                script {
                    GIT_DIGEST = sh(returnStdout: true, script: "git rev-parse --short HEAD").trim()
                }
            }
        }
        stage('Building Docker Image') {
            steps {
                echo "starting build ${env.BUILD_ID} - digest: $GIT_DIGEST"
		        sh "docker build -t $APP_NAME ."
            }
        }
        stage('Testing Docker Image') {
            steps {
                sh "docker-compose rm --force --stop -v"
                sh "docker-compose up --no-build -d"

                sleep(time:3,unit:"SECONDS")
                sh "make tests-e2e"

                sh "docker-compose down -v"
            }
        }
        stage('Push Docker Image') {
            steps {
                script {
                    if (GIT_DIGEST == 'null' || GIT_DIGEST == '') {
                        error('GIT_DIGEST variable is not set, this is required')
                    }
                }
                echo "pushing image to $REGISTRY_HOST/$APP_NAME:$GIT_DIGEST"

                sh "docker image tag $APP_NAME $REGISTRY_HOST/$APP_NAME:$GIT_DIGEST"
                sh "docker push $REGISTRY_HOST/$APP_NAME:$GIT_DIGEST"
            }
        }
        stage('Update K8 yaml files') {
            steps {
                sh "rm -r deployment/build || echo 'no directory' "
                sh "mkdir deployment/build"
                sh "cp deployment/*.yaml deployment/build/"
                sh "sed -i 's/\$GIT_TAG/$GIT_DIGEST/g' deployment/build/*.yaml"
            }
        }
        stage('Release') {
            steps {
                sshagent(credentials: ['github']) {
                    sh "git add ."
                    sh "git commit -m \"release $GIT_DIGEST\" "
                    sh "git push -f origin HEAD:release"
                }
            }
        }
    }
}