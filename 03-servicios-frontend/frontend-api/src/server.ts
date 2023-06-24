import express from "express";
import http from 'http';
import cors from 'cors';
import * as commons from './commons/commons';

const logger = commons.getLogger();
logger.info('Inicializando servidor');

const app = express();

app.options('*', cors() as any);
app.use(cors());

app.use(express.json({ limit: '50mb' }));
app.use(express.urlencoded({
    extended: false,
    limit: '50mb'
}));

import * as router from './routes/router';
router.setup('/', app);

const port = process.env.PORT || '3003';
const host = process.env.HOST || '0.0.0.0';
const server = http.createServer(app).listen(Number(port), host);

process.on('SIGTERM', () => {
    server.close(() => {
        logger.info('Servidor finalizado');
    });
});
process.on('SIGINT', () => {
    server.close(() => {
        logger.info('Servidor finalizado');
    });
});

logger.info('Servidor inicializado en %s:%s', host, port);
server.timeout = 240000;