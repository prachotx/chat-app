import { ChangeDetectionStrategy, Component, OnDestroy, OnInit, inject, signal } from '@angular/core';
import { JsonPipe } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { MessageService, MessageResponse } from '../../core/services/message.service';

type WsEventType = 'message' | 'user_joined' | 'user_left' | 'online_users';

interface WsEvent<T = unknown> {
  type: WsEventType;
  data?: T;
}

interface UserStatusData {
  user_id: number;
}

interface OnlineUsersData {
  user_ids: number[];
}

@Component({
  selector: 'app-room',
  imports: [FormsModule, JsonPipe],
  templateUrl: './room.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class RoomComponent implements OnInit, OnDestroy {
  private readonly route = inject(ActivatedRoute);
  private readonly messageService = inject(MessageService);

  private ws: WebSocket | null = null;
  private roomId!: number;

  readonly messages = signal<MessageResponse[]>([]);
  readonly onlineUserIds = signal<number[]>([]);
  readonly isConnected = signal(false);
  readonly isLoadingHistory = signal(true);

  draft = '';

  ngOnInit() {
    this.roomId = Number(this.route.snapshot.paramMap.get('id'));
    this.loadHistory();
    this.connectWs();
  }

  ngOnDestroy() {
    this.ws?.close();
  }

  private loadHistory() {
    this.isLoadingHistory.set(true);
    this.messageService.getByRoom(this.roomId).subscribe({
      next: res => {
        this.messages.set(res.data ?? []);
        this.isLoadingHistory.set(false);
      },
      error: () => this.isLoadingHistory.set(false),
    });
  }

  private connectWs() {
    const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:';
    this.ws = new WebSocket(`${protocol}//${location.host}/ws/${this.roomId}`);

    this.ws.onopen = () => this.isConnected.set(true);
    this.ws.onclose = () => this.isConnected.set(false);
    this.ws.onerror = () => this.isConnected.set(false);
    this.ws.onmessage = ({ data }: MessageEvent<string>) =>
      this.handleEvent(JSON.parse(data) as WsEvent);
  }

  private handleEvent(event: WsEvent) {
    switch (event.type) {
      case 'message':
        this.messages.update(msgs => [...msgs, event.data as MessageResponse]);
        break;
      case 'online_users':
        this.onlineUserIds.set((event.data as OnlineUsersData).user_ids);
        break;
      case 'user_joined': {
        const userId = (event.data as UserStatusData).user_id;
        this.onlineUserIds.update(ids => (ids.includes(userId) ? ids : [...ids, userId]));
        break;
      }
      case 'user_left': {
        const userId = (event.data as UserStatusData).user_id;
        this.onlineUserIds.update(ids => ids.filter(id => id !== userId));
        break;
      }
    }
  }

  sendMessage() {
    const content = this.draft.trim();
    if (!content || !this.ws || this.ws.readyState !== WebSocket.OPEN) return;
    this.ws.send(JSON.stringify({ content }));
    this.draft = '';
  }
}
