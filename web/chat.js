const chatItems = document.querySelectorAll('.chat-item');
const chatHeader = document.querySelector('.chat-header h3');
const chatMessages = document.querySelector('.chat-messages');

// Simulated chat data
const chatData = {
  Alice: [
    { type: 'received', text: "Hi! How's it going?" },
    { type: 'sent', text: "All good! You?" },
  ],
  Bob: [
    { type: 'received', text: "Let's catch up soon." },
    { type: 'sent', text: "Sure, just say when!" },
  ]
};

// Add click event listeners to chat items
chatItems.forEach(item => {
  item.addEventListener('click', () => {
    // Remove active class from all
    chatItems.forEach(i => i.classList.remove('active'));
    item.classList.add('active');

    const name = item.querySelector('.name').innerText;
    chatHeader.textContent = name;

    // Load messages
    const messages = chatData[name] || [];
    chatMessages.innerHTML = ''; // Clear previous
    messages.forEach(msg => {
      const div = document.createElement('div');
      div.classList.add('message', msg.type);
      div.textContent = msg.text;
      chatMessages.appendChild(div);
    });
  });
});
