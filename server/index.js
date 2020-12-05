import grpc from "grpc"
import * as protoLoader from "@grpc/proto-loader"

const PROTO_PATH = __dirname + "/../schema/square_service.proto"

const HOST = "0.0.0.0"
const PORT = 3000

const PackageDefintion = protoLoader.loadSync(
    PROTO_PATH, {
        keepCase: true,
        longs: String,
        enums: String,
        defaults: true,
        oneofs: true
    }
)

const rpc_proto = grpc.loadPackageDefinition(PackageDefintion).squarerpc_service

const SquareService = {
    square(call,callback){
        console.log(`get message ${call.request.message}`)
        let result = call.request.message**2
        callback(null,{
            message: result
        })
    }
}

const server = new grpc.Server()

server.addService(rpc_proto.SquareService.service,SquareService)

function main(){
    server.bind(`${HOST}:${PORT}`, grpc.ServerCredentials.createInsecure())
    console.log(`start @ ${HOST}:${PORT}`)
    server.start()
}

main()