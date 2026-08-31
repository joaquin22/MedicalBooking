<script lang="ts">
	import '$lib/app.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();

	const links = [
		{ href: '/users', label: 'Usuarios' },
		{ href: '/resources', label: 'Recursos' },
		{ href: '/reservations', label: 'Reservas' }
	];
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="app">
	<aside>
		<div class="brand">MedicalBooking</div>
		<nav>
			{#each links as link}
				<a href={link.href}>{link.label}</a>
			{/each}
		</nav>
	</aside>

	<main>
		{@render children()}
	</main>
</div>

<style>
	.app {
		display: flex;
		min-height: 100vh;
	}

	aside {
		width: 220px;
		background: #1f2937;
		color: #e5e7eb;
		display: flex;
		flex-direction: column;
		padding: 1rem 0;
	}

	.brand {
		padding: 0 1.5rem 1.5rem;
		font-size: 1.2rem;
		font-weight: 700;
		color: #fff;
	}

	nav {
		display: flex;
		flex-direction: column;
	}

	nav :global(a) {
		padding: 0.75rem 1.5rem;
		color: #9ca3af;
		text-decoration: none;
	}

	nav :global(a:hover),
	nav :global(a.active) {
		background: #374151;
		color: #fff;
	}

	main {
		flex: 1;
		padding: 2rem;
		background: #f3f4f6;
	}
</style>
