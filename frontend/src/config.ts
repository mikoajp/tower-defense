// Centralized frontend configuration
export const API_URL: string = (import.meta as any).env?.VITE_API_URL ?? 'http://localhost:8080';
export const WS_URL: string = API_URL.replace(/^http(s?):\/\//, 'ws$1://') + '/ws';
