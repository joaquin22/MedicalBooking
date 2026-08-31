<script lang="ts">
	interface Props {
		title: string;
		open?: boolean;
	}

	let { title, open = true, children, onclose }: Props & {
		children?: import('svelte').Snippet;
		onclose?: () => void;
	} = $props();

	const show = $derived(open);
</script>

{#if show}
	<div class="overlay" role="presentation" onclick={onclose} onkeydown={(e) => { if (e.key === 'Escape') onclose?.(); }}>
		<div class="modal" role="dialog" aria-modal="true" tabindex="-1" onclick={(e) => e.stopPropagation()}>
			<div class="header">
				<h2>{title}</h2>
				<button class="close" onclick={onclose} aria-label="Cerrar">×</button>
			</div>
			<div class="body">
				{@render children?.()}
			</div>
		</div>
	</div>
{/if}

<style>
	.overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 50;
		padding: 1rem;
	}

	.modal {
		background: #fff;
		border-radius: 8px;
		width: 100%;
		max-width: 480px;
		max-height: 90vh;
		overflow: auto;
		box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
	}

	.header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 1rem 1.25rem;
		border-bottom: 1px solid #e5e7eb;
	}

	.header h2 {
		margin: 0;
		font-size: 1.1rem;
	}

	.close {
		border: none;
		background: none;
		font-size: 1.5rem;
		line-height: 1;
		cursor: pointer;
		color: #6b7280;
	}

	.body {
		padding: 1.25rem;
	}
</style>
