// Minimal SDK contract used to type-check generated HTTP and WebSocket connectors.
// index.js implements the runtime fake. This file describes the public signatures only;
// payload/result/handler are test controls and intentionally are not part of that contract.
export interface Config {
	errorHandler: (error: Error) => string;
}
export interface RequestOptions {}
export interface WebSocketOptions {}
export class Ctx {}
export class UserError extends Error {}
export class Client {
	constructor(config: Config);
	getConfig(): Config;
	receive<T>(name: string, handler: (ctx: Ctx, event: T) => void): void;
	send<R, S>(
		name: string,
		payload: R,
		options?: RequestOptions,
	): Promise<S | UserError>;
	sendOk<R>(
		name: string,
		payload: R,
		options?: RequestOptions,
	): Promise<UserError | undefined>;
	sendOkNoRequest(
		name: string,
		options?: RequestOptions,
	): Promise<UserError | undefined>;
	sendNoRequest<S>(
		name: string,
		options?: RequestOptions,
	): Promise<S | UserError>;
	sendNoResponse<R>(
		name: string,
		payload: R,
		options?: RequestOptions,
	): Promise<void>;
	sendPing(name: string, options?: RequestOptions): Promise<void>;
}
export function applyHTTP(client: Client, method: string, url: string): void;

// WebSocket output is compiled against this declaration, but not instantiated by the tests.
export class WebSocketTransporter {
	constructor(client: Client, config: Config & WebSocketOptions);
}
