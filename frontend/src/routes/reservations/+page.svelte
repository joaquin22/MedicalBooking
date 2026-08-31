<script lang="ts">
	import { onMount } from 'svelte';
	import { api, type Reservation, type User, type Resource } from '$lib/api';
	import Modal from '$lib/components/Modal.svelte';

	let reservations: Reservation[] = $state([]);
	let users: User[] = $state([]);
	let resources: Resource[] = $state([]);
	let loading = $state(true);
	let error = $state('');
	let saving = $state(false);

	let showForm = $state(false);
	let editing: Reservation | null = $state(null);

	let userId = $state<number>(0);
	let resourceId = $state<number>(0);
	let startTime = $state('');
	let endTime = $state('');
	let status = $state<'pending' | 'confirmed' | 'cancelled' | 'completed'>('pending');
	let notes = $state('');

	function toLocalInput(iso: string): string {
		const d = new Date(iso);
		const pad = (n: number) => String(n).padStart(2, '0');
		return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
	}

	function fromLocalInput(v: string): string {
		return v ? new Date(v).toISOString() : '';
	}

	async function load() {
		loading = true;
		error = '';
		try {
			const [res, us, ru] = await Promise.all([
				api.get<Reservation[]>('/reservations'),
				api.get<User[]>('/users'),
				api.get<Resource[]>('/resources')
			]);
			reservations = res;
			users = us;
			resources = ru;
		} catch (e) {
			error = (e as Error).message;
		} finally {
			loading = false;
		}
	}

	onMount(load);

	function openCreate() {
		editing = null;
		userId = users[0]?.id ?? 0;
		resourceId = resources[0]?.id ?? 0;
		startTime = '';
		endTime = '';
		status = 'pending';
		notes = '';
		showForm = true;
	}

	function openEdit(r: Reservation) {
		editing = r;
		userId = r.user.id;
		resourceId = r.resource.id;
		startTime = toLocalInput(r.start_time);
		endTime = toLocalInput(r.end_time);
		status = r.status;
		notes = r.notes ?? '';
		showForm = true;
	}

	function closeForm() {
		showForm = false;
		editing = null;
	}

	async function submit() {
		saving = true;
		error = '';
		const body = {
			start_time: fromLocalInput(startTime),
			end_time: fromLocalInput(endTime),
			notes
		};
		try {
			if (editing) {
				await api.put<Reservation>(`/reservations/${editing.id}`, { ...body, status });
			} else {
				await api.post<Reservation>('/reservations', {
					user_id: userId,
					resource_id: resourceId,
					...body
				});
			}
			closeForm();
			await load();
		} catch (e) {
			error = (e as Error).message;
		} finally {
			saving = false;
		}
	}

	async function cancel(r: Reservation) {
		if (!confirm(`¿Cancelar la reserva #${r.id}?`)) return;
		error = '';
		try {
			await api.patch<Reservation>(`/reservations/${r.id}/cancel`, {});
			await load();
		} catch (e) {
			error = (e as Error).message;
		}
	}

	async function remove(r: Reservation) {
		if (!confirm(`¿Eliminar la reserva #${r.id}?`)) return;
		error = '';
		try {
			await api.del(`/reservations/${r.id}`);
			await load();
		} catch (e) {
			error = (e as Error).message;
		}
	}

	function fmt(iso: string): string {
		return new Date(iso).toLocaleString();
	}
</script>

<div class="page">
	<div class="topbar">
		<h1>Reservas</h1>
		<button class="primary" onclick={openCreate} disabled={users.length === 0 || resources.length === 0}>
			Nueva reserva
		</button>
	</div>

	{#if error}<div class="error">{error}</div>{/if}
	{#if users.length === 0 || resources.length === 0}
		<p class="empty">Crea primero usuarios y recursos para poder registrar reservas.</p>
	{/if}

	{#if loading}
		<p>Cargando...</p>
	{:else if reservations.length === 0}
		<p class="empty">No hay reservas registradas.</p>
	{:else}
		<table>
			<thead>
				<tr>
					<th>#</th>
					<th>Usuario</th>
					<th>Recurso</th>
					<th>Inicio</th>
					<th>Fin</th>
					<th>Estado</th>
					<th>Notas</th>
					<th>Acciones</th>
				</tr>
			</thead>
			<tbody>
				{#each reservations as r (r.id)}
					<tr>
						<td>{r.id}</td>
						<td>{r.user.full_name}</td>
						<td>{r.resource.name}</td>
						<td>{fmt(r.start_time)}</td>
						<td>{fmt(r.end_time)}</td>
						<td><span class="badge status-{r.status}">{r.status}</span></td>
						<td>{r.notes || '—'}</td>
						<td class="action-gap">
							<button class="ghost" onclick={() => openEdit(r)}>Editar</button>
							{#if r.status !== 'cancelled'}
								<button class="ghost" onclick={() => cancel(r)}>Cancelar</button>
							{/if}
							<button class="ghost danger" onclick={() => remove(r)}>Eliminar</button>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	{/if}

	{#if showForm}
		<Modal title={editing ? 'Editar reserva' : 'Nueva reserva'} onclose={closeForm}>
			<form onsubmit={(e) => { e.preventDefault(); submit(); }}>
				<label>
					Usuario
					<select bind:value={userId} disabled={!!editing}>
						{#each users as u}<option value={u.id}>{u.full_name}</option>{/each}
					</select>
				</label>
				<label>
					Recurso
					<select bind:value={resourceId} disabled={!!editing}>
						{#each resources as r}<option value={r.id}>{r.name}</option>{/each}
					</select>
				</label>
				<label>
					Inicio
					<input type="datetime-local" bind:value={startTime} required />
				</label>
				<label>
					Fin
					<input type="datetime-local" bind:value={endTime} required />
				</label>
				{#if editing}
					<label>
						Estado
						<select bind:value={status}>
							<option value="pending">pending</option>
							<option value="confirmed">confirmed</option>
							<option value="cancelled">cancelled</option>
							<option value="completed">completed</option>
						</select>
					</label>
				{/if}
				<label>
					Notas
					<input type="text" bind:value={notes} />
				</label>
				<div class="form-actions">
					<button type="button" class="ghost" onclick={closeForm}>Cancelar</button>
					<button type="submit" class="primary" disabled={saving}>
						{saving ? 'Guardando...' : 'Guardar'}
					</button>
				</div>
			</form>
		</Modal>
	{/if}
</div>

<style>
	h1 {
		margin: 0;
		font-size: 1.5rem;
	}

	.page {
		max-width: 1200px;
	}

	.topbar {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1.5rem;
	}

	.form-actions {
		display: flex;
		justify-content: flex-end;
		gap: 0.5rem;
		margin-top: 1rem;
	}

	.error {
		background: #fee2e2;
		color: #b91c1c;
		padding: 0.75rem 1rem;
		border-radius: 6px;
		margin-bottom: 1rem;
	}

	.empty {
		color: #6b7280;
	}

	:global(.badge.status-confirmed) {
		background: #d1fae5;
		color: #065f46;
	}

	:global(.badge.status-cancelled) {
		background: #fee2e2;
		color: #991b1b;
	}

	:global(.badge.status-completed) {
		background: #e0e7ff;
		color: #3730a3;
	}
</style>
