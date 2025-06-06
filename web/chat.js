var userid;
var chatList;

async function main(params) {
  userid = await verifyToken();
  chatList = await getChatList();




// Simulated chat data
const chatData = {
  Alice: [
    { type: "received", text: "Hi! How's it going?" },
    { type: "sent", text: "All good! You?" },
  ],
  Bob: [
    { type: "received", text: "Let's catch up soon." },
    { type: "sent", text: "Sure, just say when!" },
  ],
};

  const chatHeader = document.querySelector(".chat-header h3");
  const chatMessages = document.querySelector(".chat-messages");

  // Add click event listeners after DOM is populated
  document.querySelector(".chat-list")?.addEventListener("click", (event) => {
    const chatItem = event.target.closest(".chat-item");
    if (!chatItem) return;

    document.querySelectorAll(".chat-item").forEach((item) => {
      item.classList.remove("active");
    });

    chatItem.classList.add("active");
    const name = chatItem.querySelector(".name").innerText;
    chatHeader.textContent = name;

    const messages = chatData[name] || [];
    chatMessages.innerHTML = "";
    messages.forEach((msg) => {
      const div = document.createElement("div");
      div.classList.add("message", msg.type);
      div.textContent = msg.text;
      chatMessages.appendChild(div);
    });
  });






}

function verifyToken() {
  return new Promise((resolve, reject) => {
    if (window.localStorage.getItem("token") != null) {
      fetch("/api/verifytoken", {
        method: "POST",
        headers: {
          Authorization: window.localStorage.getItem("token"),
        },
      })
        .then((res) => {
          if (res.ok) {
            return res.json();
          } else {
            throw new Error("Token verification failed");
          }
        })
        .then((data) => {
          resolve(data.userID);
        })
        .catch((error) => {
          console.error("Error verifying token:", error);
          alert("Please log in first!");
          window.location.href = "/login";
          reject(error);
        });
    } else {
      alert("Please log in first!");
      window.location.href = "/login";
      reject(new Error("No token found"));
    }
  });
}

function getChatList() {
  return new Promise((resolve, reject) => {
    fetch("/api/getchatsbyuser", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: window.localStorage.getItem("token"),
      },
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("Failed to fetch chats");
        }
        return response.json();
      })
      .then((chatList) => {
        chatList = chatList;
        const chatListContainer = document.querySelector(".chat-list");
        chatListContainer.innerHTML = "";

        chatList.forEach((chat) => {
          const participant = chat.Participants.find((p) => p.ID !== userid);
          const chatItem = document.createElement("div");
          chatItem.classList.add("chat-item");
          chatItem.innerHTML = `
          <div class="avatar">
            <img src="${
              participant.ProfilePicture || "default-avatar.png"
            }" alt="${participant.Name || "User"}">
            <div class="status-indicator ${
              participant.IsOnline ? "online" : ""
            }"></div>
          </div>
          <div class="chat-info">
            <h4 class="name">${participant.Name || "Unknown User"}</h4>
          </div>
        `;
          chatListContainer.appendChild(chatItem);
        });
        resolve(chatList);
      })
      .catch((error) => {
        console.error("Error fetching chats:", error);
      });
  });
}

main();
