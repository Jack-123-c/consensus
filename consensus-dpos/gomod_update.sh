set -x
BRANCH=v2.3.0_qc
go get chainmaker.org/chainmaker/consensus-utils/v2@${BRANCH}
go get chainmaker.org/chainmaker/logger/v2@${BRANCH}
go get chainmaker.org/chainmaker/pb-go/v2@${BRANCH}
go get chainmaker.org/chainmaker/protocol/v2@${BRANCH}
go get chainmaker.org/chainmaker/utils/v2@${BRANCH}
go get chainmaker.org/chainmaker/vm-native/v2@${BRANCH}
go mod tidy
go build ./...
