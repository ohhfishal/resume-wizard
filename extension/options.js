
const DEFAULT_SERVER_URL = 'http://localhost:8080';
const DEFAULT_USER_ID = '';

const form = document.getElementById("form");

document.addEventListener('DOMContentLoaded', () => {
  chrome.storage.local.get(['server_url', 'user_id'], (result) => {
    document.getElementById('server_url').value = result.server_url || DEFAULT_SERVER_URL;
    document.getElementById('user_id').value = result.user_id || DEFAULT_USER_ID;
  });
});

form.addEventListener("submit", (e) => {
  e.preventDefault();
  const serverUrl = document.getElementById('server_url').value.trim();
  const userId = document.getElementById('user_id').value.trim();

  // TODO: Validate that the user_id is valid!

  chrome.storage.local.set({ 
    "server_url": serverUrl, 
    "user_id": userId 
  }).then(() => {
    showStatus("Settings saved!", "success");
  });
});

form.addEventListener("reset", (e) => {
  e.preventDefault();

  document.getElementById('server_url').value = DEFAULT_SERVER_URL;

  chrome.storage.local.set({ 
    "server_url": DEFAULT_SERVER_URL, 
  }).then(() => {
    showStatus('Reset to default values', 'success');
  });
});

function showStatus(message, type) {
  const statusEl = document.getElementById('status');
  statusEl.textContent = message;
  statusEl.className = type;

  setTimeout(() => {
    statusEl.className = '';
  }, 3000);
}
