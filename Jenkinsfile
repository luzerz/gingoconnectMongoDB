import groovy.json.JsonSlurperClassic
import org.yaml.snakeyaml.Yaml
import org.yaml.snakeyaml.DumperOptions
import groovy.json.JsonSlurper
pipeline {
    agent any
    stages {
        stage('Git & ENV') {
            steps {
                script
                {
                       def properties = readYaml file: "jenkins/properties.yaml"
                       Branch_name = "${GIT_BRANCH.split('/').size() > 1 ? GIT_BRANCH.split('/')[1] : GIT_BRANCH}"
                       Tag = "${GIT_BRANCH.split('/').size() > 2 ? GIT_BRANCH.split('/')[2] : GIT_BRANCH}"
                       println("BRANCH : ${Branch_name}")
                       if(Branch_name == 'main'){
                          App_env = 'beta'
                          DeployAppNamespace         = properties.DeployApp.NamespaceDev
                          DeployAppVersionTag       = 'latest'
                          Prefix_Tag = ""
                       }else if(Branch_name.contains('release')){
                          println("TAG RELEASE: ${Tag}")
                          App_env     = 'sit'
                          App_env_uat = 'uat'
                          App_env_stg = 'stg'
                          DeployAppNamespace         = properties.DeployApp.NamespaceSIT
                          DeployAppNamespaceUAT     = properties.DeployApp.NamespaceUAT
                          DeployAppNamespaceSTG     = properties.DeployApp.NamespaceSTG
                          DeployAppVersionTag       = Tag
                          Prefix_Tag = "release-v."
                       }else if(Branch_name == 'master'){
                          App_env     = 'production'
                          DeployAppNamespace         = properties.DeployApp.NameSpacePRD
                          //EDIT FOR PRODUCTION BUILD
                          DeployApp_Version_tag       = 'v1.0.0'
                          Prefix_Tag = "release-v."
                       }
                      DeployModeUpdateImage           = properties.DeployMode.UpdateImage
                      DeployModeConfigmap             = properties.DeployMode.Configmap
                      DeployAppProjectName            = properties.DeployApp.ProjectName
                      DeployAppServiceName            = properties.DeployApp.AppServiceName
                      DeployAppLanguage               = properties.DeployApp.Language
                      RegistryUrlRepository           = properties.Registry.UrlRepository
                      RegistryGCRCredential           = properties.Registry.GCRCredential
                      SonarScan                       = properties.Security.SonarQube.Scan
                      SonarProjectKey                 = properties.Security.SonarQube.ProjectKey
                      SonarSource                     = properties.Security.SonarQube.Sources
                      SonarServer                     = properties.Security.SonarQube.Host
                      SonarQuantityGateUrl            = SonarServer+'/api/qualitygates/project_status?projectKey='+SonarProjectKey
                      SonarLogin                      = properties.Security.SonarQube.Login
                      SonarTestsSource                = properties.Security.SonarQube.Tests
                      SonarTestInclusions             = properties.Security.SonarQube.TestInclusions
                      SonarExclusions                 = properties.Security.SonarQube.Exclusions
                      SecurityScanImages              = properties.Security.ScanImages
                      SecurityScanDAST                = properties.Security.ScanDAST
                      TestingUnittest                 = properties.Testing.Unittest
                      TestingIntegrateTest            = properties.Testing.IntegrateTest
                      TestingQA                       = properties.Testing.Qa

                      println (
                      'DeployModeUpdateImage          = ' + DeployModeUpdateImage +'\n'+
                      'DeployModeConfigmap            = ' + DeployModeConfigmap +'\n'+
                      'DeployAppProjectName           = ' + DeployAppProjectName +'\n'+
                      'DeployAppServiceName           = ' + DeployAppServiceName +'\n'+
                      'DeployAppVersionTag            = ' + DeployAppVersionTag +'\n'+
                      'DeployAppLanguage              = ' + DeployAppLanguage +'\n'+
                      'RegistryUrlRepository          = ' + RegistryUrlRepository +'\n'+
                      'RegistryGCRCredential          = ' + RegistryGCRCredential +'\n'+
                      'SonarProjectKey                = ' + SonarProjectKey +'\n'+
                      'SonarServer                    = ' + SonarServer +'\n'+
                      'SonarLogin                     = ' + SonarLogin +'\n'+
                      'SecuritySonarqube              = ' + SonarScan +'\n'+
                      'SecurityScanImages             = ' + SecurityScanImages +'\n'+
                      'SecurityScanDAST               = ' + SecurityScanDAST +'\n'+
                      'TestingUnittest                = ' + TestingUnittest +'\n'+
                      'TestingIntegrateTest           = ' + TestingIntegrateTest +'\n'+
                      'TestingQA                      = ' + TestingQA +'\n');
                }
            }
        }
       
        stage('Build images and scan') {
            steps
            {
                script
                    {
                   // sh 'CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .'
                    a = docker.build("${RegistryUrlRepository}/${DeployAppProjectName}/${DeployAppServiceName}:${DeployAppVersionTag}", "-f ./Dockerfile .")
                    if (SecurityScanImages){
                      sh '''
                          /usr/local/bin/trivy --exit-code 1 --severity CRITICAL '''+"${RegistryUrlRepository}/${DeployAppProjectName}/${DeployAppServiceName}:${DeployAppVersionTag}"+'''
                          my_exit_code=$?
                          echo "RESULT 1:--- $my_exit_code"
                          if [[ ${my_exit_code} == 1 ]]; then
                              echo "Image scanning failed. Some vulnerabilities found"
                              exit 1;
                          else
                              echo "Image is scanned Successfully. No vulnerabilities found"
                          fi;
                       '''
                    }

                }
            }
        }

        stage('Push images') {
            steps
            {
                    script
                    {
                        withCredentials([file(credentialsId: 'jobfinder-secret', variable: 'GC_KEY')]){
                            sh "cat '$GC_KEY' | docker login -u _json_key --password-stdin https://asia.gcr.io"
                            sh "gcloud auth activate-service-account --key-file='$GC_KEY'"
                            sh "gcloud auth configure-docker"
                            sh (script: 'gcloud auth print-access-token',returnStdout: true).trim()
                            echo "Pushing image To GCR"
                            a.push()
                        }
                    }
                }
        }

        stage('Deploy to Cluster') {
            steps
            {
                script {
                    withKubeConfig([credentialsId:'job-cluster-secret']){
                        if(DeployModeConfigmap){
                                sh "kubectl apply -R -f ${env.WORKSPACE}/configk8s/configmap/*.yml -n ${DeployAppNamespace}"
                                sh "kubectl apply -R -f ${env.WORKSPACE}/configk8s/deployment/*.yml -n ${DeployAppNamespace}"
                                sh "kubectl apply -R -f ${env.WORKSPACE}/configk8s/service/*.yml -n ${DeployAppNamespace}"
                                //sh "kubectl apply -R -f ${env.WORKSPACE}/configk8s/kong/*.yaml -n ${DeployAppNamespace}"
                        }
                        if(DeployModeUpdateImage) {
                                deploy     = sh(returnStdout: true, script:"kubectl set image deployment/${DeployAppServiceName}  ${DeployAppServiceName}=${RegistryUrlRepository}/${DeployAppProjectName}/${DeployAppServiceName}:${DeployAppVersionTag} --record -n ${DeployAppNamespace} || true").trim()
                                status     = sh(returnStdout: true, script:"kubectl rollout status deploy ${DeployAppServiceName} -n ${DeployAppNamespace} --timeout=60s || true").trim()
                                patch      = sh(returnStdout: true, script:"kubectl set env deployment --env=\"LAST_PIPELINE_DEPLOY=\$(date +%s)\" ${DeployAppServiceName} -n "+ DeployAppNamespace).trim()
                                podName    = sh(returnStdout: true, script:"kubectl get pod -n ${DeployAppNamespace} | grep ${DeployAppServiceName} | head -n 1 | awk '{print \$1}'").trim()
                                describePo = sh(returnStdout: true, script:"kubectl get event -n ${DeployAppNamespace} --field-selector involvedObject.name=$podName").trim()
                                allStatus  = sh(returnStdout: true, script:"kubectl get pod -n ${DeployAppNamespace}").trim()

                                echo "$deploy  $status $podName $describePo $allStatus"
                                boolean running = status.contains("successfully rolled out")
                                if (running == true){
                                        echo "\nStatus all pods is \n $allStatus"
                                        currentBuild.result = 'SUCCESS'
                                }
                                else if ( running == false) {
                                    error('Deploy false ')
                                }

                        }
                    }
                }
            }
        }

        // post {
        //     always {
        //         echo 'One way or another, I have finished'
        //         sh "docker rmi -f ${RegistryUrlRepository}/${DeployAppProjectName}/${DeployAppServiceName}:${DeployAppVersionTag}  || echo 'ignore' "
        //         cleanWs()
        //     }
        //     success {
        //         echo 'Succeeded!'
        //     }
        //     failure {
        //         echo 'Failed!!!'
        //     }
        // }
    }
}
