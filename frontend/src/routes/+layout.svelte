<slot />

<style>
    /* Import fontów */
    @import url('https://fonts.googleapis.com/css2?family=Mountains+of+Christmas:wght@700&family=Poppins:wght@300;400;600&display=swap');

    :global(:root) {
        --bg-gradient-start: #0f2027;
        --bg-gradient-mid: #203a43;
        --bg-gradient-end: #2c5364;
        
        --card-bg: rgba(255, 255, 255, 0.95); /* Prawie białe szkło */
        --text-main: #2c3e50;
        --text-light: #7f8c8d;
        
        --santa-red: #c0392b;
        --elf-green: #27ae60;
        --gold: #f1c40f;
        
        --radius: 16px;
        --shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
    }

    :global(body) {
        font-family: 'Poppins', sans-serif;
        color: var(--text-main);
        margin: 0;
        min-height: 100vh;
        display: flex;
        justify-content: center;
        align-items: center;
        padding: 20px;
        box-sizing: border-box;

        /* --- MAGICZNE TŁO --- */
        background-color: #0f2027; /* Kolor bazowy (ciemny) */
        
        /* Tutaj definiujemy dwie warstwy tła:
           1. Gradient (półprzezroczysty), żeby przyciemnić wzorki.
           2. Wzorek SVG zdefiniowany bezpośrednio w kodzie (emoji).
        */
        background-image: 
            linear-gradient(to bottom, rgba(15, 32, 39, 0.9), rgba(44, 83, 100, 0.85)),
            url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='100' height='100' viewBox='0 0 100 100'%3E%3Ctext x='10' y='35' font-size='25'%3E🎅%3C/text%3E%3Ctext x='60' y='35' font-size='25'%3E🎄%3C/text%3E%3Ctext x='10' y='85' font-size='25'%3E🎁%3C/text%3E%3Ctext x='60' y='85' font-size='25'%3E🍪%3C/text%3E%3C/svg%3E");
            
        background-repeat: repeat; /* Powtarzaj wzorek */
        background-size: cover 100px 100px; /* Rozmiar kafelka wzoru */
        background-attachment: fixed; /* Tło nie przesuwa się przy scrollowaniu */
    }

    :global(.glass-card) {
        background: var(--card-bg);
        padding: 40px;
        border-radius: var(--radius);
        box-shadow: var(--shadow);
        width: 100%;
        max-width: 500px;
        position: relative;
        overflow: hidden;
        animation: slideIn 0.6s cubic-bezier(0.25, 1, 0.5, 1);
        
        /* Dodatkowy efekt rozmycia tła pod kartą (frosted glass) */
        backdrop-filter: blur(5px); 
    }

    /* Złoty pasek na górze */
    :global(.glass-card::before) {
        content: '';
        position: absolute;
        top: 0; left: 0; right: 0;
        height: 6px;
        background: linear-gradient(90deg, var(--santa-red), var(--gold), var(--elf-green));
    }

    @keyframes slideIn {
        from { opacity: 0; transform: translateY(30px); }
        to { opacity: 1; transform: translateY(0); }
    }

    :global(h1, h2, h3) {
        font-family: 'Mountains of Christmas', cursive;
        font-weight: 700;
    }
</style>