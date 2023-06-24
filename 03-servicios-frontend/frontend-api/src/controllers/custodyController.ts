import { Request, Response } from "express";
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';

const PROTO_PATH = __dirname + '../../../protos/system/custody/custody.proto';

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
    keepCase: true,
	defaults: true,
	oneofs: true
});

const protoDescriptor = grpc.loadPackageDefinition(packageDefinition);

const endpoint = process.env.CUSTODY_BACKED || 'localhost:5000';

const creds = grpc.credentials.createInsecure();
const service = (protoDescriptor.lab as any).system.custody.CustodyService;
let stub = new service(endpoint, creds);

export const create = async (req: Request, res: Response) => {
    const msg = req.body;

    const now = Math.floor(Date.now()/1000);

    const p = new Promise((resolve, reject) =>
        stub.AddCustodyStock({
            period: msg.period,
            stock: msg.stock,
            client_id: msg.client_id,
            quantity: msg.quantity
        }, (err: any, response: any) => {
            if (err)
                return reject(err);
            resolve(response);
        })
    );

    const result = await p;

    return res.json(result);
}


export const get = async (req: Request, res: Response) => {
    const msg = req.body;

    //const now = Math.floor(Date.now()/1000);

    const p = new Promise((resolve, reject) =>
        stub.GetCustody({
            period:        	msg.period,
		    client_id:      	msg.client_id,
		    stock: 			msg.stock
            
        }, (err: any, response: any) => {
            if (err)
                return reject(err);
            resolve(response);
        })
    );

    const result = await p;

    return res.json(result);
}