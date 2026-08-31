<script lang="ts">
	import { onMount } from 'svelte';
	import { api, type Resource } from '$lib/api';
	import Modal from '$lib/components/Modal.svelte';

	let resources: Resource[] = $state([]);
	let loading = $state(true);
	let error = $state('');
	let saving = $state(false);

	let showForm = $state(false);
	let editing: Resource | null = $state(null);

	let name = $state('');
	let type = $state<'doctor' | 'room' | 'other'>('doctor');
	let description = $state('');
	let location = $state('');
	let capacity = $state<number>(1);
	let active = $state(true);

	async function load() {
		loading = true;
		error = '';
		try {
			resources = await api.get<Resource[]>('/resources');
		} catch (e) {
			error = (e as Error).message;
		} finally {
			loading = false;
		}
	}

	onMount(load);

	function openCreate() {
		editing = null;
		name = '';
		type = 'doctor';
		description = '';
		location = '';
		capacity = 1;
		active = true;
		showForm = true;
	}

	function openEdit(r: Resource) {
		editing = r;
		name = r.name;
		type = r.type;
		description = r.description ?? '';
		location = r.location ?? '';
		capacity = r.capacity;
		active = r.active;
		showForm = true;
	}

	function closeForm() {
		showForm = false;
		editing = null;
	}

	async function submit() {
		saving = true;
		error = '';
		const base = { name, description, location, capacity };
		try {
			if (editing) {
				await api.put<Resource>(`/resources/${editing.id}`, { ...base, active });
			} else {
				await api.post<Resource>('/resources', { ...base, type });
			}
			closeForm();
			await load();
		} catch (e) {
			error = (e as Error).message;
		} finally {
			saving = false;
		}
	}

	async function remove(r: Resource) {
		if (!confirm(`¿Eliminar el recurso "${r.name}"?`)) return;
		error = '';
		try {
			await api.del(`/resources/${r.id}`);
			await load();
		} catch (e) {
			error = (e as Error).message;
		}
	}
</script>

<div class="page">
	<div class="topbar">
		<h1>Recursos</h1>
		<button class="primary" onclick={openCreate}>Nuevo recurso</button>
	</div>

	{#if error}<div class="error">{error}</div>{/if}

	{#if loading}
		<p>Cargando...</p>
	{:else if resources.length === 0}
		<p class="empty">No hay recursos registrados.</p>
	{:else}
		<table>
			<thead>
				<tr>
					<th>Nombre</th>
					<th>Tipo</th>
					<th>Descripción</th>
					<th>Ubicación</th>
					<th>Capacidad</th>
					<th>Estado</th>
					<th>Acciones</th>
				</tr>
			</thead>
			<tbody>
				{#each resources as r (r.id)}
					<tr>
						<td>{r.name}</td>
						<td><span class="badge">{r.type}</span></td>
						<td>{r.description || '—'}</td>
						<td>{r.location || '—'}</td>
						<td>{r.capacity}</td>
						<td>{r.active ? 'Activo' : 'Inactivo'}</td>
						<td class="action-gap">
							<button class="ghost" onclick={() => openEdit(r)}>Editar</button>
							<button class="ghost danger" onclick={() => remove(r)}>Eliminar</button>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	{/if}

	{#if showForm}
		<Modal title={editing ? 'Editar recurso' : 'Nuevo recurso'} onclose={closeForm}>
			<form onsubmit={(e) => { e.preventDefault(); submit(); }}>
				<label>
					Nombre
					<input type="text" bind:value={name} required minlength={2} />
				</label>
				{#if !editing}
					<label>
						Tipo
						<select bind:value={type}>
							<option value="doctor">doctor</option>
							<option value="room">room</option>
							<option value="other">other</option>
						</select>
					</label>
				{/if}
				<label>
					Descripción
					<input type="text" bind:value={description} />
				</label>
				<label>
					Ubicación
					<input type="text" bind:value={location} />
				</label>
				<label>
					Capacidad
					<input type="number" min="1" bind:value={capacity} />
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
		max-width: 1100px;
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
