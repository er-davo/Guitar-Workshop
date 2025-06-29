import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import path from 'path';
import config from './config.js';
import * as magenta from '@magenta/music';
import createAnalyzer from './controllers/analyze.js'; 

const PROTO_PATH = path.join(__dirname, 'proto/onsetsframes.proto');

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
});

const protoDescriptor = grpc.loadPackageDefinition(packageDefinition);
const onsetsframes = protoDescriptor.onsetsframes;

async function main() {
    const model = new magenta.OnsetsAndFrames();
    await model.initialize();
    console.log('O&F модель загружена');

    const server = new grpc.Server();

    server.addService(onsetsframes.OnsetsAndFrames.service, {
        Analyze: createAnalyzer(model),
    });

    server.bindAsync(`0.0.0.0:${config.PORT}`, grpc.ServerCredentials.createInsecure(), (err, port) => {
        if (err) throw err;
        console.log(`gRPC сервер работает на порту ${port}`);
    });
}

main().catch(err => {
    console.error('Ошибка запуска сервера:', err);
});
