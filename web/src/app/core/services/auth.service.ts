import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { tap, catchError, of } from 'rxjs';

export interface User {
  id: number;
  username: string;
  email: string;
}

interface ApiResponse<T = null> {
  message: string;
  data?: T;
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly http = inject(HttpClient);

  readonly currentUser = signal<User | null>(null);

  login(email: string, password: string) {
    return this.http.post<ApiResponse>('/api/auth/login', { email, password }, { withCredentials: true });
  }

  register(username: string, email: string, password: string) {
    return this.http.post<ApiResponse>('/api/auth/register', { username, email, password }, { withCredentials: true });
  }

  me() {
    return this.http.get<ApiResponse<User>>('/api/auth/me', { withCredentials: true }).pipe(
      tap(res => this.currentUser.set(res.data ?? null)),
      catchError(() => {
        this.currentUser.set(null);
        return of(null);
      }),
    );
  }

  logout() {
    return this.http.get<ApiResponse>('/api/auth/logout', { withCredentials: true }).pipe(
      tap(() => this.currentUser.set(null)),
    );
  }
}
