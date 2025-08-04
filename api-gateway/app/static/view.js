document.addEventListener('DOMContentLoaded', () => {
  const id = window.location.pathname.split("/").pop();

  const tabDiv = document.getElementById('tabResult');
  const tabNameEl = document.getElementById("tabName");
  const loadingEl = document.getElementById("loading");

  fetch(`/tab/${id}`)
    .then(res => {
      if (!res.ok) throw new Error("Не удалось загрузить метаданные");
      return res.json();
    })
    .then(tab => {
      loadingEl.classList.add("hidden");
      tabNameEl.textContent = tab.name;
      displayTab(tab.body);
    })
    .catch(err => {
      tabNameEl.textContent = "Ошибка загрузки табулатуры";
      tabDiv.innerHTML = `
        <div class="error-text card" style="text-align:center; padding:20px;">
          ${err.message}
        </div>`;
      console.error(err);
    });

  function displayTab(tabText) {
    if (!tabText) {
      tabDiv.innerHTML = `
        <div class="error-text card" style="text-align:center; padding:20px;">
          Ошибка загрузки табулатуры
        </div>`;
      return;
    }

    tabDiv.innerHTML = '';
    const lines = tabText.split('\n').filter(line => line.trim());

    lines.forEach(line => {
      const lineDiv = document.createElement('div');
      lineDiv.className = 'tab-line font-mono text-white';

      const highlightedLine = line.replace(/(\d+)/g, '<span class="note">$1</span>');
      lineDiv.innerHTML = highlightedLine;

      tabDiv.appendChild(lineDiv);
    });
  }
});