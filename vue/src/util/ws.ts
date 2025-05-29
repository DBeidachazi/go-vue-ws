export class WebSocketClient {
    private socket: WebSocket | null = null;
    private url: string;
    private listeners: { [key: string]: Function[] } = {};

    /**
     * 创建WebSocket客户端
     * @param url WebSocket服务器地址
     */
    constructor(url: string) {
        this.url = url;
    }

    connect(): void {
        if (this.socket?.readyState === WebSocket.OPEN) return;

        try {
            const wsUrl = new URL(this.url);

            this.socket = new WebSocket(wsUrl.toString());

            this.socket.onopen = () => {
                console.log('WebSocket连接已建立');
                this.emit('open');
            };

            this.socket.onmessage = (event) => {
                try {
                    const data = JSON.parse(event.data);
                    this.emit('message', data);
                } catch (e) {
                    this.emit('message', event.data);
                }
            };

            this.socket.onclose = (event) => {
                this.emit('close', event);
            };

            this.socket.onerror = (error) => {
                console.error('WebSocket错误:', error);
                this.emit('error', error);
            };
        } catch (error) {
            console.error('WebSocket连接失败:', error);
        }
    }


    close(): void {
        if (this.socket) {
            this.socket.close();
            this.socket = null;
        }
    }

    on(event: string, callback: Function): void {
        if (!this.listeners[event]) {
            this.listeners[event] = [];
        }
        this.listeners[event].push(callback);
    }

    off(event: string, callback: Function): void {
        if (!this.listeners[event]) return;

        const index = this.listeners[event].indexOf(callback);
        if (index !== -1) {
            this.listeners[event].splice(index, 1);
        }
    }

    private emit(event: string, ...args: any[]): void {
        if (!this.listeners[event]) return;

        this.listeners[event].forEach(callback => {
            callback(...args);
        });
    }
}