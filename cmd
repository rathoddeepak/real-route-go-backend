echo "\n"

cd account
protoc --proto_path=../thirdparty --proto_path=proto:. --micro_out=. --go_out=. account.proto
protoc --proto_path=../thirdparty --proto_path=proto:. --micro_out=. --go_out=. health.proto
echo "Account Proto..Done";

cd ../outservice
protoc --proto_path=../thirdparty --proto_path=proto:. --micro_out=. --go_out=. outservice.proto
protoc --proto_path=../thirdparty --proto_path=proto:. --micro_out=. --go_out=. health.proto
echo "Out Service Proto..Done";


cd ../session
protoc --proto_path=../thirdparty --proto_path=proto:. --micro_out=. --go_out=. session.proto
protoc --proto_path=../thirdparty --proto_path=proto:. --micro_out=. --go_out=. health.proto
echo "Session Proto..Done";


cd ../city
protoc --proto_path=../thirdparty --proto_path=proto:. --micro_out=. --go_out=. city.proto
protoc --proto_path=../thirdparty --proto_path=proto:. --micro_out=. --go_out=. health.proto
echo "City Proto..Done";

cd ../hub
protoc --proto_path=../thirdparty --proto_path=proto:. --micro_out=. --go_out=. hub.proto
protoc --proto_path=../thirdparty --proto_path=proto:. --micro_out=. --go_out=. health.proto
echo "Hub Proto..Done";

cd ../logistics
protoc --proto_path=../thirdparty --proto_path=protofiles:. --micro_out=. --go_out=. logistics.proto
protoc --proto_path=../thirdparty --proto_path=protofiles:. --micro_out=. --go_out=. accountservice.proto
protoc --proto_path=../thirdparty --proto_path=protofiles:. --micro_out=. --go_out=. hubservice.proto
protoc --proto_path=../thirdparty --proto_path=protofiles:. --micro_out=. --go_out=. cityservice.proto
protoc --proto_path=../thirdparty --proto_path=protofiles:. --micro_out=. --go_out=. outservice.proto
protoc --proto_path=../thirdparty --proto_path=protofiles:. --micro_out=. --go_out=. health.proto
echo "Logistics Proto..Done";



echo "\n\n"
echo "-----------------Gateway-------------"
echo "\n\n"

cd ../proto
protoc --plugin=protoc-gen-grpc-gateway=/Users/Clufter/go/bin/protoc-gen-grpc-gateway --proto_path=../thirdparty --proto_path=account:. --grpc-gateway_out=account --grpc-gateway_opt=paths=source_relative account.proto
protoc --proto_path=../thirdparty --proto_path=account:.  --go_out=account --go_opt=paths=source_relative --go-grpc_out=account --go-grpc_opt=paths=source_relative account.proto
echo "GateWay Account Proto..Done";


protoc --plugin=protoc-gen-grpc-gateway=/Users/Clufter/go/bin/protoc-gen-grpc-gateway --proto_path=../thirdparty --proto_path=session:. --grpc-gateway_out=session --grpc-gateway_opt=paths=source_relative session.proto
protoc --proto_path=../thirdparty --proto_path=session:.  --go_out=session --go_opt=paths=source_relative --go-grpc_out=session --go-grpc_opt=paths=source_relative session.proto
echo "GateWay Session Proto..Done";

protoc --plugin=protoc-gen-grpc-gateway=/Users/Clufter/go/bin/protoc-gen-grpc-gateway --proto_path=../thirdparty --proto_path=city:. --grpc-gateway_out=city --grpc-gateway_opt=paths=source_relative city.proto
protoc --proto_path=../thirdparty --proto_path=city:.  --go_out=city --go_opt=paths=source_relative --go-grpc_out=city --go-grpc_opt=paths=source_relative city.proto
echo "GateWay city Proto..Done";

protoc --plugin=protoc-gen-grpc-gateway=/Users/Clufter/go/bin/protoc-gen-grpc-gateway --proto_path=../thirdparty --proto_path=hub:. --grpc-gateway_out=hub --grpc-gateway_opt=paths=source_relative hub.proto
protoc --proto_path=../thirdparty --proto_path=hub:.  --go_out=hub --go_opt=paths=source_relative --go-grpc_out=hub --go-grpc_opt=paths=source_relative hub.proto
echo "GateWay hub Proto..Done";

protoc --plugin=protoc-gen-grpc-gateway=/Users/Clufter/go/bin/protoc-gen-grpc-gateway --proto_path=../thirdparty --proto_path=logistics:. --grpc-gateway_out=logistics --grpc-gateway_opt=paths=source_relative logistics.proto
protoc --proto_path=../thirdparty --proto_path=logistics:.  --go_out=logistics --go_opt=paths=source_relative --go-grpc_out=logistics --go-grpc_opt=paths=source_relative logistics.proto
echo "GateWay Logistics Proto..Done";

echo "\n\n"