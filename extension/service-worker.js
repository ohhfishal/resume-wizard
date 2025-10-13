const DEFAULT_SERVER_URL = "https://localhost:8080"
const rootId = "resumeWizard"
const findResumeId = "findResume"

chrome.runtime.onInstalled.addListener(async ({ reason }) => {
  switch (reason) {
    case "update":
    case "install":
      // chrome.contextMenus.create({
      //   id: rootId,
      //   title: "Resume Wizard",
      //   contexts: ["page", "selection"]
      // });

      chrome.contextMenus.create({
        id: findResumeId,
        title: "Match resume to selection",
        contexts: ["selection"],
        // parentId: rootId
      });

      let options = await chrome.storage.local.get(["server_url", "user_id"])

      if (options.server_url === undefined) {
        options.server_url = DEFAULT_SERVER_URL
        chrome.storage.local.set({ server_url: options.server_url}).then(() => {
          console.log("Stored value: server_url = " + options.server_url);
        });
      }

      if (options.user_id === undefined) {
        console.log("Getting new user_id");
        fetch(`${options.server_url}}/api/register`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
        })
        .then(response => {
          if (!response.ok) {
            throw new Error('registering: ' + response.statusText);
          }
          return response.json();
        })
        .then(user => {
          chrome.storage.local.set({ user_id: user.id }).then(() => {
            console.log("Stored value:", user);
          });
        })
        .catch(error => console.error('Error:', error));
      }

      return
    default:
      console.log("unknown reason being handed:", reason)
      return
  }
});

chrome.contextMenus.onClicked.addListener(async (info, tab) => {
  switch (info.menuItemId) {
    case rootId:
      console.log("ROOT")
      return
    case findResumeId:
      let options = await chrome.storage.local.get(["server_url", "user_id"]);
      fetch(`${options.server_url}/api/match`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${options.user_id}`
        },
      })
        .then(response => {
          console.log(response)
          if (!response.ok) {
            throw new Error('requesting match: ' + response.statusText);
          }
          return response.json();
        })
        .then(data => console.log(data))
        .catch(error => console.error('Error:', error));
      return
    default:
      console.log("unknown menu id", info.menuItemId)
      return
  }
});
