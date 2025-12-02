<script>
    let name = '';
    let email = '';
    let preferences = '';
    let isSubmitting = false;

    async function handleSubmit() {
        isSubmitting = true;
        try {
            const response = await fetch('http://localhost:8080/api/signup', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name, email, preferences })
            });

            if (response.ok) {
                alert("Ho ho ho! Zapisano pomyślnie! 🎄");
                name = ''; email = ''; preferences = '';
            } else {
                const err = await response.json();
                alert("Ups: " + err.error);
            }
        } catch (e) {
            alert("Błąd połączenia z serwerem.");
        } finally {
            isSubmitting = false;
        }
    }
</script>

<div class="glass-card">
    <header>
        <h1>Świąteczne Losowanie</h1>
        <p>Wypełnij list do Mikołaja, a magia zrobi resztę ✨</p>
    </header>

    <form on:submit|preventDefault={handleSubmit}>
        <div class="input-group">
            <label for="name">Twoje Imię</label>
            <input id="name" type="text" bind:value={name} required placeholder="np. Elf Bartłomiej">
        </div>

        <div class="input-group">
            <label for="email">Twój Email</label>
            <input id="email" type="email" bind:value={email} required placeholder="bartlomiej@biegun-polnocny.pl">
        </div>

        <div class="input-group">
            <label for="prefs">List do Mikołaja (Preferencje)</label>
            <textarea id="prefs" bind:value={preferences} rows="4" required placeholder="Napisz co lubisz, podaj rozmiar skarpetek..."></textarea>
        </div>

        <button type="submit" disabled={isSubmitting}>
            {#if isSubmitting}
                Wysyłanie saniami... 🛷
            {:else}
                🎄 Wrzuć do worka
            {/if}
        </button>
    </form>
    
    <footer>
        <a href="/admin">Panel Admina</a>
    </footer>
</div>

<style>
    header { text-align: center; margin-bottom: 30px; }
    h1 { color: var(--santa-red); font-size: 2.5rem; margin: 0; letter-spacing: 1px; }
    p { color: var(--text-light); margin-top: 5px; font-size: 0.95rem; }

    .input-group { margin-bottom: 20px; }
    label { display: block; font-weight: 600; margin-bottom: 8px; color: var(--text-main); font-size: 0.9rem; }
    
    input, textarea {
        width: 100%;
        padding: 12px 15px;
        border: 2px solid #e0e0e0;
        border-radius: 8px;
        font-family: inherit;
        font-size: 1rem;
        transition: all 0.3s ease;
        background: #f9f9f9;
        box-sizing: border-box; /* Ważne, żeby padding nie rozpychał */
    }

    input:focus, textarea:focus {
        outline: none;
        border-color: var(--elf-green);
        background: white;
        box-shadow: 0 0 0 4px rgba(39, 174, 96, 0.1);
    }

    button {
        width: 100%;
        padding: 16px;
        background: linear-gradient(135deg, var(--santa-red), #922b21);
        color: white;
        border: none;
        border-radius: 8px;
        font-size: 1.1rem;
        font-weight: 600;
        cursor: pointer;
        transition: transform 0.2s, box-shadow 0.2s;
        box-shadow: 0 4px 15px rgba(192, 57, 43, 0.3);
    }

    button:hover:not(:disabled) {
        transform: translateY(-2px);
        box-shadow: 0 6px 20px rgba(192, 57, 43, 0.4);
    }
    
    button:disabled { opacity: 0.7; cursor: wait; }

    footer { text-align: center; margin-top: 25px; border-top: 1px solid #eee; padding-top: 15px; }
    footer a { color: #bbb; text-decoration: none; font-size: 0.8rem; transition: color 0.3s; }
    footer a:hover { color: var(--text-main); }
</style>