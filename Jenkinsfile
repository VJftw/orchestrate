node {
    stage 'Commander: Unit tests'
    env.CI = "true"
    checkout scm
    wrap([$class: 'AnsiColorBuildWrapper', 'colorMapName': 'XTerm', 'defaultFg': 1, 'defaultBg': 2]) {
      sh '''
        set +x
        cd commander
        invoke test
      '''
    }
}
