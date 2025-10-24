// NOTE: listeners wrapped in funcs to avoid namespace conflicts

// filtering/search
function srch() {
  document.addEventListener("DOMContentLoaded", function () {
    const searchInput = document.getElementById("searchInput");
  
    searchInput.addEventListener("input", function () {
      const rows = document.querySelectorAll("#datatable tbody tr");
      const filter = this.value.toLowerCase();
  
      const filter_elem_ors = filter.split(/\|\|/g)
  
      rows.forEach(row => {
        const text = row.textContent.toLowerCase();
        let seen_or = false; // NOTE: ors can't make false if we start from true
  
        for (const filter_elem_or of filter_elem_ors) {
          const filter_elems_ands = filter_elem_or.split(/&&/g)
          let seen_and = true; // NOTE: ands can't make true, if we start from false
          for (const filter_elem_and of filter_elems_ands) {
            seen_and = seen_and && text.includes(filter_elem_and.trim());
          }
          seen_or = seen_or || seen_and;
        }
        row.style.display = seen_or ? "" : "none";
      });
    });
  });
}


// sorting
function srt() {
  const tables = document.querySelectorAll("#datatable");
  tables.forEach(table => {
    table.querySelectorAll("th").forEach((th, idx) => {
      th.addEventListener("click", function () {
        const tbody = table.querySelector("tbody");
        const columnIndex = idx;;
        const rows = Array.from(tbody.querySelectorAll("tr"));
        const ascending = th.classList.toggle("asc");
  
        rows.sort((a, b) => {
          const aText = a.children[columnIndex].innerText.trim();
          const bText = b.children[columnIndex].innerText.trim();
  
          // Try numeric comparison, fallback to string
          const aNum = parseFloat(aText);
          const bNum = parseFloat(bText);
          const aVal = isNaN(aNum) ? aText : aNum;
          const bVal = isNaN(bNum) ? bText : bNum;
  
          if (aVal > bVal) return ascending ? 1 : -1;
          if (aVal < bVal) return ascending ? -1 : 1;
          return 0;
        });
  
        tbody.innerHTML = "";
        rows.forEach(r => tbody.appendChild(r));
      });
    });
  });
}

srch()
srt()