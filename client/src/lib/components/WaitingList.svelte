<script lang="ts">
	import type { RoomMember } from '$lib/ws/messages';

	interface Props {
		playerID: string;
		roomMembers: RoomMember[];
	}
	let { playerID, roomMembers }: Props = $props();
</script>

<div class="list-wrapper">
	<h3>Players in Room ({roomMembers.length})</h3>
	<ul>
		{#each roomMembers as player (player.id)}
			<li class="player-item" class:current-player={player.id === playerID}>
				<div class="player-info">
					<span class="player-nickname">{player.nickname}</span>
					{#if player.id === playerID}
						<span class="you-badge">You</span>
					{/if}
				</div>
				<div class="player-status">
					<div class="status-indicator waiting"></div>
					<span class="status-text">Waiting...</span>
				</div>
			</li>
		{/each}
	</ul>
	{#if roomMembers.length === 0}
		<div class="empty-state">
			<p>No players in the room yet</p>
		</div>
	{/if}
</div>

<style>
	.list-wrapper {
    max-height: 20rem;

		background: #f8f9fa;
		border-radius: 8px;
		padding: 1.5rem;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

    overflow-y: scroll;
	}

	h3 {
		margin: 0 0 1rem 0;
		color: #333;
		font-size: 1.2rem;
		font-weight: 600;
	}

	ul {
		list-style: none;
		padding: 0;
		margin: 0;
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.player-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem;
		background: white;
		border-radius: 6px;
		border: 2px solid #e9ecef;
		transition: all 0.2s ease;
	}

	.player-item:hover {
		border-color: #dee2e6;
		transform: translateY(-1px);
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
	}

	.current-player {
		border-color: #007bff;
		background: #f0f8ff;
	}

	.player-info {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.player-nickname {
		font-weight: 500;
		color: #333;
		font-size: 1rem;
	}

	.you-badge {
		background: #007bff;
		color: white;
		padding: 0.25rem 0.5rem;
		border-radius: 12px;
		font-size: 0.75rem;
		font-weight: 600;
		text-transform: uppercase;
	}

	.player-status {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.status-indicator {
		width: 8px;
		height: 8px;
		border-radius: 50%;
		animation: pulse 2s infinite;
	}

	.status-indicator.waiting {
		background: #ffc107;
	}

	.status-text {
		color: #6c757d;
		font-size: 0.875rem;
		font-style: italic;
	}

	.empty-state {
		text-align: center;
		padding: 2rem;
		color: #6c757d;
	}

	.empty-state p {
		margin: 0;
		font-style: italic;
	}

	@keyframes pulse {
		0%,
		100% {
			opacity: 1;
		}
		50% {
			opacity: 0.5;
		}
	}

	@media (max-width: 480px) {
		.list-wrapper {
			padding: 1rem;
		}

		.player-item {
			flex-direction: column;
			align-items: flex-start;
			gap: 0.5rem;
		}

		.player-status {
			align-self: flex-end;
		}
	}
</style>
