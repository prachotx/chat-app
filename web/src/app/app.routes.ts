import { Routes } from '@angular/router';
import { authGuard } from './core/guards/auth.guard';

export const routes: Routes = [
  { path: '', redirectTo: '/home', pathMatch: 'full' },
  {
    path: 'login',
    loadComponent: () => import('./pages/login/login').then(m => m.LoginComponent),
  },
  {
    path: 'register',
    loadComponent: () => import('./pages/register/register').then(m => m.RegisterComponent),
  },
  {
    path: 'rooms',
    canActivate: [authGuard],
    loadComponent: () => import('./pages/rooms/rooms').then(m => m.RoomsComponent),
  },
  {
    path: 'rooms/:id',
    canActivate: [authGuard],
    loadComponent: () => import('./pages/room/room').then(m => m.RoomComponent),
  },
];
