import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';

export interface Room {
  id: number;
  name: string;
  created_by_id: number;
  created_by: { id: number; username: string };
}

interface ApiResponse<T = null> {
  message: string;
  data?: T;
}

@Injectable({ providedIn: 'root' })
export class RoomService {
  private readonly http = inject(HttpClient);

  getAll() {
    return this.http.get<ApiResponse<Room[]>>('/api/rooms', { withCredentials: true });
  }

  create(name: string) {
    return this.http.post<ApiResponse<Room>>('/api/rooms', { name }, { withCredentials: true });
  }
}
