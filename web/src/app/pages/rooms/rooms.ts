import { ChangeDetectionStrategy, Component, OnInit, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../../core/services/auth.service';
import { Room, RoomService } from '../../core/services/room.service';

@Component({
  selector: 'app-rooms',
  imports: [ReactiveFormsModule, RouterLink],
  templateUrl: './rooms.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class RoomsComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly authService = inject(AuthService);
  private readonly roomService = inject(RoomService);
  private readonly router = inject(Router);

  readonly currentUser = this.authService.currentUser;
  readonly rooms = signal<Room[]>([]);
  readonly isLoadingRooms = signal(true);
  readonly showCreateModal = signal(false);
  readonly isCreating = signal(false);
  readonly createError = signal('');

  readonly createForm = this.fb.nonNullable.group({
    name: ['', [Validators.required]],
  });

  ngOnInit() {
    this.loadRooms();
  }

  private loadRooms() {
    this.isLoadingRooms.set(true);
    this.roomService.getAll().subscribe({
      next: res => {
        this.rooms.set(res.data ?? []);
        this.isLoadingRooms.set(false);
      },
      error: () => this.isLoadingRooms.set(false),
    });
  }

  openCreateModal() {
    this.createForm.reset();
    this.createError.set('');
    this.showCreateModal.set(true);
  }

  closeCreateModal() {
    this.showCreateModal.set(false);
  }

  createRoom() {
    if (this.createForm.invalid) {
      this.createForm.markAllAsTouched();
      return;
    }
    const { name } = this.createForm.getRawValue();
    this.isCreating.set(true);
    this.createError.set('');
    this.roomService.create(name).subscribe({
      next: () => {
        this.isCreating.set(false);
        this.showCreateModal.set(false);
        this.loadRooms();
      },
      error: () => {
        this.createError.set('Failed to create room. Please try again.');
        this.isCreating.set(false);
      },
    });
  }

  logout() {
    this.authService.logout();
    this.router.navigate(['/login']);
  }
}
