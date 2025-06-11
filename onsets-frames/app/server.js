const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const path = require('path');
const config = require('./config');
const mm = require('@magenta/music');
const createAnalyzer = require('./controllers/analyze');

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
    const model = new mm.OnsetsAndFrames();
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
