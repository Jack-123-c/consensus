module chainmaker.org/chainmaker/consensus-dpos/v2

go 1.15

require (
	chainmaker.org/chainmaker/consensus-utils/v2 v2.3.0
	chainmaker.org/chainmaker/logger/v2 v2.3.0
	chainmaker.org/chainmaker/pb-go/v2 v2.3.0
	chainmaker.org/chainmaker/protocol/v2 v2.3.0
	chainmaker.org/chainmaker/utils/v2 v2.3.0
	chainmaker.org/chainmaker/vm-native/v2 v2.3.0
	github.com/gogo/protobuf v1.3.2
	github.com/golang/mock v1.6.0
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.1-0.20200815110645-5c35d600f0ca
)

replace (
	github.com/linvon/cuckoo-filter => chainmaker.org/third_party/cuckoo-filter v1.0.0
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
