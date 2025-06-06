document.forms.loginform.addEventListener("submit", async (event) => {
  event.preventDefault();

  const form = event.target;
  const email = form.querySelector('input[type="email"]').value;
  const password = form.querySelector('input[type="password"]').value;
  const errorMessage = document.getElementById("error-message");

  errorMessage.style.display = "none"; // Hide error message initially

  try {
    const response = await fetch("/api/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) {
      const error = await response.json();
      errorMessage.textContent = `Login failed: ${error.message}`;
      errorMessage.style.display = "block"; // Show error message
      return;
    }

    const data = await response.json();
    localStorage.setItem("token", data.token); // Save token in local storage
    // alert("Login successful!");
    // console.log("User data:", data);

    // Optionally redirect to another page
    window.location.href = "/chat";
  } catch (error) {
    errorMessage.textContent = "An error occurred during login.";
    errorMessage.style.display = "block"; // Show error message
    console.error("Login error:", error);
  }
});
