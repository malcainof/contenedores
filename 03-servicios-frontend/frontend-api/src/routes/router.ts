import { Express, NextFunction, Router, Request, Response } from "express";
import * as billingController from '../controllers/billingController';
import * as custodyController from '../controllers/custodyController';

import * as commons from '../commons/commons';

const logger = commons.getLogger();

export const setup = (root: string, app: Express) => {

    const asyncMiddleware = (fn: Function) =>
        (req: Request, res: Response, next: NextFunction) => {
            Promise.resolve(fn(req, res, next)).catch(err => next(err));
        }

    const router = Router();

    router.get('/', (req, res, next) => {
        res.send('API Lab v1.0.0');
    });

    router.post('/billing/create', asyncMiddleware(billingController.create));
    router.post('/custody/create', asyncMiddleware(custodyController.create));
    router.post('/custody/get', asyncMiddleware(custodyController.get));

    router.use((err: Error, req: Request, res: Response, next: NextFunction) => {
        logger.error(err);
        res.status(500).json(err.message);
    });

    app.use(root, router);
}