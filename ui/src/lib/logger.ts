/**
 * Logger utility for debugging
 * Provides consistent logging across the application
 */

type LogLevel = 'debug' | 'info' | 'warn' | 'error';

const isDevelopment = process.env.NODE_ENV === 'development';

class Logger {
  private prefix: string;

  constructor(prefix: string) {
    this.prefix = prefix;
  }

  private formatMessage(level: LogLevel, message: string, data?: any): string {
    const timestamp = new Date().toISOString();
    return `[${timestamp}] [${level.toUpperCase()}] [${this.prefix}] ${message}`;
  }

  debug(message: string, data?: any) {
    if (isDevelopment) {
      console.log(this.formatMessage('debug', message), data || '');
    }
  }

  info(message: string, data?: any) {
    console.log(this.formatMessage('info', message), data || '');
  }

  warn(message: string, data?: any) {
    console.warn(this.formatMessage('warn', message), data || '');
  }

  error(message: string, error?: any) {
    console.error(this.formatMessage('error', message), error || '');
  }

  group(label: string) {
    if (isDevelopment) {
      console.group(`[${this.prefix}] ${label}`);
    }
  }

  groupEnd() {
    if (isDevelopment) {
      console.groupEnd();
    }
  }
}

// Factory function to create loggers for different parts of the app
export function createLogger(prefix: string): Logger {
  return new Logger(prefix);
}

// Pre-configured loggers
export const authLogger = createLogger('AUTH');
export const apiLogger = createLogger('API');
export const middlewareLogger = createLogger('MIDDLEWARE');
export const supabaseLogger = createLogger('SUPABASE');
export const storeLogger = createLogger('STORE');
