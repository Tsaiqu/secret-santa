<script>
    import { onMount } from "svelte";
    import { participantsStore } from "$lib/stores";

    let isDrawing = false;

    // Pobierz listę ludzi z serwera po wejściu na stronę
    onMount(async () => {
        try {
            const res = await fetch('http://localhost:8080/api/admin/participants');
            if (res.ok) {
                const data = await res.json();
                participantsStore.set(data);
            }
        } catch (e) {
            console.error("Nie udało się pobrać uczestników:", e);
        }
    });

    // Big red button
    async function runLottery() {
        if (!confirm("⚠️ Czy na pewno chcesz rozlosować pary i wysłać maile")) return;

        isDrawing = true;
        try {
            const res = await fetch('http://localhost:8080/api/admin/draw', { method: 'POST' });
            const data = await res.json();
            alert(data.message || data.error);
        } catch (error) {
            alert("Wystąpił błąd podczas losowania.");
        } finally {
            isDrawing = false;
        }
    }
</script>

<div class="glass-card" style="max-width: 900px;">
    <div class="header-row">
        <h2>⚙️ Panel Admina</h2>
        <a href="/" class="back-link">← Wróć</a>
    </div>

    <div class="stats">
        <div class="stat-item">
            <span class="label">Elfów w bazie:</span>
            <span class="value">{$participantsStore.length}</span>
        </div>
    </div>

    <div class="table-wrapper">
        <table>
            <thead>
                <tr>
                    <th>Imię</th>
                    <th>Email</th>
                    <th>Preferencje</th>
                </tr>
            </thead>
            <tbody>
                {#each $participantsStore as p}
                    <tr>
                        <td class="fw-bold">{p.name}</td>
                        <td>{p.email}</td>
                        <td class="prefs">{p.preferences}</td>
                    </tr>
                {:else}
                    <tr><td colspan="3" class="empty">Brak zapisanych elfów.</td></tr>
                {/each}
            </tbody>
        </table>
    </div>

    <div class="danger-zone">
        <div class="dz-header">
            <h3>⚠️ Strefa Zagrożenia</h3>
            <p>Akcja jest nieodwracalna. Maile zostaną wysłane 😱</p>
        </div>
        <button class="btn-danger" on:click={runLottery} disabled={isDrawing || $participantsStore.length < 2}>
            {isDrawing ? 'Maszyna losująca pracuje...' : '🚀 ROZLOSUJ I WYŚLIJ'}
        </button>
    </div>
</div>

<style>
    .header-row { display: flex; justify-content: space-between; align-items: center; margin-bottom: 25px; }
    h2 { color: var(--text-main); margin: 0; font-size: 2rem; }
    .back-link { color: var(--elf-green); text-decoration: none; font-weight: 600; }

    .stats { display: flex; gap: 20px; margin-bottom: 20px; }
    .stat-item { background: #f4f6f7; padding: 10px 20px; border-radius: 8px; display: flex; gap: 10px; align-items: center; }
    .stat-item .value { font-weight: bold; color: var(--santa-red); font-size: 1.2rem; }

    .table-wrapper { overflow-x: auto; border-radius: 8px; border: 1px solid #eee; margin-bottom: 30px; }
    table { width: 100%; border-collapse: collapse; min-width: 600px; }
    th { background: var(--bg-gradient-end); color: white; text-align: left; padding: 15px; font-weight: 400; font-size: 0.9rem; letter-spacing: 0.5px; }
    td { padding: 15px; border-bottom: 1px solid #eee; color: var(--text-main); font-size: 0.95rem; }
    tr:last-child td { border-bottom: none; }
    
    .fw-bold { font-weight: 600; }
    .prefs { color: #7f8c8d; font-style: italic; max-width: 300px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
    .empty { text-align: center; padding: 30px; color: #999; }

    /* Danger Zone Design */
    .danger-zone { 
        border: 2px dashed #fab1a0; 
        background: #fff5f5; 
        padding: 25px; 
        border-radius: 12px; 
        display: flex; 
        justify-content: space-between; 
        align-items: center;
        flex-wrap: wrap;
        gap: 20px;
    }
    .dz-header h3 { color: #c0392b; margin: 0 0 5px 0; font-size: 1.4rem; }
    .dz-header p { margin: 0; color: #7f8c8d; font-size: 0.9rem; }
    
    .btn-danger {
        background: #c0392b;
        color: white;
        border: none;
        padding: 12px 24px;
        border-radius: 6px;
        font-weight: bold;
        cursor: pointer;
        transition: background 0.2s;
        box-shadow: 0 4px 10px rgba(192, 57, 43, 0.2);
    }
    .btn-danger:hover:not(:disabled) { background: #a93226; }
    .btn-danger:disabled { background: #e6b0aa; cursor: not-allowed; box-shadow: none; }
</style>