<script>
	import { onMount } from "svelte";
	import { page } from "$app/stores";

	let loading = true;
	let error = "";
	let data = null;

	onMount(async () => {
		// Pobieramy token z adresu URL (?token=...)
		const token = $page.url.searchParams.get("token");

		if (!token) {
			error = "Brak magicznego linku. Sprawdź swój email!";
			loading = false;
			return;
		}

		try {
			const res = await fetch(
				`http://localhost:8080/api/my-status?token=${token}`,
			);
			if (res.ok) {
				data = await res.json();
			} else {
				error = "Nieprawidłowy link lub wygasł.";
			}
		} catch (e) {
			error = "Błąd połączenia z serwerem.";
		} finally {
			loading = false;
		}
	});
</script>

<div class="glass-card">
	{#if loading}
		<div class="center">
			<h2>Szukam Twojego prezentu... 🎁</h2>
		</div>
	{:else if error}
		<div class="center">
			<h2 style="color: var(--santa-red)">Ojej! ☹️</h2>
			<p>{error}</p>
			<a href="/">Wróć na stronę główną</a>
		</div>
	{:else if data}
		<header>
			<h1>Cześć, {data.me.name}! 👋</h1>
			<p>To jest Twój tajny panel dowodzenia.</p>
		</header>

		<div class="section">
			<h3>Twoje preferencje:</h3>
			<div class="box info">
				{data.me.preferences}
			</div>
			<small>Chcesz to zmienić? Napisz do Admina.</small>
		</div>

		<hr />

		<div class="section">
			<h3>Kogo wylosowałeś?</h3>

			{#if data.is_draw_done}
				<div class="target-reveal">
					<p class="intro">
						Twoim celem w tym roku jest:
					</p>
					<h2 class="target-name">
						✨ {data.target_name} ✨
					</h2>

					<p class="intro">
						List do Mikołaja tej osoby:
					</p>
					<div class="box gift">
						{data.target_prefs}
					</div>
				</div>
			{:else}
				<div class="box waiting">
					<p>
						⏳ Maszyna losująca jeszcze nie
						ruszyła.
					</p>
					<p>
						Sprawdź ten link ponownie, gdy
						otrzymasz maila o losowaniu!
					</p>
				</div>
			{/if}
		</div>
	{/if}
</div>

<style>
	.center {
		text-align: center;
	}
	header {
		text-align: center;
		margin-bottom: 30px;
	}
	h1 {
		color: var(--text-main);
		margin: 0;
	}

	.section {
		margin-bottom: 25px;
	}
	h3 {
		color: var(--elf-green);
		margin-bottom: 10px;
		font-size: 1.2rem;
	}

	.box {
		padding: 15px;
		border-radius: 8px;
		font-size: 0.95rem;
		line-height: 1.5;
	}
	.info {
		background: #f0f4f8;
		border-left: 4px solid #888;
		color: #555;
	}
	.waiting {
		background: #fff3cd;
		border: 1px dashed #cca000;
		text-align: center;
		color: #856404;
	}

	.target-reveal {
		text-align: center;
		animation: fadeIn 1s;
	}
	.target-name {
		color: var(--santa-red);
		font-size: 2.2rem;
		margin: 15px 0;
		text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	}
	.gift {
		background: #d4edda;
		border: 2px solid #c3e6cb;
		color: #155724;
		text-align: left;
		margin-top: 10px;
	}

	hr {
		border: 0;
		border-top: 1px solid #eee;
		margin: 30px 0;
	}

	a {
		color: var(--text-main);
		font-weight: bold;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
			transform: scale(0.9);
		}
		to {
			opacity: 1;
			transform: scale(1);
		}
	}
</style>

