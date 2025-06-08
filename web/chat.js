var userid;
var chatList;
var chatID;
const socket = new WebSocket("ws://localhost:8080/api/ws");

async function main(params) {
  userid = await verifyToken();
  chatList = await getChatList();

  setupsocket();

  //verify socket connection
  socket.send(
    JSON.stringify({
      type: "verify",
      token: window.localStorage.getItem("token"),
    })
  );

  const chatHeader = document.querySelector(".chat-header h3");

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

    // Extract chatID from the clicked chat item
    chatID = chatItem.querySelector(".chatID").innerText;

    socket.send(
      JSON.stringify({
        type: "subscribe",
        chatID: parseInt(chatID),
      })
    );

    socket.send(
      JSON.stringify({
        type: "getmessages",
        chatID: parseInt(chatID),
      })
    );

    // Show the chat input after selecting a chat
    document.querySelector(".chat-input").classList.remove("hidden");
  });

  document.querySelector(".chat-input button").addEventListener("click", () => {
    const input = document.querySelector(".chat-input textarea"); // Revert to input
    const message = input.value;
    if (message) {
      // Send message to the server
      socket.send(
        JSON.stringify({
          type: "message",
          content: message,
          chatID: parseInt(chatID),
          messagetype: "text",
        })
      );

      showmessage(userid, message);

      input.value = ""; // Clear the input field
    }
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
            <p class="chatID">${chat.ID}</p>
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

function setupsocket() {
  // Connection opened
  socket.addEventListener("open", (event) => {
    // console.log("WebSocket connection established");
  });

  // Listen for messages from the server
  socket.addEventListener("message", (event) => {
    const data = JSON.parse(event.data);
    if (data.type === "message") {
      if (data.senderID != userid) {
        showmessage(data.senderID, data.text);
      }
    }
    if (data.type === "messages") {
      data.messages.forEach((message) => {
        showmessage(message.SenderID, message.Content);
      });
    }
  });
}

function showmessage(id, text) {
  var messagetype = "received"; // Default to received message
  if (id == userid) {
    messagetype = "sent";
  }
  // Display the sent message in the chat
  const chatMessages = document.querySelector(".chat-messages");
  const div = document.createElement("div");
  div.classList.add("message", messagetype);
  div.textContent = text;
  chatMessages.appendChild(div);
  chatMessages.scrollTop = chatMessages.scrollHeight; // Auto-scroll to the bottom
}

function scrollToBottom() {
  const chatMessages = document.querySelector(".chat-messages");
  chatMessages.scrollTop = chatMessages.scrollHeight; // Auto-scroll to the bottom
}

main();
