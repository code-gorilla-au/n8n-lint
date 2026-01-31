import pino from 'pino';

const isCI = process.env.CI === 'true';


export const logger = pino({
	transport: {
		target: 'pino-pretty',
		options: { colorize: !isCI }
	}
});
