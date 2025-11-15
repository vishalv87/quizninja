/**
 * Global TypeScript type declarations
 * This file tells TypeScript how to handle various file imports
 */

// CSS Modules
declare module '*.css' {
  const content: { [className: string]: string };
  export default content;
}

declare module '*.scss' {
  const content: { [className: string]: string };
  export default content;
}

// Image assets
declare module '*.png';
declare module '*.jpg';
declare module '*.jpeg';
declare module '*.gif';
declare module '*.svg';
declare module '*.webp';
declare module '*.ico';

// Other assets
declare module '*.json';
declare module '*.woff';
declare module '*.woff2';
declare module '*.ttf';
declare module '*.eot';
