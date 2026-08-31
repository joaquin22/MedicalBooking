<script lang="ts">
	import { onMount } from 'svelte';
	import { api, type User } from '$lib/api';
	import Modal from '$lib/components/Modal.svelte';

	let users: User[] = $state([]);
	let loading = $state(true);
	let error = $state('');
	let saving = $state(false);

	let showForm = $state(false);
	let editing: User | null = $state(null);

	let fullName = $state('');
	let email = $state('');
	let phone = $state('');
	let password = $state('');
	let role = $state<'admin' | 'staff' | 'customer'>('customer');
	let active = $state(true);

	async function load() {
		loading = true;
		error = '';
		try {
			users = await api.get<User[]>('/users');
		} catch (e) {
			error = (e as Error).message;
		} finally {
			loading = false;
		}
	}

	onMount(load);

	function openCreate() {
		editing = null;
		fullName = '';
		email = '';
		phone = '';
		password = '';
		role = 'customer';
		active = true;
		showForm = true;
	}

	function openEdit(u: User) {
		editing = u;
		fullName = u.full_name;
		email = u.email;
		phone = u.phone ?? '';
		password = '';
		role = u.role;
		active = u.active;
		showForm = true;
	}

	function closeForm() {
		showForm = false;
		editing = null;
	}

	async function submit() {
		saving = true;
		error = '';
		const base = { full_name: fullName, phone };
		try {
			if (editing) {
				await api.put<User>(`/users/${editing.id}`, { ...base, active });
			} else {
				await api.post<User>('/users', { ...base, email, password, role });
			}
			closeForm();
			await load();
		} catch (e) {
			error = (e as Error).message;
		} finally {
			saving = false;
		}
	}

	async function remove(u: User) {
		if (!confirm(`¿Eliminar al usuario "${u.full_name}"?`)) return;
		error = '';
		try {
			await api.del(`/users/${u.id}`);
			await load();
		} catch (e) {
			error = (e as Error).message;
		}
	}
</script>

<div class="page">
	<div class="topbar">
		<h1>Usuarios</h1>
		<button class="primary" onclick={openCreate}>Nuevo usuario</button>
	</div>

	{#if error}<div class="error">{error}</div>{/if}

	{#if loading}
		<p>Cargando...</p>
	{:else if users.length === 0}
		<p class="empty">No hay usuarios registrados.</p>
	{:else}
		<table>
			<thead>
				<tr>
					<th>Nombre</th>
					<th>Email</th>
					<th>Teléfono</th>
					<th>Rol</th>
					<th>Estado</th>
					<th>Acciones</th>
				</tr>
			</thead>
			<tbody>
				{#each users as u (u.id)}
					<tr>
						<td>{u.full_name}</td>
						<td>{u.email}</td>
						<td>{u.phone || '—'}</td>
						<td><span class="badge">{u.role}</span></td>
						<td>{u.active ? 'Activo' : 'Inactivo'}</td>
						<td class="action-gap">
							<button class="ghost" onclick={() => openEdit(u)}>Editar</button>
							<button class="ghost danger" onclick={() => remove(u)}>Eliminar</button>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	{/if}

	{#if showForm}
		<Modal title={editing ? 'Editar usuario' : 'Nuevo usuario'} onclose={closeForm}>
			<form onsubmit={(e) => { e.preventDefault(); submit(); }}>
				<label>
					Nombre completo
					<input type="text" bind:value={fullName} required minlength={3} />
				</label>
				<label>
					Email
					<input type="email" bind:value={email} required disabled={!!editing} />
				</label>
				{#if !editing}
					<label>
						Contraseña
						<input type="password" bind:value={password} required minlength={6} />
					</label>
				{/if}
				<label>
					Teléfono
					<input type="text" bind:value={phone} />
				</label>
				<label>
					Rol
					<select bind:value={role} disabled={!!editing}>
						<option value="customer">customer</option>
						<option value="staff">staff</option>
						<option value="admin">admin</option>
					</select>
				</label>
				{#if editing}
					<label class="checkbox">
						<input type="checkbox" bind:checked={active} />
						Activo
					</label>
				{/if}
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
		max-width: 900px;
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
</style>
