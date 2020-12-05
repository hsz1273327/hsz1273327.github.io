import grpc from "grpc"
import * as protoLoader from "@grpc/proto-loader"
import * as bluebird from 'bluebird'


const PROTO_PATH = __dirname + "/../schema/square_service.proto"

const HOST = "localhost"
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

const clientcb = new rpc_proto.SquareService(
    `${HOST}:${PORT}`,
    grpc.credentials.createInsecure()
)

const client = bluebird.promisifyAll(clientcb)

async function main() {
    let result = await client.squareAsync({
        message: 12.3
    })
    console.log(result.message)
}

main()