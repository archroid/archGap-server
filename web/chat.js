function verifyToken() {
  if (window.localStorage.getItem("token") != null) {
    let request = new XMLHttpRequest();
    request.open("POST", "/api/verifytoken");
    // console.log(window.localStorage.getItem("token"))
    request.setRequestHeader(
      "Authorization",
      window.localStorage.getItem("token")
    );
    request.send();
    if (!request.ok) {
      return;
    } else {
      alert("Please log in first!");
      window.location.href = "/login";
    }
  } else {
    alert("Please log in first!");
    window.location.href = "/login";
  }
}

verifyToken();

try {
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
      const chatListContainer = document.querySelector(".chat-list");
      chatListContainer.innerHTML = "";

      chatList.forEach((chat) => {
        const participant = chat.Participants.find(
          (p) => p.ID !== 1 // Assuming the current user has ID 1
        );
        const chatItem = document.createElement("div");
        chatItem.classList.add("chat-item");
        chatItem.innerHTML = `
            <div class="avatar">
              <img  src="${
                participant.ProfilePicture || "default-avatar.png"
              }" width="40" height="40" alt="${participant.Name || "User"}">
            </div>
            <div class="chat-info">
              <h4 class="name">${participant.Name || "Unknown User"}</h4>
              <p class="last-message">Last message preview...</p>
            </div>
          `;
        chatListContainer.appendChild(chatItem);
      });
    })
    .catch((error) => {
      console.error("Error fetching chats:", error);
    });
} catch (error) {
  console.error("Unexpected error:", error);
}

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

document.addEventListener("DOMContentLoaded", () => {
  const chatHeader = document.querySelector(".chat-header h3");
  const chatMessages = document.querySelector(".chat-messages");

  // Add click event listeners to chat items
  document.querySelector(".chat-list").addEventListener("click", (event) => {
    const chatItem = event.target.closest(".chat-item");
    if (!chatItem) return;

    // Remove active class from all chat items
    document.querySelectorAll(".chat-item").forEach((item) => {
      item.classList.remove("active");
    });

    // Add active class to the clicked chat item
    chatItem.classList.add("active");

    // Update chat header with the selected chat name
    const name = chatItem.querySelector(".name").innerText;
    chatHeader.textContent = name;

    // Load messages (simulated for now)
    const messages = chatData[name] || [];
    chatMessages.innerHTML = ""; // Clear previous messages
    messages.forEach((msg) => {
      const div = document.createElement("div");
      div.classList.add("message", msg.type);
      div.textContent = msg.text;
      chatMessages.appendChild(div);
    });
  });
});
