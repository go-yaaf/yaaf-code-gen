// This implementation is based on Medium articles: WebSockets in Angular: A Comprehensive Guide
// https://medium.com/@saranipeiris17/websockets-in-angular-a-comprehensive-guide-e92ca33f5d67
// https://medium.com/@saranipeiris17/websockets-in-angular-a-comprehensive-guide-part-2-bd8021a9be09

import { Injectable } from '@angular/core';
import { WebSocketSubject, webSocket } from 'rxjs/webSocket';
import { Observable, timer } from 'rxjs';
import { switchMap, retryWhen } from 'rxjs/operators';

@Injectable({
    providedIn: 'root'
})
export class WebSocketService {
    private socket$: WebSocketSubject<any>;
    private reconnectInterval = 5000; // 5 seconds

    constructor() {
        this.connect();
    }

    private connect() {
        this.socket$ = webSocket('ws://your-websocket-url');

        this.socket$.pipe(
            retryWhen(errors =>
                errors.pipe(
                    switchMap(() => {
                        console.log(`WebSocket connection failed. Retrying in ${this.reconnectInterval / 1000} seconds...`);
                        return timer(this.reconnectInterval);
                    })
                )
            )
        ).subscribe(
            message => this.handleMessage(message),
            error => console.error('WebSocket error:', error),
            () => console.log('WebSocket connection closed')
        );
    }

    // Send a message to the server
    sendMessage(message: any) {
        this.socket$.next(message);
    }

    // Receive messages from the server
    getMessages(): Observable<any> {
        return this.socket$.asObservable();
    }

    // Close the WebSocket connection
    closeConnection() {
        this.socket$.complete();
    }
}