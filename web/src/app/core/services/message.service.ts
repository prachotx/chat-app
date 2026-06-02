import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';

export interface MessageResponse {
  id: number;
  created_at: string;
  updated_at: string;
  content: string;
  room_id: number;
  sender_id: number;
  sender: { id: number; username: string; email: string };
}

interface ApiResponse<T = null> {
  message: string;
  data?: T;
}

@Injectable({ providedIn: 'root' })
export class MessageService {
  private readonly http = inject(HttpClient);

  getByRoom(roomId: number) {
    return this.http.get<ApiResponse<MessageResponse[]>>(`/api/rooms/${roomId}/messages`, { withCredentials: true });
  }
}
