
const rootId = "resumeWizard"
const findResumeId = "findResume"

chrome.runtime.onInstalled.addListener(({ reason }) => {
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
      return
    default:
      console.log("unknown reason being handed:", reason)
      return
  }
});

chrome.contextMenus.onClicked.addListener((info, tab) => {
  switch (info.menuItemId) {
    case rootId:
      console.log("ROOT")
      return
    case findResumeId:
      fetch('http://localhost:8080/health')
        .then(response => {
          if (!response.ok) {
            throw new Error('Network response was not ok');
          }
          console.log(response)
          // return response.json();
        })
        // .then(data => console.log(data))
        .catch(error => console.error('Error:', error));
      return
    default:
      console.log("unknown menu id", info.menuItemId)
      return
  }
});
