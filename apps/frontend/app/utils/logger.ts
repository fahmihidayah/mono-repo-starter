/**
 * Logger utility for easier debugging
 *
 * Usage:
 * import { logger } from '~/utils/logger';
 *
 * logger.info('Component', 'Action description', { data });
 * logger.warn('Component', 'Warning message', { details });
 * logger.error('Component', 'Error occurred', error);
 */

type LogLevel = 'info' | 'warn' | 'error' | 'debug';

interface LogOptions {
  component: string;
  action?: string;
  data?: any;
  error?: any;
}

class Logger {
  private isDev = import.meta.env.DEV;
  private isEnabled = true;

  /**
   * Enable or disable logging
   */
  setEnabled(enabled: boolean) {
    this.isEnabled = enabled;
  }

  /**
   * Log info message
   */
  info(component: string, message: string, data?: any) {
    this.log('info', component, message, data);
  }

  /**
   * Log warning message
   */
  warn(component: string, message: string, data?: any) {
    this.log('warn', component, message, data);
  }

  /**
   * Log error message
   */
  error(component: string, message: string, error?: any) {
    this.log('error', component, message, error);
  }

  /**
   * Log debug message (only in development)
   */
  debug(component: string, message: string, data?: any) {
    if (this.isDev) {
      this.log('debug', component, message, data);
    }
  }

  /**
   * Internal log method
   */
  private log(level: LogLevel, component: string, message: string, data?: any) {
    if (!this.isEnabled) return;

    const timestamp = new Date().toISOString();
    const prefix = `[${component}]`;

    switch (level) {
      case 'info':
        console.log(`${prefix} ${message}`, data !== undefined ? data : '');
        break;
      case 'warn':
        console.warn(`${prefix} ${message}`, data !== undefined ? data : '');
        break;
      case 'error':
        console.error(`${prefix} ${message}`, data !== undefined ? data : '');
        break;
      case 'debug':
        console.debug(`${prefix} ${message}`, data !== undefined ? data : '');
        break;
    }
  }

  /**
   * Create a scoped logger for a specific component
   */
  scope(component: string) {
    return {
      info: (message: string, data?: any) => this.info(component, message, data),
      warn: (message: string, data?: any) => this.warn(component, message, data),
      error: (message: string, error?: any) => this.error(component, message, error),
      debug: (message: string, data?: any) => this.debug(component, message, data),
    };
  }

  /**
   * Log API request
   */
  apiRequest(method: string, url: string, data?: any) {
    this.debug('API', `${method} ${url}`, data);
  }

  /**
   * Log API response
   */
  apiResponse(method: string, url: string, status: number, data?: any) {
    const level = status >= 400 ? 'error' : status >= 200 && status < 300 ? 'info' : 'warn';
    this.log(level, 'API', `${method} ${url} → ${status}`, data);
  }
}

export const logger = new Logger();

// Create common component loggers
export const authLogger = logger.scope('Auth');
export const apiLogger = logger.scope('API');
export const routeLogger = logger.scope('Route');
