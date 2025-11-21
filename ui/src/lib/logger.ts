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

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  debug(_message: string, _data?: any) {
    // Disabled - only warn and error logs are enabled
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  info(_message: string, _data?: any) {
    // Disabled - only warn and error logs are enabled
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
