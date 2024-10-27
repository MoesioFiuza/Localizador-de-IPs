// Função para carregar o JSON com os dados de localização
async function carregarDados() {
    try {
        // Faz a requisição ao endpoint /dados
        const response = await fetch('/dados');
        const data = await response.json();

        // Extrai as coordenadas e a cidade dos dados
        const coordenadas = data.coordenadas.split(',').map(Number); // Converte para [latitude, longitude]
        const cidade = data.cidade;

        // Inicializa o mapa no centro das coordenadas
        const mapa = L.map('map').setView(coordenadas, 13);

        // Adiciona o tile layer do mapa (OpenStreetMap)
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            maxZoom: 18,
            attribution: 'Map data © <a href="https://openstreetmap.org">OpenStreetMap</a> contributors'
        }).addTo(mapa);

        // Adiciona um marcador na posição com popup mostrando a cidade e as coordenadas
        L.marker(coordenadas).addTo(mapa)
            .bindPopup(`<b>${cidade}</b><br>Coordenadas: ${data.coordenadas}`)
            .openPopup();

    } catch (error) {
        console.error("Erro ao carregar os dados de localização:", error);
    }
}

// Carrega os dados e inicializa o mapa
carregarDados();
